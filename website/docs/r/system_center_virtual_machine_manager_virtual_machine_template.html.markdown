---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_virtual_machine_template"
description: |-
  Manages a System Center Virtual Machine Manager Virtual Machine Template.
---

# azurerm_system_center_virtual_machine_manager_virtual_machine_template

Manages a System Center Virtual Machine Manager Virtual Machine Template.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_system_center_virtual_machine_manager_server" "example" {
  name                = "example-scvmmms"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ExtendedLocation/customLocations/customLocation1"
  fqdn                = "example.labtest"
  username            = "testUser"
  password            = "H@Sh1CoR3!"
}

data "azurerm_system_center_virtual_machine_manager_inventory_items" "example" {
  inventory_type                                  = "VirtualMachineTemplate"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.example.id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_template" "example" {
  name                                                           = "example-scvmmvmtemplate"
  resource_group_name                                            = azurerm_resource_group.example.name
  location                                                       = azurerm_resource_group.example.location
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.example.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.example.inventory_items[0].id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the System Center Virtual Machine Manager Virtual Machine Template. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the System Center Virtual Machine Manager Virtual Machine Template should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the System Center Virtual Machine Manager Virtual Machine Template should exist. Changing this forces a new resource to be created.

* `custom_location_id` - (Required) The ID of the Custom Location for the System Center Virtual Machine Manager Virtual Machine Template. Changing this forces a new resource to be created.

* `system_center_virtual_machine_manager_server_inventory_item_id` - (Required) The ID of the System Center Virtual Machine Manager Server Inventory Item. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the System Center Virtual Machine Manager Virtual Machine Template.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the System Center Virtual Machine Manager Virtual Machine Template.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the System Center Virtual Machine Manager Virtual Machine Template.
* `read` - (Defaults to 5 minutes) Used when retrieving the System Center Virtual Machine Manager Virtual Machine Template.
* `update` - (Defaults to 30 minutes) Used when updating the System Center Virtual Machine Manager Virtual Machine Template.
* `delete` - (Defaults to 30 minutes) Used when deleting the System Center Virtual Machine Manager Virtual Machine Template.

## Import

System Center Virtual Machine Manager Virtual Machine Templates can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_system_center_virtual_machine_manager_virtual_machine_template.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.ScVmm/virtualMachineTemplates/virtualMachineTemplate1
```
