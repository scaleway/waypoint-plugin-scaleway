package container

import (
	"fmt"

	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (p *Platform) scalewayClient() (*scw.Client, error) {
	cfg, err := scw.LoadConfig()
	if _, isNotFoundError := err.(*scw.ConfigFileNotFoundError); isNotFoundError {
		cfg = &scw.Config{}
	} else if err != nil {
		return nil, fmt.Errorf("failed to load scaleway's config: %w", err)
	}
	clientOptions := []scw.ClientOption{
		scw.WithProfile(&cfg.Profile),
		scw.WithEnv(),
		scw.WithUserAgent(p.PluginConfig.UserAgent),
	}
	if p.overrideHttpClient != nil {
		clientOptions = append(clientOptions, scw.WithHTTPClient(p.overrideHttpClient))
	}
	client, err := scw.NewClient(clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to init scaleway's client: %w", err)
	}
	return client, err
}

func (p *Platform) scalewayContainerAPI() (*containerSDK.API, error) {
	client, err := p.scalewayClient()
	if err != nil {
		return nil, err
	}
	return containerSDK.NewAPI(client), nil
}

func prepareContainerCreateRequest(
	name string,
	registryImage string,
	cfg *PlatformConfig,
	entrypointEnv map[string]string,
) *containerSDK.CreateContainerRequest {
	containerSecretEnv := []*containerSDK.Secret(nil)

	for key, value := range entrypointEnv {
		containerSecretEnv = append(containerSecretEnv, &containerSDK.Secret{
			Key:   key,
			Value: scw.StringPtr(value),
		})
	}

	req := &containerSDK.CreateContainerRequest{
		Region:                     scw.Region(cfg.Region),
		NamespaceID:                cfg.NamespaceID,
		Name:                       name,
		EnvironmentVariables:       &cfg.Env,
		RegistryImage:              &registryImage,
		Port:                       createContainerValue(cfg.Port),
		SecretEnvironmentVariables: containerSecretEnv,
		MinScale:                   createContainerValue(cfg.MinScale),
		MaxScale:                   createContainerValue(cfg.MaxScale),
		MemoryLimit:                createContainerValue(cfg.MemoryLimit),
		MaxConcurrency:             createContainerValue(cfg.MaxConcurrency),
	}

	if cfg.Privacy != "" {
		req.Privacy = containerSDK.ContainerPrivacy(cfg.Privacy)
	}

	if cfg.Timeout != 0 {
		req.Timeout = &scw.Duration{
			Seconds: int64(cfg.Timeout),
		}
	}

	return req
}

func createContainerValue(value uint32) *uint32 {
	if value == 0 {
		return nil
	}
	return &value
}
