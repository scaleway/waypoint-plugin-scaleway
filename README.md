# Waypoint Plugins Scaleway

Plugins for waypoint that add support for Scaleway
Currently the only plugin is container

## Usage

```hcl
deploy {
  use "scaleway-container" {
    port = 80
    namespace_id = "xxxx-xxxx-xxx-xxxx"
    region = "fr-par"
  }
}
```
API keys are loaded from [scaleway's config](https://github.com/scaleway/scaleway-sdk-go/tree/master/scw#scaleway-config) default profile and override by environment variables

List of all options can be found in [container's documentation](./docs/container.md)

## Install

### From releases

- [Download the zip](https://github.com/scaleway/waypoint-plugin-scaleway/releases) of the latest version for your architecture
- `unzip waypoint-plugin-scaleway-container_*.zip -d ~/.config/.waypoint/plugins/`

> **Note**
> On macOS, you will have to execute the following command to ignore developer authenticity verification :
> ```
> xattr -d com.apple.quarantine ~/.config/.waypoint/plugins/waypoint-plugin-scaleway-container
> ```

### From sources

- `make all install`

