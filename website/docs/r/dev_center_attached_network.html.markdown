---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center_attached_network"
description: |-
  Manages a Dev Center Attached Network.
---

# azurerm_dev_center_attached_network

Manages a Dev Center Attached Network.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-dcan"
  location = "West Europe"
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

resource "azurerm_dev_center" "example" {
  name                = "example-dc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_dev_center_network_connection" "example" {
  name                = "example-dcnc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  domain_join_type    = "AzureADJoin"
  subnet_id           = azurerm_subnet.example.id
}

resource "azurerm_dev_center_attached_network" "example" {
  name                  = "example-dcet"
  dev_center_id         = azurerm_dev_center.example.id
  network_connection_id = azurerm_dev_center_network_connection.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Dev Center Attached Network. Changing this forces a new resource to be created.

* `dev_center_id` - (Required) The ID of the associated Dev Center. Changing this forces a new resource to be created.

* `network_connection_id` - (Required) The ID of the Dev Center Network Connection you want to attach. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Attached Network.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center Attached Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Attached Network.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center Attached Network.

## Import

An existing Dev Center Attached Network can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_attached_network.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevCenter/devCenters/dc1/attachedNetworks/et1
```
