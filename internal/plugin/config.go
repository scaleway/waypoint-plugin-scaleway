package plugin

import "fmt"

type Config struct {
	Version   string
	UserAgent string
}

func InitConfig(version string, pluginName string) Config {
	formattedVersion := formatVersion(version)
	return Config{
		Version:   formattedVersion,
		UserAgent: fmt.Sprintf("waypoint-plugin-scaleway-%s/%s", pluginName, formattedVersion),
	}
}

func formatVersion(version string) string {
	if version == "" {
		return "dev"
	}
	return version
}
