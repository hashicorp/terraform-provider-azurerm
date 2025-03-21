---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_static_web_app_link_backend"
description: |-
  Manages a Static Web App Link Backend.
---

# azurerm_static_web_app_link_backend

Manages an App Service Static Web App Link Backend.

~> **NOTE:** This resource registers the specified API Management, App Service, or Container App to the `Production` build of the Static Web App.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "westus"
}

resource "azurerm_static_web_app" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_tier            = "Standard"
  sku_size            = "Standard"
}

resource "azurerm_api_management" "example" {
  name                = "example-api-management"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_static_web_app_link_backend" "example" {
  static_web_app_id   = azurerm_static_web_app.example.id
  backend_resource_id = azurerm_api_management.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `static_web_app_id` - (Required) The ID of the Static Web App to register the API Management, App Service, or Container App to as a backend. Changing this forces a new Static Web App Link Backend to be created.

* `backend_resource_id` - (Required) The ID of an API Management, App Service, or Container App to connect to the Static Web App as a Backend. Changing this forces a new Static Web App Link Backend to be created.

~> **NOTE:** Only one Backend resource can be connected to a Static Web App. Multiple resources are not currently supported.

~> **NOTE:** Connecting an App Service resource to a Static Web App resource updates the App Service to use AuthV2 and configures the `azure_static_web_app_v2` which may need to be accounted for by the use of `ignore_changes` depending on the existing `auth_settings_v2` configuration of the target App Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Static Web App Link Backend.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Static Web App Link Backend.
* `read` - (Defaults to 5 minutes) Used when retrieving the Static Web App Link Backend.
* `delete` - (Defaults to 30 minutes) Used when deleting the Static Web App Link Backend.

## Import

Static Web App Link Backend can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_static_web_app_link_backend.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/staticSites/swa1/linkedBackends/linkedbackend1
```
