# Waypoint Plugins Scaleway

Plugins for waypoint that add support for Scaleway.
Currently, the only plugin available is container.

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
API keys are loaded from [Scaleway's config](https://github.com/scaleway/scaleway-sdk-go/tree/master/scw#scaleway-config) default profile and can be overwritten by environment variables.

A list of all options can be found in [container's documentation](./docs/container.md)

## Install

### From releases

- [Download the zip](https://github.com/scaleway/waypoint-plugin-scaleway/releases) of the latest version for your architecture.
- Unzip the plugin by running the following command: `unzip waypoint-plugin-scaleway-container_*.zip -d ~/.config/.waypoint/plugins/`

> **Note**
> On macOS, you will have to execute the following command to ignore Apple's developer authenticity verification:
> ```
> xattr -d com.apple.quarantine ~/.config/.waypoint/plugins/waypoint-plugin-scaleway-container
> ```

### From sources
- Clone the repository on your local machine.
- Run the following command to build and install the plugin: `make all install`

