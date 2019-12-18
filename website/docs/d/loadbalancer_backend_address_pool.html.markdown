---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_backend_address_pool"
sidebar_current: "docs-azurerm-datasource-load-balancer-backend-address-pool"
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

* `name` - (Required) Specifies the name of the Backend Address Pool.

* `loadbalancer_id` - (Required) The ID of the Load Balancer in which the Backend Address Pool exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Backend Address Pool.

* `name` - The name of the Backend Address Pool.
 
* `backend_ip_configurations` - An array of references to IP addresses defined in network interfaces.
