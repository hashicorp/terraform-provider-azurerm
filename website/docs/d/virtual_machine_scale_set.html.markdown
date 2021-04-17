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

A `network_profile` block exports the list of network profile blocks that contain the following:

* `name` - The network interface configuration name.
* `primary` - Indicates whether the network interface is the the primary.
* `ip_configuration` - An ip_configuration block as documented below.
* `accelerated_networking` - Whether accelerated networking is enabled or not.
* `dns_settings` - A dns_settings block as documented below.
* `ip_forwarding` - Whether IP forwarding is enabled or not on this NIC.
* `network_security_group_id` - The network security group identifier.

`dns_settings` block exports the following:

* `dns_servers` - An list of dns servers.

`ip_configuration` block exports the following:

* `name` - The IP configuration name.
* `subnet_id` - The subnet identifier.
* `application_gateway_backend_address_pool_ids` - A list of references to backend address pools of application gateways if configured.
* `load_balancer_backend_address_pool_ids` - A list of references to backend address pools of load balancers if configured.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Scale Set.
