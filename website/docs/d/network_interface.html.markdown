---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_interface"
sidebar_current: "docs-azurerm-datasource-network-interface-x"
description: |-
  Get information about the specified Network Interface.
---

# Data Source: azurerm_network_interface

Use this data source to access the properties of an Azure Network Interface.

## Example Usage

```hcl
data "azurerm_network_interface" "test" {
  name                 = "acctest-nic"
  resource_group_name  = "networking"
}

output "network_interface_id" {
  value = "${data.azurerm_network_interface.test.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Network Interface.
* `resource_group_name` - (Required) Specifies the name of the resource group the Network Interface is located in.

## Attributes Reference

* `id` - The ID of the virtual network.
* `applied_dns_servers` - List of applied DNS servers.
* `dns_servers` - The list of DNS servers used by this network interface.
* `mac_address` - The MAC address used by the network interface.
* `private_ip_address` - The private ip address associated to the network interface.
* `virtual_machine_id` - The ID of the virtual machine this network interface is attached to.
* `network_security_group_id` - The ID of the network security group associated to the network interface

