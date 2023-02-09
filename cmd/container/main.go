package main

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
	"github.com/scaleway/waypoint-plugin-scaleway/container"
)

func main() {
	sdk.Main(sdk.WithComponents(
		&container.Platform{},
	))
}
