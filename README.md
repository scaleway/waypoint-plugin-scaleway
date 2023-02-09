# Waypoint Plugins Scaleway

Plugins for waypoint that add support for Scaleway
Currently the only plugin is container

## Usage

```hcl
deploy {
  use "scaleway-container" {
    port = 1337
    namespace = "xxxx-xxxx-xxx-xxxx"
    region = "fr-par"
  }
}
```
API keys are loaded from [scaleway's config](https://github.com/scaleway/scaleway-sdk-go/tree/master/scw#scaleway-config) default profile and override by environment variables

## Install

### From sources

- `make all install`
