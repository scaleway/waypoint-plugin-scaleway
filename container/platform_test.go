package container

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/scaleway/waypoint-plugin-scaleway/internal/plugin"
	"github.com/scaleway/waypoint-plugin-scaleway/internal/plugintesting"
)

func TestDeploy(t *testing.T) {
	tt := plugintesting.Init(t)
	defer tt.CleanUp()
	namespace := tt.ContainerNamespace()

	imageTag, err := tt.UploadTestImage(namespace)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to upload test image to registry namespace (%s): %w", namespace.RegistryNamespaceID, err))
	}

	p := Platform{
		PluginConfig: plugin.InitConfig("test", "container"),
		config: PlatformConfig{
			Namespace: namespace.ID,
			Port:      80,
		},
		overrideHttpClient: tt.HttpClient,
	}

	semicolonIndex := strings.Index(imageTag, ":")
	img := docker.Image{
		Image:        imageTag[:semicolonIndex],
		Tag:          imageTag[semicolonIndex+1:],
		Architecture: "amd64",
		Location: &docker.Image_Registry{
			Registry: &docker.Image_RegistryLocation{
				WaypointGenerated: false,
			},
		},
	}

	logger := hclog.New(nil)

	src := component.Source{
		App:  "app-name",
		Path: "path-name",
	}

	ui := terminal.ConsoleUI(context.Background())
	dcrResp := component.DeclaredResourcesResp{}
	deployConfig := component.DeploymentConfig{}
	_, err = p.deploy(ui, logger, &dcrResp, &src, &img, context.Background(), &deployConfig)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to deploy container: %w", err))
	}
}
