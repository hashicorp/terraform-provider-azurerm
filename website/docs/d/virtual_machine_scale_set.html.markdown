---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_virtual_machine_scale_set"
description: |-
  Gets information about an existing Virtual Machine Scale Set.
---

# Data Source: azurerm_virtual_machine_scale_set

Use this data source to access information about an existing Virtual Machine Scale Set.

## Example Usage

```hcl
data "azurerm_virtual_machine_scale_set" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_virtual_machine_scale_set.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Virtual Machine Scale Set.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Machine Scale Set exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Virtual Machine Scale Set.

* `identity` - A `identity` block as defined below.

* `network_profile` - A collection of network profile blocks as documented below.

---

A `identity` block exports the following:

* `identity_ids` -  The list of User Managed Identity ID's which are assigned to the Virtual Machine Scale Set.

* `principal_id` - The ID of the System Managed Service Principal assigned to the Virtual Machine Scale Set.

* `type` - The identity type of the Managed Identity assigned to the Virtual Machine Scale Set.

---

A `network_interface` block exports the list of network profile blocks that contain the following:

* `name` - The network interface configuration name.
* `ip_configuration` - An ip_configuration block as documented below.
* `dns_settings` - A dns_settings block as documented below.
* `enable_accelerated_networking` - Whether accelerated networking is enabled or not.
* `enable_ip_forwarding` - Whether accelerated networking is enabled or not.
* `network_security_group_id` - The network security group identifier.
* `primary` - Indicates whether the network interface is the primary one.
* `dns_servers` - A list of dns servers.

`ip_configuration` block exports the following:

* `name` - The IP configuration name.
* `subnet_id` - The subnet identifier.
* `version` - The Internet Protocol Version which is used for this IP Configuration.
* `application_gateway_backend_address_pool_ids` - A list of references to backend address pools of application gateways if configured.
* `application_security_group_ids` - A list of Application Security Group ID's which this Virtual Machine Scale Set is connected to.
* `load_balancer_backend_address_pool_ids` - A list of references to backend address pools of load balancers if configured.
* `load_balancer_inbound_nat_rules_ids` - A list of NAT Rule ID's from a Load Balancer which this Virtual Machine Scale Set is connected to.
* `primary` - Indicates whether this ip configuration is the primary one.
* `public_ip_address` - A `public_ip_address` block as documented below.

`public_ip_address` block exports the following:

* `name` - The name of the public ip address configuration
* `idle_timeout` - The idle timeout in minutes.
* `domain_name_label` - The domain name label for the dns settings.
* `ip_tag` - A `ip_tag` block as documented below.

`ip_tag` block exports the following:

* `tag` - The IP Tag associated with the Public IP.
* `type` - The Type of IP Tag.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Scale Set.
