---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_backend_address_pool"
description: |-
    Lists Load Balancer Backend Address Pool resources.
---

# List resource: azurerm_lb_backend_address_pool

Lists Load Balancer Backend Address Pool resources.

## Example Usage

### List all Load Balancer Backend Address Pools in a Load Balancer

```hcl
list "azurerm_lb_backend_address_pool" "example" {
  provider = azurerm
  config {
    loadbalancer_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/loadBalancers/example-lb"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `loadbalancer_id` - (Required) The ID of the Load Balancer to query.
