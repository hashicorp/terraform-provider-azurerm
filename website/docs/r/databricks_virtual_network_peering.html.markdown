---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_virtual_network_peering"
description: |-
  Manages a Databricks Virtual Network Peering
---

# azurerm_databricks_virtual_network_peering

Manages a Databricks Virtual Network Peering

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "remote" {
  name                = "remote-vnet"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.example.location
}

resource "azurerm_databricks_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "standard"
}

resource "azurerm_databricks_virtual_network_peering" "example" {
  name                = "databricks-vnet-peer"
  resource_group_name = azurerm_resource_group.example.name
  workspace_id        = azurerm_databricks_workspace.example.id

  remote_address_space_prefixes = azurerm_virtual_network.remote.address_space
  remote_virtual_network_id     = azurerm_virtual_network.remote.id
  allow_virtual_network_access  = true
}

resource "azurerm_virtual_network_peering" "remote" {
  name                         = "peer-to-databricks"
  resource_group_name          = azurerm_resource_group.example.name
  virtual_network_name         = azurerm_virtual_network.remote.name
  remote_virtual_network_id    = azurerm_databricks_virtual_network_peering.example.virtual_network_id
  allow_virtual_network_access = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Databricks Virtual Network Peering resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Databricks Virtual Network Peering should exist. Changing this forces a new resource to be created.

* `workspace_id` - (Required) The ID of the Databricks Workspace that this Databricks Virtual Network Peering is bound. Changing this forces a new resource to be created.

* `address_space_prefixes` - (Required) A list of address blocks reserved for this virtual network in CIDR notation.

* `remote_address_space_prefixes` - (Required) A list of address blocks reserved for the remote virtual network in CIDR notation. Changing this forces a new resource to be created.

* `remote_virtual_network_id` - (Required) The ID of the remote virtual network. Changing this forces a new resource to be created.

~> **Note:** The remote virtual network should be in the same region as the databricks workspace. Please see the [product documentation](https://learn.microsoft.com/azure/databricks/administration-guide/cloud-configurations/azure/vnet-peering) for more information.

* `allow_virtual_network_access` - (Optional) Can the VMs in the local virtual network space access the VMs in the remote virtual network space? Defaults to `true`.

* `allow_forwarded_traffic` - (Optional) Can the forwarded traffic from the VMs in the local virtual network be forwarded to the remote virtual network? Defaults to `false`.

* `allow_gateway_transit` - (Optional) Can the gateway links be used in the remote virtual network to link to the Databricks virtual network? Defaults to `false`.

* `use_remote_gateways` - (Optional) Can remote gateways be used on the Databricks virtual network? Defaults to `false`.

~> **Note:** If the `use_remote_gateways` is set to `true`, and `allow_gateway_transit` on the remote peering is also `true`, the virtual network will use the gateways of the remote virtual network for transit. Only one peering can have this flag set to `true`. `use_remote_gateways` cannot be set if the virtual network already has a gateway.

* `virtual_network_id` - The ID of the internal Virtual Network used by the DataBricks Workspace.

~> **Note:** The `virtual_network_id` field is the value you must supply to the `azurerm_virtual_network_peering` resources `remote_virtual_network_id` field to successfully peer the Databricks Virtual Network with the remote virtual network.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Databricks Virtual Network Peering.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Databricks Virtual Network Peering.
* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Virtual Network Peering.
* `update` - (Defaults to 30 minutes) Used when updating the Databricks Virtual Network Peering.
* `delete` - (Defaults to 30 minutes) Used when deleting the Databricks Virtual Network Peering.

## Import

Databrick Virtual Network Peerings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databricks_virtual_network_peering.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Databricks/workspaces/workspace1/virtualNetworkPeerings/peering1
```
