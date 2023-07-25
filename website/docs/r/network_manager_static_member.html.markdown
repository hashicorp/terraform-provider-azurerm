---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_static_member"
description: |-
  Manages a Network Manager Static Member.
---

# azurerm_network_manager_static_member

Manages a Network Manager Static Member.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin"]
  description    = "example network manager"
}

resource "azurerm_network_manager_network_group" "example" {
  name               = "example-group"
  network_manager_id = azurerm_network_manager.example.id
  description        = "example network group"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.example.location
}

resource "azurerm_network_manager_static_member" "example" {
  name                      = "example-nmsm"
  network_group_id          = azurerm_network_manager_network_group.example.id
  target_virtual_network_id = azurerm_virtual_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager Static Member. Changing this forces a new Network Manager Static Member to be created.

* `network_group_id` - (Required) Specifies the ID of the Network Manager Group. Changing this forces a new Network Manager Static Member to be created.

* `target_virtual_network_id` - (Required) Specifies the Resource ID of the Virtual Network using as the Static Member. Changing this forces a new Network Manager Static Member to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Static Member.
 
* `region` - The region of the Network Manager Static Member.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Static Member.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Static Member.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Static Member.

## Import

Network Manager Static Member can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_static_member.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/networkGroups/networkGroup1/staticMembers/staticMember1
```
