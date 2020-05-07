---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network"
description: |-
  Gets information about an existing Virtual Network.
---

# Data Source: azurerm_virtual_network

Use this data source to access information about an existing Virtual Network.

## Example Usage

```hcl
data "azurerm_virtual_network" "example" {
  name                = "production"
  resource_group_name = "networking"
}

output "virtual_network_id" {
  value = data.azurerm_virtual_network.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Virtual Network.
* `resource_group_name` - Specifies the name of the resource group the Virtual Network is located in.

## Attributes Reference

* `id` - The ID of the virtual network.
* `location` - Location of the virtual network.
* `address_space` - The list of address spaces used by the virtual network.
* `dns_servers` - The list of DNS servers used by the virtual network.
* `guid` - The GUID of the virtual network.
* `subnets` - The list of name of the subnets that are attached to this virtual network.
* `vnet_peerings` - A mapping of name - virtual network id of the virtual network peerings.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network.
