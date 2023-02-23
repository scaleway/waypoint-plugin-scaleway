package container

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/waypoint-plugin-scaleway/internal/plugin"
)

// Platform is the Platform implementation for Scaleway Container.
type Platform struct {
	overrideHttpClient *http.Client
	PluginConfig       plugin.Config
	config             PlatformConfig
}

var (
	_ component.Configurable       = (*Platform)(nil)
	_ component.ConfigurableNotify = (*Platform)(nil)
	_ component.Platform           = (*Platform)(nil)
	_ component.Status             = (*Platform)(nil)
)

// PlatformConfig is the config for the Scaleway Container Platform
type PlatformConfig struct {
	// NamespaceID is the ID of the container namespace used to deploy container
	NamespaceID string `hcl:"namespace_id"`

	// Region where the container namespace is located, will default to profile's namespace
	Region string `hcl:"region,optional"`
	// Port is the listening port of your container, will default to API's default
	Port uint32 `hcl:"port,optional"`
	// Env is a map of static env variables to add to your container, static env variable are not secrets
	Env map[string]string `hcl:"env,optional"`

	// Timeout is the maximum amount of time in seconds during which your container can process a request before being stopped
	Timeout uint32 `hcl:"timeout,optional"`

	// Privacy mode of your container, defaults to public
	// your container may still remain public using waypoint url service
	Privacy string `hcl:"privacy,optional"`

	// MaxConcurrency is the maximum number of simultaneous requests your container can handle at the same time
	MaxConcurrency uint32 `hcl:"max_concurrency,optional"`

	MinScale uint32 `hcl:"min_scale,optional"`
	MaxScale uint32 `hcl:"max_scale,optional"`

	MemoryLimit uint32 `hcl:"memory_limit,optional"`
}

func (p *Platform) ConfigSet(i interface{}) error {
	cfg := i.(*PlatformConfig)

	privacy := containerSDK.ContainerPrivacy(cfg.Privacy)
	if privacy != containerSDK.ContainerPrivacyPublic &&
		privacy != containerSDK.ContainerPrivacyPrivate {
		return fmt.Errorf("invalid container privacy %q, only public or private allowed", privacy)
	}
	
	return nil
}

func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

func (p *Platform) StatusFunc() interface{} {
	return p.status
}

func (p *Platform) DeployFunc() interface{} {
	return p.deploy
}

func (p *Platform) DestroyFunc() interface{} {
	return p.destroy
}

func (p *Platform) resourceManager(
	log hclog.Logger,
	dcr *component.DeclaredResourcesResp,
	dtr *component.DestroyedResourcesResp,
) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(p.scalewayContainerAPI),
		resource.WithDeclaredResourcesResp(dcr),
		resource.WithDestroyedResourcesResp(dtr),
		resource.WithResource(resource.NewResource(
			resource.WithName("container"),
			resource.WithPlatform("scaleway"),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER),
			resource.WithState(&Resource_Container{}),
			resource.WithCreate(p.resourceContainerCreate),
			resource.WithStatus(p.resourceContainerStatus),
			resource.WithDestroy(p.resourceContainerDestroy),
		)),
	)
}

func (p *Platform) status(
	ctx context.Context,
	ji *component.JobInfo,
	ui terminal.UI,
	log hclog.Logger,
	container *Container,
) (*sdk.StatusReport, error) {
	sg := ui.StepGroup()
	s := sg.Add("Checking the status of the container deployment...")

	rm := p.resourceManager(log, nil, nil)

	if container.ResourceState == nil {
		s.Update("Creating state")
		err := rm.Resource("container").SetState(&Resource_Container{
			Id: container.Id,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to set resource container state: %w", err)
		}
	} else {
		s.Update("Loading state")
		if err := rm.LoadState(container.ResourceState); err != nil {
			return nil, err
		}
	}

	report, err := rm.StatusReport(ctx, log, sg, ui)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resource manager failed to generate resource statuses: %s", err)
	}

	//report.Health = sdk.StatusReport_READY
	//s.Update("Deployment no implemented: " + container.Image)
	s.Done()
	return report, nil
}

func (p *Platform) deploy(
	ui terminal.UI,
	log hclog.Logger,
	dcr *component.DeclaredResourcesResp,
	src *component.Source,
	img *docker.Image,
	ctx context.Context,
	deployConfig *component.DeploymentConfig,
) (*Container, error) {
	st := ui.Status()
	defer st.Close()
	st.Update("Deploying container")

	container := &Container{
		Region: p.config.Region,
	}

	id, err := component.Id()
	if err != nil {
		return nil, err
	}

	container.DeploymentId = id
	container.Name = strings.ToLower(fmt.Sprintf("%s-v%v", src.App, deployConfig.Sequence))

	rm := p.resourceManager(log, dcr, nil)

	err = rm.CreateAll(ctx, container, log, st, deployConfig, img)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	st.Step(terminal.StatusOK, "Created resource")

	container.ResourceState = rm.State()

	servState := rm.Resource("container").State().(*Resource_Container)
	if servState == nil {
		return nil, status.Errorf(codes.Internal, "service state is nil")
	}

	return container, nil
}

func (p *Platform) destroy(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	container *Container,
	dtr *component.DestroyedResourcesResp,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	rm := p.resourceManager(log, nil, dtr)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if container.ResourceState == nil {
		err := rm.Resource("deployment").SetState(&Resource_Container{
			Id:     container.Id,
			Region: container.Region,
		})
		if err != nil {
			return err
		}
	} else {
		// Load our set state
		if err := rm.LoadState(container.ResourceState); err != nil {
			return err
		}
	}

	// Destroy
	return rm.DestroyAll(ctx, log, sg, ui)
}
