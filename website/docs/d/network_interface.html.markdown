---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_interface"
sidebar_current: "docs-azurerm-datasource-network-interface"
description: |-
  Gets information about an existing Network Interface.
---

# Data Source: azurerm_network_interface

Use this data source to access information about an existing Network Interface.

## Example Usage

```hcl
data "azurerm_network_interface" "example" {
  name                = "acctest-nic"
  resource_group_name = "networking"
}

output "network_interface_id" {
  value = "${data.azurerm_network_interface.example.id}"
}
```

## Argument Reference


* `name` - (Required) Specifies the name of the Network Interface.
* `resource_group_name` - (Required) Specifies the name of the resource group the Network Interface is located in.

## Attributes Reference

* `id` - The ID of the Network Interface.
* `applied_dns_servers` - List of DNS servers applied to the specified Network Interface.
* `enable_accelerated_networking` - Indicates if accelerated networking is set on the specified Network Interface.
* `enable_ip_forwarding` - Indicate if IP forwarding is set on the specified Network Interface.
* `dns_servers` - The list of DNS servers used by the specified Network Interface.
* `internal_dns_name_label` - The internal dns name label of the specified Network Interface.
* `ip_configuration` - One or more `ip_configuration` blocks as defined below.
* `location` - The location of the specified Network Interface.
* `mac_address` - The MAC address used by the specified Network Interface.
* `network_security_group_id` - The ID of the network security group associated to the specified Network Interface.
* `private_ip_address` - The primary private ip address associated to the specified Network Interface.
* `private_ip_addresses` - The list of private ip addresses associates to the specified Network Interface.
* `tags` - List the tags associated to the specified Network Interface.
* `virtual_machine_id` - The ID of the virtual machine that the specified Network Interface is attached to.

---

A `ip_configuration` block contains:

* `name` - The name of the IP Configuration.
* `subnet_id` - The ID of the Subnet which the Network Interface is connected to.
* `private_ip_address` - The Private IP Address assigned to this Network Interface.
* `private_ip_address_allocation` - The IP Address allocation type for the Private address, such as `Dynamic` or `Static`.
* `public_ip_address_id` - The ID of the Public IP Address which is connected to this Network Interface.
* `application_gateway_backend_address_pools_ids` - A list of Backend Address Pool ID's within a Application Gateway that this Network Interface is connected to.
* `load_balancer_backend_address_pools_ids` - A list of Backend Address Pool ID's within a Load Balancer that this Network Interface is connected to.
* `load_balancer_inbound_nat_rules_ids` - A list of Inbound NAT Rule ID's within a Load Balancer that this Network Interface is connected to.
* `primary` - is this the Primary IP Configuration for this Network Interface?
