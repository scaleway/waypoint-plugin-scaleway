# Container

## Deployment

### Attributes

#### Required

- `namespace_id` ID of the container namespace to use when deploying your containers

#### Optional

- `port` - The listening port of your container
- `env` - Static non-secret env variables
- `timeout` - Maximum amount of time in seconds for requests before being stopped
- `privacy` - Privacy mode of your container, public or private, defaults to public
- `max_concurrency` - Maximum number of simultaneous requests your container can handle at the same time, defaults to 50
- `min_scale` - Minimum scaling value of your container, defaults to 0
    
    __<!>__ If your container scales down to 0, the waypoint entrypoint will not be reachable

- `max_scale` - Maximum scaling value of your container, defaults to 5
- `memory_limit` - Memory allocated to your container, defaults to 256Mi, this is the value that change the price you pay per container.
- `region` - Region of your container namespace, defaults to your scw [profile](scw-config.md)'s default region
- `profile` - The [config](scw-config.md)'s profile to use

### Examples

```hcl
deploy {
  use "scaleway-container" {
    profile = "dev"
    namespace_id = "xxxx-xxxx-xxx-xxxx"
    region = "nl-ams"

    port = 80
    
    env = {
      key = "value"
    }
    
    privacy = "public"
    min_scale = 0
    max_scale = 1
    memory_limit = 128
  }
}
```

### References

- [Container API](https://developers.scaleway.com/en/products/containers/api/#introduction)
- [How am I billed for Serverless Containers?](https://www.scaleway.com/en/docs/faq/serverless-containers/#how-am-i-billed-for-serverless-containers)
