---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network"
sidebar_current: "docs-azurerm-datasource-virtual-network-x"
description: |-
  Gets information about an existing Virtual Network.
---

# Data Source: azurerm_virtual_network

Use this data source to access information about an existing Virtual Network.

## Example Usage

```hcl
data "azurerm_virtual_network" "test" {
  name                = "production"
  resource_group_name = "networking"
}

output "virtual_network_id" {
  value = "${data.azurerm_virtual_network.test.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Virtual Network.
* `resource_group_name` - (Required) Specifies the name of the resource group the Virtual Network is located in.

## Attributes Reference

* `id` - The ID of the virtual network.
* `location` - Location of the virtual network.
* `address_space` - The list of address spaces used by the virtual network.
* `dns_servers` - The list of DNS servers used by the virtual network.
* `subnets` - The list of name of the subnets that are attached to this virtual network.
* `vnet_peerings` - A mapping of name - virtual network id of the virtual network peerings.
