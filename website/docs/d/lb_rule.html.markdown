---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_rule"
description: |-
  Gets information about an existing Load Balancer Rule.
---

# Data Source: azurerm_lb_rule

Use this data source to access information about an existing Load Balancer Rule.

## Example Usage

```hcl
data "azurerm_lb" "example" {
  name                = "example-lb"
  resource_group_name = "example-resources"
}

data "azurerm_lb_rule" "example" {
  name                = "first"
  resource_group_name = "example-resources"
  loadbalancer_id     = data.azurerm_lb.example.id
}

output "lb_rule_id" {
  value = data.azurerm_lb_rule.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `loadbalancer_id` - (Required) The ID of the Load Balancer Rule.

* `name` - (Required) The name of this Load Balancer Rule.

* `resource_group_name` - (Required) The name of the Resource Group where the Load Balancer Rule exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Load Balancer Rule.

* `backend_address_pool_id` - A reference to a Backend Address Pool over which this Load Balancing Rule operates.

* `probe_id` - A reference to a Probe used by this Load Balancing Rule.
 
* `frontend_ip_configuration_name` - The name of the frontend IP configuration to which the rule is associated.

* `protocol` - The transport protocol for the external endpoint.

* `frontend_port` - The port for the external endpoint.

* `backend_port` - The port used for internal connections on the endpoint.

* `enable_floating_ip` - If Floating IPs are enabled for this Load Balancer Rule

* `idle_timeout_in_minutes` - Specifies the idle timeout in minutes for TCP connections.

* `load_distribution` - Specifies the load balancing distribution type used by the Load Balancer. 

* `disable_outbound_snat` - If outbound SNAT is enabled for this Load Balancer Rule.

* `enable_tcp_reset` - If TCP Reset is enabled for this Load Balancer Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Load Balancer Rule.
