---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_backend_address_pool_address"
description: |-
  Manages a Backend Address within a Backend Address Pool.
---

# azurerm_lb_backend_address_pool_address

Manages a Backend Address within a Backend Address Pool.

-> **Note:** Backend Addresses can only be added to a `Standard` SKU Load Balancer.

## Example Usage

```hcl
data "azurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = "example-resources"
}

data "azurerm_lb" "example" {
  name                = "example-lb"
  resource_group_name = "example-resources"
}

data "azurerm_lb_backend_address_pool" "example" {
  name            = "first"
  loadbalancer_id = data.azurerm_lb.example.id
}

resource "azurerm_lb_backend_address_pool_address" "example" {
  name                    = "example"
  backend_address_pool_id = data.azurerm_lb_backend_address_pool.example.id
  virtual_network_id      = data.azurerm_virtual_network.example.id
  ip_address              = "10.0.0.1"
}

data "azurerm_lb_backend_address_pool" "backend-pool-cr" {
  name            = "globalLBBackendPool"
  loadbalancer_id = data.azurerm_lb.example.id
}

resource "azurerm_lb_backend_address_pool_address" "example-1" {
  name                                = "address1"
  backend_address_pool_id             = data.azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R1.frontend_ip_configuration[0].id
}

resource "azurerm_lb_backend_address_pool_address" "example-2" {
  name                                = "address2"
  backend_address_pool_id             = data.azurerm_lb_backend_address_pool.backend-pool-cr.id
  backend_address_ip_configuration_id = azurerm_lb.backend-lb-R2.frontend_ip_configuration[0].id
}
```

## Arguments Reference

-> **Note:** Backend Addresses can only be added to a `Standard` SKU Load Balancer. Cross region load balancer is for Load Balancer with `Global` SKU.

The following arguments are supported:

* `backend_address_pool_id` - (Required) The ID of the Backend Address Pool. Changing this forces a new Backend Address Pool Address to be created.

* `ip_address` - (Optional) The Static IP Address which should be allocated to this Backend Address Pool.

* `name` - (Required) The name which should be used for this Backend Address Pool Address. Changing this forces a new Backend Address Pool Address to be created.

* `virtual_network_id` - (Optional) The ID of the Virtual Network within which the Backend Address Pool should exist.

* `backend_address_ip_configuration_id` - (Optional) The ip config ID of the regional load balancer that's added to the global load balancer's backend address pool.

-> **Note:** For cross-region load balancer, please append the name of the load balancers, virtual machines, and other resources in each region with a -R1 and -R2.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backend Address Pool Address.

* `inbound_nat_rule_port_mapping` - A list of `inbound_nat_rule_port_mapping` block as defined below.

---

A `inbound_nat_rule_port_mapping` block exports the following:

* `inbound_nat_rule_name` - The name of the Load Balancing Inbound NAT Rules associated with this Backend Address Pool Address.

* `frontend_port` - The Frontend Port of the Load Balancing Inbound NAT Rules associated with this Backend Address Pool Address.

* `backend_port` - The Backend Port of the Load Balancing Inbound NAT Rules associated with this Backend Address Pool Address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backend Address Pool Address.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backend Address Pool Address.
* `update` - (Defaults to 30 minutes) Used when updating the Backend Address Pool Address.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backend Address Pool Address.

## Import

Backend Address Pool Addresses can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb_backend_address_pool_address.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/addresses/address1
```
