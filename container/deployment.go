package container

import (
	"github.com/hashicorp/waypoint-plugin-sdk/component"
)

var (
	_ component.Deployment        = (*Container)(nil)
	_ component.DeploymentWithUrl = (*Container)(nil)
)

func (c *Container) URL() string {
	return c.Url
}
