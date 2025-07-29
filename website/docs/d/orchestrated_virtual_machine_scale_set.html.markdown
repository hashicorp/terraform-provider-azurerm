---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_orchestrated_virtual_machine_scale_set"
description: |-
  Gets information about an existing Orchestrated Virtual Machine Scale Set.
---

# Data Source: azurerm_orchestrated_virtual_machine_scale_set

Use this data source to access information about an existing Orchestrated Virtual Machine Scale Set.

## Example Usage

```hcl
data "azurerm_orchestrated_virtual_machine_scale_set" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_orchestrated_virtual_machine_scale_set.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Orchestrated Virtual Machine Scale Set.

* `resource_group_name` - (Required) The name of the Resource Group where the Orchestrated Virtual Machine Scale Set exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Machine Scale Set.

* `location` - The Azure Region in which this Orchestrated Virtual Machine Scale Set exists.

* `identity` - A `identity` block as defined below.

* `network_interface` - A list of `network_interface` blocks as defined below.

* `sku_profile` - A `sku_profile` block as defined below.

---

An `identity` block exports the following:

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Orchestrated Virtual Machine Scale Set.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Orchestrated Virtual Machine Scale Set.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Orchestrated Virtual Machine Scale Set.

* `type` - The type of Managed Service Identity that is configured on this Orchestrated Virtual Machine Scale Set.

---

A `network_interface` block exports the following:

* `accelerated_networking_enabled` - Is accelerated networking enabled?

* `dns_servers` - An array of the DNS servers in use.

* `ip_configuration` - An `ip_configuration` block as documented below.

* `ip_forwarding_enabled` - Is IP forwarding enabled?

* `name` - The name of the network interface configuration.

* `network_security_group_id` - The identifier for the network security group.

* `primary` - Whether network interfaces created from the network interface configuration will be the primary NIC of the VM.

---

An `ip_configuration` block exports the following:

* `application_gateway_backend_address_pool_ids` - An array of references to backend address pools of application gateways.

* `application_security_group_ids` -  The application security group IDs to use.

* `load_balancer_backend_address_pool_ids` - An array of references to backend address pools of load balancers.

* `load_balancer_inbound_nat_rules_ids` - An array of references to inbound NAT pools for load balancers.

* `name` - The name of the IP configuration.

* `primary` -  If this ip_configuration is the primary one.

* `public_ip_address` - The virtual machines scale set IP Configuration's PublicIPAddress configuration. The `public_ip_address` is documented below.

* `subnet_id` - The identifier of the subnet.

---

A `public_ip_address` block exports the following:

* `domain_name_label` - The domain name label for the DNS settings.

* `idle_timeout_in_minutes` - The idle timeout in minutes.

* `ip_tag` - A list of `ip_tag` blocks as defined below.

* `name` - The name of the public IP address configuration

* `public_ip_prefix_id` - The ID of the public IP prefix.

* `version` - The Internet Protocol Version of the public IP address.

---

An `ip_tag` block exports the following:

* `tag` - The IP Tag associated with the Public IP.

* `type` - The Type of IP Tag.

---

A `sku_profile` block exports the following:

* `allocation_strategy` - The allocation strategy used by this Orchestrated Virtual Machine Scale Set.

* `vm_sizes` - A list of `vm_sizes` blocks as defined below.

---

A `vm_sizes` block exports the following:

* `name` - The name of the VM size.

* `rank` - The rank of the VM size when `allocation_strategy` is set to `Prioritized`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Orchestrated Virtual Machine Scale Set.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Compute`: 2024-11-01
