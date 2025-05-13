---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_app_active_slot"
description: |-
  Manages a Web App Active Slot.
---

# azurerm_web_app_active_slot

Manages a Web App Active Slot.

## Example Usage

### Windows Web App

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_plan" "example" {
  name                = "example-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Windows"
  sku_name            = "P1v2"
}

resource "azurerm_windows_web_app" "example" {
  name                = "example-windows-web-app"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_service_plan.example.location
  service_plan_id     = azurerm_service_plan.example.id

  site_config {}
}

resource "azurerm_windows_web_app_slot" "example" {
  name           = "example-windows-web-app-slot"
  app_service_id = azurerm_windows_web_app.example.name

  site_config {}
}

resource "azurerm_web_app_active_slot" "example" {
  slot_id = azurerm_windows_web_app_slot.example.id

}
```

### Linux Web App

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_plan" "example" {
  name                = "example-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "P1v2"
}

resource "azurerm_linux_web_app" "example" {
  name                = "example-linux-web-app"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_service_plan.example.location
  service_plan_id     = azurerm_service_plan.example.id

  site_config {}
}

resource "azurerm_linux_web_app_slot" "example" {
  name             = "example-linux-web-app-slot"
  app_service_name = azurerm_linux_web_app.example.name
  location         = azurerm_service_plan.example.location
  service_plan_id  = azurerm_service_plan.example.id

  site_config {}
}

resource "azurerm_web_app_active_slot" "example" {
  slot_id = azurerm_linux_web_app_slot.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `slot_id` - (Required) The ID of the Slot to swap with `Production`.

---

* `overwrite_network_config` - (Optional) The swap action should overwrite the Production slot's network configuration with the configuration from this slot. Defaults to `true`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Web App Active Slot

* `last_successful_swap` - The timestamp of the last successful swap with `Production`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web App Active Slot.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web App Active Slot.
* `update` - (Defaults to 30 minutes) Used when updating the Web App Active Slot.
* `delete` - (Defaults to 5 minutes) Used when deleting the Web App Active Slot.

## Import

a Web App Active Slot can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_app_active_slot.example "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1"
```
