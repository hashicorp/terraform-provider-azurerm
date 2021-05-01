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

* `network_interface` - A list of `network_interface` blocks as defined below.

---

A `identity` block exports the following:

* `identity_ids` -  The list of User Managed Identity ID's which are assigned to the Virtual Machine Scale Set.

* `principal_id` - The ID of the System Managed Service Principal assigned to the Virtual Machine Scale Set.

* `type` - The identity type of the Managed Identity assigned to the Virtual Machine Scale Set.

---

`network_profile` exports the following:

* `name` - The name of the network interface configuration.
* `primary` - Whether network interfaces created from the network interface configuration will be the primary NIC of the VM.
* `ip_configuration` - An ip_configuration block as documented below.
* `accelerated_networking` - Whether to enable accelerated networking or not.
* `dns_settings` - A dns_settings block as documented below.
* `ip_forwarding` - Whether IP forwarding is enabled on this NIC.
* `network_security_group_id` - The identifier for the network security group.

`dns_settings` exports the following:

* `dns_servers` - The dns servers in use.

`ip_configuration` exports the following:

* `name` - The name of the IP configuration.
* `subnet_id` - The the identifier of the subnet.
* `application_gateway_backend_address_pool_ids` - An array of references to backend address pools of application gateways.
* `load_balancer_backend_address_pool_ids` - An array of references to backend address pools of load balancers.
* `load_balancer_inbound_nat_rules_ids` - An array of references to inbound NAT pools for load balancers.
* `primary` -  If this ip_configuration is the primary one.
* `application_security_group_ids` -  The application security group IDs to use.
* `public_ip_address_configuration` - The virtual machines scale set IP Configuration's PublicIPAddress configuration. The `public_ip_address_configuration` is documented below.

`public_ip_address_configuration` exports the following:

* `name` - The name of the public ip address configuration
* `idle_timeout` - The idle timeout in minutes.
* `domain_name_label` - The domain name label for the dns settings.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Scale Set.
