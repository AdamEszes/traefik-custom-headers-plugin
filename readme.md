# Rename Header

Traefik custom headers plugin is a middleware plugin for [Traefik](https://traefik.io) which renames headers in the response, while keeping their values.

## Configuration

### Static

```yaml
pilot:
  token: "xxxx"

experimental:
  plugins:
    renameHeaders:
      modulename: "gitlab.com/traefik-custom-headers-plugin/traefik-custom-headers-plugin"
      version: "v0.0.1"
```

### Dynamic

To configure the Rename Headers plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). 
The following example creates and uses the renameHeaders middleware plugin to rename the "custom_id" header

```yaml
http:
  routes:
    my-router:
      rule: "Host(`localhost`)"
      service: "my-service"
      middlewares : 
        - "renameHeaders"
  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1"
  middlewares:
    renameHeaders:
      plugin:
        renameData:
          - existingHeaderName: "Customheader"
          newHeaderName: "customheader"
```
