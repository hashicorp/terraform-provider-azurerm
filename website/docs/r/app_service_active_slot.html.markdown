---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_active_slot"
sidebar_current: "docs-azurerm-resource-app-service-active-slot"
description: |-
  Swaps an App Service Slot (within an App Service) with the Production slot.

---

# azurerm_app_service_slot

Swaps an App Service Slot (within an App Service) with the Production slot.

## Example Usage

```hcl
resource "random_id" "server" {
  # ...
}

resource "azurerm_resource_group" "test" {
  # ...
}

resource "azurerm_app_service_plan" "test" {
  # ...
}

resource "azurerm_app_service" "test" {
  # ...
}

resource "azurerm_app_service_slot" "test" {
  # ...
}

resource "azurerm_app_service_active_slot" "test" {
    resource_group_name = "${azurerm_resource_group.test.name}"
    app_service_name    = "${azurerm_app_service.test.name}"
    preserve_vnet       = true
    source_slot_name    = "${azurerm_app_service_slot.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which to create the App Service Slot component.

* `app_service_name` - (Required) The name of the App Service within which to create the App Service Slot.  Changing this forces a new resource to be created.

* `preserve_vnet` - (Required) `true` to preserve Virtual Network to the slot during swap; otherwise, `false`.

* `source_slot_name` - (Required) Source deployment slot for the swap operation.
