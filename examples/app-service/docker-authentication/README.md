# Example: a Linux App Service running a Docker container with AD and Microsoft authentication enabled.

This example provisions a Linux App Service which runs a single Docker container with AD and Microsoft authentication enabled.

### Notes

* The Container is launched on the first HTTP Request, which can take a while.
* Continuous Deployment of a single Docker Container can be achieved using the App Setting `DOCKER_ENABLE_CI` to `true`.

```hcl
resource "azurerm_linux_web_app" "example" {
  # ...
  site_config {
    application_stack {
      docker_image     = "jackofallops/azure-containerapps-python-acctest"
      docker_image_tag = "v0.0.1"
    }
  }
}
```
