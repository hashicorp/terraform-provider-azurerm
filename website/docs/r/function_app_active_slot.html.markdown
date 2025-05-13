---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app_active_slot"
description: |-
  Manages a Function App Active Slot.
---

# azurerm_function_app_active_slot

Manages a Function App Active Slot.

## Example Usage

### Windows Function App

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "windowsfunctionappsa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "example" {
  name                = "example-app-service-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Windows"
  sku_name            = "Y1"
}

resource "azurerm_windows_function_app" "example" {
  name                 = "example-windows-function-app"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  storage_account_name = azurerm_storage_account.example.name
  service_plan_id      = azurerm_service_plan.example.id

  site_config {}
}

resource "azurerm_windows_function_app_slot" "example" {
  name                 = "example-windows-function-app-slot"
  function_app_id      = azurerm_windows_function_app.example.id
  storage_account_name = azurerm_storage_account.example.name

  site_config {}
}

resource "azurerm_function_app_active_slot" "example" {
  slot_id = azurerm_windows_function_app_slot.example.id
}
```

### Linux Function App

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "linuxfunctionappsa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "example" {
  name                = "example-app-service-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "Y1"
}

resource "azurerm_linux_function_app" "example" {
  name                 = "example-linux-function-app"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  service_plan_id      = azurerm_service_plan.example.id
  storage_account_name = azurerm_storage_account.example.name

  site_config {}
}

resource "azurerm_linux_function_app_slot" "example" {
  name                 = "example-linux-function-app-slot"
  function_app_id      = azurerm_linux_function_app.example.name
  storage_account_name = azurerm_storage_account.example.name

  site_config {}
}

resource "azurerm_function_app_active_slot" "example" {
  slot_id = azurerm_linux_function_app_slot.example.id
}

```

## Arguments Reference

The following arguments are supported:

* `slot_id` - (Required) The ID of the Slot to swap with `Production`.

---

* `overwrite_network_config` - (Optional) The swap action should overwrite the Production slot's network configuration with the configuration from this slot. Defaults to `true`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Function App Active Slot

* `last_successful_swap` - The timestamp of the last successful swap with `Production`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Function App Active Slot.
* `read` - (Defaults to 5 minutes) Used when retrieving the Function App Active Slot.
* `update` - (Defaults to 30 minutes) Used when updating the Function App Active Slot.
* `delete` - (Defaults to 5 minutes) Used when deleting the Function App Active Slot.

## Import

a Function App Active Slot can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_function_app_active_slot.example "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1"
```
