# Example: a Linux App Service running a Docker container with AD and Microsoft authentication enabled.

This example provisions a Linux App Service which runs a single Docker container with AD and Microsoft authentication enabled.

### Notes

* The Container is launched on the first HTTP Request, which can take a while.
* Continuous Deployment of a single Docker Container can be achieved using the App Setting `DOCKER_ENABLE_CI` to `true`.
* If you're not using App Service Slots and Deployments are handled outside of Terraform - [it's possible to ignore changes to specific fields in the configuration using `ignore_changes` within Terraform's `lifecycle` block](https://www.terraform.io/docs/configuration/resources.html#lifecycle), for example:

```hcl
resource "azurerm_app_service" "test" {
  # ...
  site_config = {
    # ...
    linux_fx_version = "DOCKER|appsvcsample/python-helloworld:0.1.2"
  }

  lifecycle {
    ignore_changes = [
      "site_config.0.linux_fx_version", # deployments are made outside of Terraform
    ]
  }
}
```
