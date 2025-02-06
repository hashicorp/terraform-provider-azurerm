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

* `location` - The Azure Region in which this Virtual Machine Scale Set exists.

* `identity` - A `identity` block as defined below.

* `instances` - A list of `instances` blocks as defined below.

* `network_interface` - A list of `network_interface` blocks as defined below.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Virtual Machine Scale Set.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Virtual Machine Scale Set.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Virtual Machine Scale Set.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Virtual Machine Scale Set.

---

`instances` exports the following:

* `computer_name` - The Hostname of this Virtual Machine.
* `instance_id` - The Instance ID of this Virtual Machine.
* `latest_model_applied` - Whether the latest model has been applied to this Virtual Machine.
* `name` - The name of the this Virtual Machine.
* `private_ip_address` - The Primary Private IP Address assigned to this Virtual Machine.
* `private_ip_addresses` - A list of Private IP Addresses assigned to this Virtual Machine.
* `public_ip_address` - The Primary Public IP Address assigned to this Virtual Machine.
* `public_ip_addresses` - A list of the Public IP Addresses assigned to this Virtual Machine.
* `power_state` - The power state of the virtual machine.
* `virtual_machine_id` - The unique ID of the virtual machine.
* `zone` - The zones of the virtual machine.

---

`network_interface` exports the following:

* `name` - The name of the network interface configuration.
* `primary` - Whether network interfaces created from the network interface configuration will be the primary NIC of the VM.
* `ip_configuration` - An `ip_configuration` block as documented below.
* `enable_accelerated_networking` - Whether to enable accelerated networking or not.
* `dns_servers` - An array of the DNS servers in use.
* `enable_ip_forwarding` - Whether IP forwarding is enabled on this NIC.
* `network_security_group_id` - The identifier for the network security group.

`ip_configuration` exports the following:

* `name` - The name of the IP configuration.
* `subnet_id` - The identifier of the subnet.
* `application_gateway_backend_address_pool_ids` - An array of references to backend address pools of application gateways.
* `load_balancer_backend_address_pool_ids` - An array of references to backend address pools of load balancers.
* `load_balancer_inbound_nat_rules_ids` - An array of references to inbound NAT pools for load balancers.
* `primary` -  If this ip_configuration is the primary one.
* `application_security_group_ids` -  The application security group IDs to use.
* `public_ip_address` - The virtual machines scale set IP Configuration's PublicIPAddress configuration. The `public_ip_address` is documented below.

`public_ip_address` exports the following:

* `name` - The name of the public IP address configuration
* `idle_timeout_in_minutes` - The idle timeout in minutes.
* `domain_name_label` - The domain name label for the DNS settings.
* `ip_tag` - A list of `ip_tag` blocks as defined below.
* `public_ip_prefix_id` - The ID of the public IP prefix.
* `version` - The Internet Protocol Version of the public IP address.

`ip_tag` exports the following:

* `tag` - The IP Tag associated with the Public IP.
* `type` - The Type of IP Tag.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Scale Set.
