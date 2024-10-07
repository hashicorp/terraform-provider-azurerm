---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_inventory_items"
description: |-
  Gets information about existing System Center Virtual Machine Manager Inventory Items.
---

# Data Source: azurerm_system_center_virtual_machine_manager_inventory_items

Use this data source to access information about existing System Center Virtual Machine Manager Inventory Items.

## Example Usage

```hcl
data "azurerm_system_center_virtual_machine_manager_inventory_items" "example" {
  inventory_type                                  = "Cloud"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.example.id
}
```

## Argument Reference

* `inventory_type` - (Required) The inventory type of the System Center Virtual Machine Manager Inventory Item. Possible values are `Cloud`, `VirtualMachine`, `VirtualMachineTemplate` and `VirtualNetwork`.

* `system_center_virtual_machine_manager_server_id` - (Required) The ID of the System Center Virtual Machine Manager Server.

## Attributes Reference

* `id` - The ID of the System Center Virtual Machine Manager Inventory Items.

* `inventory_items` - One or more `inventory_items` blocks as defined below.

---

A `inventory_items` block exports the following:

* `id` - The ID of the System Center Virtual Machine Manager Inventory Item.

* `name` - The name of the System Center Virtual Machine Manager Inventory Item.

* `uuid` - The UUID of the System Center Virtual Machine Manager Inventory Item that is assigned by System Center Virtual Machine Manager.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the System Center Virtual Machine Manager Inventory Items.
