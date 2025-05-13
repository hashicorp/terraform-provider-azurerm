---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center_project_pool"
description: |-
  Manages a Dev Center Project Pool.
---

# azurerm_dev_center_project_pool

Manages a Dev Center Project Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dev_center" "example" {
  name                = "example-dc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_dev_center_network_connection" "example" {
  name                = "example-dcnc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  subnet_id           = azurerm_subnet.example.id
  domain_join_type    = "AzureADJoin"
}

resource "azurerm_dev_center_attached_network" "example" {
  name                  = "example-dcet"
  dev_center_id         = azurerm_dev_center.example.id
  network_connection_id = azurerm_dev_center_network_connection.example.id
}

resource "azurerm_dev_center_project" "example" {
  name                = "example-dcp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  dev_center_id       = azurerm_dev_center.example.id
}

resource "azurerm_dev_center_dev_box_definition" "example" {
  name               = "example-dcet"
  location           = azurerm_resource_group.example.location
  dev_center_id      = azurerm_dev_center.example.id
  image_reference_id = "${azurerm_dev_center.example.id}/galleries/default/images/microsoftvisualstudio_visualstudioplustools_vs-2022-ent-general-win10-m365-gen2"
  sku_name           = "general_i_8c32gb256ssd_v2"
}

resource "azurerm_dev_center_project_pool" "example" {
  name                                    = "example-dcpl"
  location                                = azurerm_resource_group.example.location
  dev_center_project_id                   = azurerm_dev_center_project.example.id
  dev_box_definition_name                 = azurerm_dev_center_dev_box_definition.example.name
  local_administrator_enabled             = true
  dev_center_attached_network_name        = azurerm_dev_center_attached_network.example.name
  stop_on_disconnect_grace_period_minutes = 60
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Dev Center Project Pool. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Dev Center Project Pool should exist. Changing this forces a new resource to be created.

* `dev_center_project_id` - (Required) The ID of the associated Dev Center Project. Changing this forces a new resource to be created.

* `dev_box_definition_name` - (Required) The name of the Dev Center Dev Box Definition.

* `local_administrator_enabled` - (Required) Specifies whether owners of Dev Boxes in the Dev Center Project Pool are added as local administrators on the Dev Box.

* `dev_center_attached_network_name` - (Required) The name of the Dev Center Attached Network in parent Project of the Dev Center Project Pool.

* `stop_on_disconnect_grace_period_minutes` - (Optional) The specified time in minutes to wait before stopping a Dev Center Dev Box once disconnect is detected. Possible values are between `60` and `480`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Dev Center Project Pool.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Project Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center Project Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Project Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Dev Center Project Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center Project Pool.

## Import

An existing Dev Center Project Pool can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_project_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevCenter/projects/project1/pools/pool1
```
