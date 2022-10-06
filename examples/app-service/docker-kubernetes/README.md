# Example: a Linux App Service running multiple containers from a Kubernetes Manifest

This example provisions a Linux App Service which runs multiple Docker Containers from a Kubernetes Manifest.

### Notes

* The Container is launched on the first HTTP Request, which can take a while.
* If you're not using App Service Slots and Deployments are handled outside of Terraform - [it's possible to ignore changes to specific fields in the configuration using `ignore_changes` within Terraform's `lifecycle` block](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changes), for example:

```hcl
resource "azurerm_app_service" "test" {
  # ...
  site_config = {
    # ...
    linux_fx_version = "KUBE|${filebase64("kubernetes.yml")}"
  }

  lifecycle {
    ignore_changes = [
      "site_config.0.linux_fx_version", # deployments are made outside of Terraform
    ]
  }
}
```
