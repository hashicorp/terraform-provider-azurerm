---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_backend_address_pool"
sidebar_current: "docs-azurerm-resource-loadbalancer-backend-address-pool"
description: |-
  Manages a LoadBalancer Backend Address Pool.
---

# azurerm_lb_backend_address_pool

Create a LoadBalancer Backend Address Pool.

~> **NOTE:** When using this resource, the LoadBalancer needs to have a FrontEnd IP Configuration Attached

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_public_ip" "example" {
  # ...
}

resource "azurerm_lb" "example" {
  # ...
}

resource "azurerm_lb_backend_address_pool" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  loadbalancer_id     = "${azurerm_lb.example.id}"
  name                = "BackEndAddressPool"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Backend Address Pool.
* `resource_group_name` - (Required) The name of the resource group in which to create the resource.
* `loadbalancer_id` - (Required) The ID of the LoadBalancer in which to create the Backend Address Pool.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the LoadBalancer to which the resource is attached.

## Import

Load Balancer Backend Address Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb_backend_address_pool.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/backendAddressPools/pool1
```
