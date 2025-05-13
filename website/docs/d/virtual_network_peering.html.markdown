---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_virtual_network_peering"
description: |-
  Gets information about an existing virtual network peering.
---

# Data Source: azurerm_virtual_network_peering

Use this data source to access information about an existing virtual network peering.

## Example Usage

```hcl
data "azurerm_virtual_network" "example" {
  name                = "vnet01"
  resource_group_name = "networking"
}

data "azurerm_virtual_network_peering" "example" {
  name               = "peer-vnet01-to-vnet02"
  virtual_network_id = data.azurerm_virtual_network.example.id
}

output "id" {
  value = data.azurerm_virtual_network_peering.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this virtual network peering.

* `virtual_network_id` - (Required) The resource ID of the virtual network.  

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the virtual network peering.

* `allow_forwarded_traffic` - Controls if forwarded traffic from VMs in the remote virtual network is allowed.

* `allow_gateway_transit` - Controls gatewayLinks can be used in the remote virtual network’s link to the local virtual network.

* `allow_virtual_network_access` - Controls if the traffic from the local virtual network can reach the remote virtual network.

* `only_ipv6_peering_enabled` - Specifies whether only IPv6 address space is peered for Subnet peering.

* `peer_complete_virtual_networks_enabled` - Specifies whether complete Virtual Network address space is peered.

* `remote_virtual_network_id` - The full Azure resource ID of the remote virtual network.

* `use_remote_gateways` - Controls if remote gateways can be used on the local virtual network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the virtual network peering.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Network`: 2024-05-01
