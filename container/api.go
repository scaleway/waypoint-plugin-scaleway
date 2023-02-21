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
	client, err := scw.NewClient(
		scw.WithEnv(),
		scw.WithProfile(&cfg.Profile),
		scw.WithUserAgent(p.PluginConfig.UserAgent))
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

func createContainerValue(value uint32) *uint32 {
	if value == 0 {
		return nil
	}
	return &value
}
