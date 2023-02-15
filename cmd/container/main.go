package main

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
	"github.com/scaleway/waypoint-plugin-scaleway/container"
	"github.com/scaleway/waypoint-plugin-scaleway/internal/plugin"
)

var (
	Version = ""
)

func main() {
	pluginConfig := plugin.InitConfig(Version, "container")
	sdk.Main(sdk.WithComponents(
		&container.Platform{
			PluginConfig: pluginConfig,
		},
	))
}
