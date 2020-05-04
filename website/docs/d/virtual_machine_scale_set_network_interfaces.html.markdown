---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_scale_set_network_interfaces"
description: |-
  Gets information about the network interfaces of a virtual machine scale set.
---

# Data Source: azurerm_virtual_machine_scale_set_network_interfaces

Use this data source to access information about the network interfaces of a virtual machine scale set.

## Example Usage

```hcl
data "azurerm_virtual_machine_scale_set_network_interfaces" "example" {
  virtual_machine_scale_set_name = "scaleset"
  resource_group_name            = "networking"
}

output "network_interfaces" {
  value = data.azurerm_virtual_machine_scale_set_network_interfaces.network_interfaces
}
```

## Argument Reference


* `virtual_machine_scale_set_name` - Specifies the name of the virtual machine scale set.
* `resource_group_name` - Specifies the name of the resource group the virtual machine scale set is located in.

## Attributes Reference

* `network_interfaces` - A list of `network_interfaces` blocks associated with the virtual machine scale set.

A `network_interfaces` block contains:

* `applied_dns_servers` - List of DNS servers applied to the network interface.
* `enable_accelerated_networking` - Indicates if accelerated networking is set on the network interface.
* `enable_ip_forwarding` - Indicate if IP forwarding is set on the network interface.
* `dns_servers` - The list of DNS servers used by the network interface.
* `internal_dns_name_label` - The internal dns name label of the network interface.
* `ip_configuration` - One or more `ip_configuration` blocks as defined below.
* `mac_address` - The MAC address used by the network interface.
* `name` - The name of the network interface.
* `network_security_group_id` - The ID of the network security group associated to the network interface.
* `private_ip_address` - The primary private ip address associated to the network interface.
* `private_ip_addresses` - The list of private ip addresses associates to the network interface.
* `virtual_machine_id` - The ID of the virtual machine that the network interface is attached to.

---

An `ip_configuration` block contains:

* `name` - The name of the IP Configuration.
* `subnet_id` - The ID of the Subnet which the network interface is connected to.
* `private_ip_address` - The Private IP Address assigned to this network interface.
* `private_ip_address_allocation` - The IP Address allocation type for the Private address, such as `Dynamic` or `Static`.
* `public_ip_address_id` - The ID of the Public IP Address which is connected to this network interface.
* `application_gateway_backend_address_pools_ids` - A list of Backend Address Pool ID's within a Application Gateway that this network interface is connected to.
* `load_balancer_backend_address_pools_ids` - A list of Backend Address Pool ID's within a Load Balancer that this network interface is connected to.
* `load_balancer_inbound_nat_rules_ids` - A list of Inbound NAT Rule ID's within a Load Balancer that this network interface is connected to.
* `primary` - is this the Primary IP Configuration for this network interface?

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the network interfaces.
