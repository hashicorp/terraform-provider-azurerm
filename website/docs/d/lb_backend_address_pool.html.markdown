---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_backend_address_pool"
description: |-
  Get information about an existing Load Balancer Backend Address Pool

---

# Data Source: azurerm_lb_backend_address_pool

Use this data source to access information about an existing Load Balancer's Backend Address Pool.

## Example Usage

```hcl
data "azurerm_lb" "example" {
  name                = "example-lb"
  resource_group_name = "example-resources"
}

data "azurerm_lb_backend_address_pool" "example" {
  name            = "first"
  loadbalancer_id = data.azurerm_lb.example.id
}

output "backend_address_pool_id" {
  value = data.azurerm_lb_backend_address_pool.example.id
}

output "backend_ip_configuration_ids" {
  value = data.azurerm_lb_backend_address_pool.beap.backend_ip_configurations.*.id
}
```

## Argument Reference

* `name` - Specifies the name of the Backend Address Pool.

* `loadbalancer_id` - The ID of the Load Balancer in which the Backend Address Pool exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Backend Address Pool.

* `name` - The name of the Backend Address Pool.

* `backend_address` - A list of `backend_address` block as defined below.

* `backend_ip_configurations` - A list of references to IP addresses defined in network interfaces.

* `load_balancing_rules` - A list of the Load Balancing Rules associated with this Backend Address Pool.

* `inbound_nat_rules` - A list of the Load Balancing Inbound NAT Rules associated with this Backend Address Pool.

* `outbound_rules` - A list of the Load Balancing Outbound Rules associated with this Backend Address Pool.

---

A `backend_address` block exports the following:

* `name` - The name of the Backend Address.

* `virtual_network_id` - The ID of the Virtual Network where the Backend Address of the Load Balancer exists.

* `ip_address` - The Static IP address for this Load Balancer within the Virtual Network.

* `inbound_nat_rule_port_mapping` - A list of `inbound_nat_rule_port_mapping` block as defined below.

---

A `inbound_nat_rule_port_mapping` block exports the following:

* `inbound_nat_rule_name` - The name of the Load Balancing Inbound NAT Rules associated with this Backend Address Pool Address.

* `frontend_port` - The Frontend Port of the Load Balancing Inbound NAT Rules associated with this Backend Address Pool Address.

* `backend_port` - The Backend Port of the Load Balancing Inbound NAT Rules associated with this Backend Address Pool Address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Backend Address Pool.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Network`: 2023-09-01
