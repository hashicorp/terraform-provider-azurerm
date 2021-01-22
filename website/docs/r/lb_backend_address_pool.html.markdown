---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_backend_address_pool"
description: |-
  Manages a Load Balancer Backend Address Pool.
---

# azurerm_lb_backend_address_pool

Manages a Load Balancer Backend Address Pool.

~> **NOTE:** When using this resource, the Load Balancer needs to have a FrontEnd IP Configuration Attached

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "LoadBalancerRG"
  location = "West US"
}

resource "azurerm_public_ip" "example" {
  name                = "PublicIPForLB"
  location            = "West US"
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "TestLoadBalancer"
  location            = "West US"
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_lb_backend_address_pool" "example" {
  loadbalancer_id = azurerm_lb.example.id
  name            = "BackEndAddressPool"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Backend Address Pool.
  
* `resource_group_name` - (Optional / **Deprecated**) The name of the resource group in which to create the resource.
  
* `loadbalancer_id` - (Required) The ID of the Load Balancer in which to create the Backend Address Pool.

* `backend_address` - (Optional) An array of `backend_address` block as defined below.

-> **NOTE**: The `backend_address` is only allowed when the Load Balancer in which this Backend Address Pool resides is of `Standard` sku. See [the Azure document](https://docs.microsoft.com/en-us/azure/load-balancer/backend-pool-management#limitations) for more details.

---

A `backend_address` supports the following:

* `name` - (Required) The name of the Backend Address.

* `virtual_network_id` - (Required) The ID of the Virtual Network that is pre-allocated for this Backend Address.

* `ip_address` - (Required) The IP address pre-allocated for this Backend Address with in the Virtual Network of `virtual_network_id`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Backend Address Pool.
  
* `backend_address` - An array of `backend_address` block as defined below.

* `backend_ip_configurations` - The Backend IP Configurations associated with this Backend Address Pool.

* `load_balancing_rules` - The Load Balancing Rules associated with this Backend Address Pool.

* `outbound_rules` - An array of the Load Balancing Outbound Rules associated with this Backend Address Pool.

---

A `backend_address` block exports the following:

* `network_interface_ip_configuration` - The Network Interface's ip configuration reference of this Backend Address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Load Balancer Backend Address Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Load Balancer Backend Address Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Load Balancer Backend Address Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Load Balancer Backend Address Pool.

## Import

Load Balancer Backend Address Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb_backend_address_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/backendAddressPools/pool1
```
