---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_peering"
description: |-
  Manages a virtual network peering which allows resources to access other
  resources in the linked virtual network.
---

# azurerm_virtual_network_peering

Manages a virtual network peering which allows resources to access other
resources in the linked virtual network.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "peeredvnets-rg"
  location = "West US"
}

resource "azurerm_virtual_network" "example-1" {
  name                = "peternetwork1"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.1.0/24"]
  location            = "West US"
}

resource "azurerm_virtual_network" "example-2" {
  name                = "peternetwork2"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.2.0/24"]
  location            = "West US"
}

resource "azurerm_virtual_network_peering" "example-1" {
  name                      = "peer1to2"
  resource_group_name       = azurerm_resource_group.example.name
  virtual_network_name      = azurerm_virtual_network.example-1.name
  remote_virtual_network_id = azurerm_virtual_network.example-2.id
}

resource "azurerm_virtual_network_peering" "example-2" {
  name                      = "peer2to1"
  resource_group_name       = azurerm_resource_group.example.name
  virtual_network_name      = azurerm_virtual_network.example-2.name
  remote_virtual_network_id = azurerm_virtual_network.example-1.id
}
```

## Example Usage (Global virtual network peering)

```hcl
variable "location" {
  default = [
    "uksouth",
    "southeastasia",
  ]
}

variable "vnet_address_space" {
  default = [
    "10.0.0.0/16",
    "10.1.0.0/16",
  ]
}

resource "azurerm_resource_group" "vnet" {
  count    = length(var.location)
  name     = "rg-global-vnet-peering-${count.index}"
  location = element(var.location, count.index)
}

resource "azurerm_virtual_network" "vnet" {
  count               = length(var.location)
  name                = "vnet-${count.index}"
  resource_group_name = element(azurerm_resource_group.vnet.*.name, count.index)
  address_space       = [element(var.vnet_address_space, count.index)]
  location            = element(azurerm_resource_group.vnet.*.location, count.index)
}

resource "azurerm_subnet" "nva" {
  count                = length(var.location)
  name                 = "nva"
  resource_group_name  = element(azurerm_resource_group.vnet.*.name, count.index)
  virtual_network_name = element(azurerm_virtual_network.vnet.*.name, count.index)
  address_prefix = cidrsubnet(
    element(
      azurerm_virtual_network.vnet[count.index].address_space,
      count.index,
    ),
    13,
    0,
  ) # /29
}

# enable global peering between the two virtual network
resource "azurerm_virtual_network_peering" "peering" {
  count                        = length(var.location)
  name                         = "peering-to-${element(azurerm_virtual_network.vnet.*.name, 1 - count.index)}"
  resource_group_name          = element(azurerm_resource_group.vnet.*.name, count.index)
  virtual_network_name         = element(azurerm_virtual_network.vnet.*.name, count.index)
  remote_virtual_network_id    = element(azurerm_virtual_network.vnet.*.id, 1 - count.index)
  allow_virtual_network_access = true
  allow_forwarded_traffic      = true

  # `allow_gateway_transit` must be set to false for vnet Global Peering
  allow_gateway_transit = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual network peering. Changing this
    forces a new resource to be created.

* `virtual_network_name` - (Required) The name of the virtual network. Changing
    this forces a new resource to be created.

* `remote_virtual_network_id` - (Required) The full Azure resource ID of the
    remote virtual network.  Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the virtual network peering. Changing this forces a new resource to be
    created.

* `allow_virtual_network_access` - (Optional) Controls if the VMs in the remote
    virtual network can access VMs in the local virtual network. Defaults to
    true.

* `allow_forwarded_traffic` - (Optional) Controls if forwarded traffic from  VMs
    in the remote virtual network is allowed. Defaults to false.

* `allow_gateway_transit` - (Optional) Controls gatewayLinks can be used in the
    remote virtual network’s link to the local virtual network.

* `use_remote_gateways` - (Optional) Controls if remote gateways can be used on
    the local virtual network. If the flag is set to `true`, and
    `allow_gateway_transit` on the remote peering is also `true`, virtual network will
    use gateways of remote virtual network for transit. Only one peering can
    have this flag set to `true`. This flag cannot be set if virtual network
    already has a gateway. Defaults to `false`.

-> **NOTE:** `use_remote_gateways` must be set to `false` if using Global Virtual Network Peerings.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Network Peering.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Network Peering.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Network Peering.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network Peering.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Network Peering.

## Note

Virtual Network peerings cannot be created, updated or deleted concurrently.

## Import

Virtual Network Peerings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_network_peering.examplePeering /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/virtualNetworkPeerings/myvnet1peering
```
