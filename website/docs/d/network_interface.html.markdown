---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_interface"
sidebar_current: "docs-azurerm-datasource-network-interface"
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

* `applied_dns_servers` - List of DNS servers applied to the specified network interface.
* `dns_servers` - The list of DNS servers used by the specified network interface.
* `enable_accelerated_networking` - Indicates if accelerated networking is set on the specified network interface.
* `enable_ip_forwarding` - Indicate if IP forwarding is set on the specified network interface.
* `id` - The ID of the virtual network that the specified network interface is associated to.
* `internal_dns_name_label` - The internal dns name label of the specified network interface.
* `internal_fqdn` - The internal FQDN associated to the specified network interface.
* `ip_configuration` - The list of IP configurations associated to the specified network interface.
* `location` - The location of the specified network interface.
* `network_security_group_id` - The ID of the network security group associated to the specified network interface.
* `mac_address` - The MAC address used by the specified network interface.
* `private_ip_address` - The primary private ip address associated to the specified network interface.
* `private_ip_addresses` - The list of private ip addresses associates to the specified network interface.
* `tags` - List the tags assocatied to the specified network interface.
* `virtual_machine_id` - The ID of the virtual machine that the specified network interface is attached to.

