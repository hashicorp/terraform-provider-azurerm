---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_active_slot"
description: |-
  Promotes an App Service Slot to Production within an App Service

---

# azurerm_app_service_active_slot

Promotes an App Service Slot to Production within an App Service.

!> **Note:** This resource has been deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use [`azurerm_web_app_active_slot`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/web_app_active_slot) resource instead.

-> **Note:** When using Slots - the `app_settings`, `connection_string` and `site_config` blocks on the `azurerm_app_service` resource will be overwritten when promoting a Slot using the `azurerm_app_service_active_slot` resource.

## Example Usage

```hcl
resource "random_id" "server" {
  # ...
}

resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_app_service_plan" "example" {
  # ...
}

resource "azurerm_app_service" "example" {
  # ...
}

resource "azurerm_app_service_slot" "example" {
  # ...
}

resource "azurerm_app_service_active_slot" "example" {
  resource_group_name   = azurerm_resource_group.example.name
  app_service_name      = azurerm_app_service.example.name
  app_service_slot_name = azurerm_app_service_slot.example.name
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which the App Service exists. Changing this forces a new resource to be created.

* `app_service_name` - (Required) The name of the App Service within which the Slot exists. Changing this forces a new resource to be created.

* `app_service_slot_name` - (Required) The name of the App Service Slot which should be promoted to the Production Slot within the App Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Active Slot.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Active Slot.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Active Slot.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Active Slot.
