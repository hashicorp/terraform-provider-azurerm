---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_cloud"
description: |-
  Manages a System Center Virtual Machine Manager Cloud.
---

# azurerm_system_center_virtual_machine_manager_cloud

Manages a System Center Virtual Machine Manager Cloud.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_custom_location" "example" {
  name = "example-cl"
}

resource "azurerm_system_center_virtual_machine_manager_server" "example" {
  name                = "example-scvmms"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = azurerm_custom_location.example.id
  fqdn                = "exampledomain.com"
}

resource "azurerm_system_center_virtual_machine_manager_server_inventory_item" "example" {
  name                = "example-scvmmsii"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_system_center_virtual_machine_manager_cloud" "example" {
  name                                                           = "example-scvmmc"
  location                                                       = azurerm_resource_group.example.location
  resource_group_name                                            = azurerm_resource_group.example.name
  custom_location_id                                             = azurerm_custom_location.example.id
  system_center_virtual_machine_manager_server_id                = azurerm_system_center_virtual_machine_manager_server.example.id
  system_center_virtual_machine_manager_server_inventory_item_id = azurerm_system_center_virtual_machine_manager_server_inventory_item.example.id
  uuid                                                           = "37437563-5205-4fd3-aa74-6bc16702748a"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this System Center Virtual Machine Manager Cloud. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the System Center Virtual Machine Manager Cloud should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group within which this System Center Virtual Machine Manager Cloud should exist. Changing this forces a new resource to be created.

* `custom_location_id` - (Required) The ID of the Custom Location for the System Center Virtual Machine Manager Cloud.

* `system_center_virtual_machine_manager_server_id` - (Optional) The ID of the System Center Virtual Machine Manager Server.

* `system_center_virtual_machine_manager_server_inventory_item_id` - (Optional) The ID of the System Center Virtual Machine Manager Server Inventory Item.

* `uuid` - (Optional) The unique ID of the System Center Virtual Machine Manager Cloud.

* `tags` - (Optional) A mapping of tags which should be assigned to the System Center Virtual Machine Manager Cloud.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the System Center Virtual Machine Manager Cloud.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this System Center Virtual Machine Manager Cloud.
* `read` - (Defaults to 5 minutes) Used when retrieving this System Center Virtual Machine Manager Cloud.
* `update` - (Defaults to 30 minutes) Used when updating this System Center Virtual Machine Manager Cloud.
* `delete` - (Defaults to 30 minutes) Used when deleting this System Center Virtual Machine Manager Cloud.

## Import

System Center Virtual Machine Manager Clouds can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_system_center_virtual_machine_manager_cloud.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.ScVmm/clouds/cloud1
```
