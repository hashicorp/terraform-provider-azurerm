---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_linux_web_app_site_container"
description: |-
  Manages a Site Container for a Linux Web App.
---

# azurerm_linux_web_app_site_container

Manages a Site Container (sidecar) for a Linux Web App.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_plan" "example" {
  name                = "example-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Linux"
  sku_name            = "P1v2"
}

resource "azurerm_linux_web_app" "example" {
  name                = "example-web-app"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id

  site_config {
    application_stack {
      site_containers_enabled = true
    }
  }
}

resource "azurerm_linux_web_app_site_container" "example" {
  name             = "main"
  linux_web_app_id = azurerm_linux_web_app.example.id
  image            = "mcr.microsoft.com/appsvc/sample-hello-world:latest"
  target_port      = 80
  primary          = true
}

resource "azurerm_linux_web_app_site_container" "sidecar" {
  name             = "sidecar"
  linux_web_app_id = azurerm_linux_web_app.example.id
  image            = "mcr.microsoft.com/appsvc/sample-hello-world:latest"
  target_port      = 8080
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Site Container. Changing this forces a new resource to be created.

-> **Note:** The `name` must start and end with an alphanumeric character and may contain hyphens.

* `linux_web_app_id` - (Required) The ID of the Linux Web App that this Site Container is associated with. Changing this forces a new resource to be created.

* `image` - (Required) The fully qualified container image (including tag) that should run inside this Site Container.

* `target_port` - (Required) The port exposed by the container image that should receive traffic. Possible values range between `1` and `65535`.

---

* `authentication_type` - (Optional) The authentication strategy used to pull the image. Possible values are `Anonymous`, `SystemIdentity`, `UserAssigned`, and `UserCredentials`. Defaults to `Anonymous`.

* `environment_variable` - (Optional) One or more `environment_variable` blocks as defined below.

* `password_secret` - (Optional) The password to use when `authentication_type` is set to `UserCredentials`.

-> **Note:** Azure does not return values supplied to `password_secret`, so Terraform cannot detect drift for this property.

* `primary` - (Optional) Should this container serve the primary site traffic? Defaults to `false`.

* `startup_command` - (Optional) The command that should be executed when the container starts.

* `user_managed_identity_client_id` - (Optional) The Client ID of the user-assigned managed identity that should be used when `authentication_type` is set to `UserAssigned`.

* `username` - (Optional) The username to use when `authentication_type` is set to `UserCredentials`.

~> **Note:** When `authentication_type` is set to `UserCredentials`, both `username` and `password_secret` must be specified. When `authentication_type` is set to `UserAssigned`, `user_managed_identity_client_id` must be specified.

* `volume_mount` - (Optional) One or more `volume_mount` blocks as defined below.

---

An `environment_variable` block supports the following:

* `app_setting_name` - (Required) The name of an App Setting on the parent Linux Web App whose value is exposed to the container as the environment variable named above. The actual value is resolved from the App Setting at runtime; if the App Setting is not defined the environment variable is set to an empty string.

* `name` - (Required) The name of the environment variable as it appears inside the container.

---

A `volume_mount` block supports the following:

* `container_mount_path` - (Required) The absolute path inside the container where the volume is mounted.

* `data` - (Optional) The opaque data supplied to Azure for the mount. The contents depend on the selected storage option.

* `read_only` - (Optional) Should the mounted volume be read only? Defaults to `false`.

* `volume_sub_path` - (Optional) The path inside the Linux Web App volume that should be exposed to the container.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Linux Web App Site Container.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Linux Web App Site Container.
* `read` - (Defaults to 5 minutes) Used when retrieving the Linux Web App Site Container.
* `update` - (Defaults to 30 minutes) Used when updating the Linux Web App Site Container.
* `delete` - (Defaults to 5 minutes) Used when deleting the Linux Web App Site Container.

## Import

Linux Web App Site Containers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_linux_web_app_site_container.example "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Web/sites/site1/sitecontainers/container1"
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Web` - 2023-12-01
