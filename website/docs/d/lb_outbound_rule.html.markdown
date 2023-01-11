---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_lb_outbound_rule"
description: |-
  Gets information about an existing Load Balancer Outbound Rule.
---

# Data Source: azurerm_lb_outbound_rule

Use this data source to access information about an existing Load Balancer Outbound Rule.

## Example Usage

```hcl
data "azurerm_lb_outbound_rule" "example" {
  name            = "existing_lb_outbound_rule"
  loadbalancer_id = "existing_load_balancer_id"
}

output "id" {
  value = data.azurerm_lb_outbound_rule.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Load Balancer Outbound Rule.

* `loadbalancer_id` - (Required) The ID of the Load Balancer in which the Outbound Rule exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Load Balancer Outbound Rule.

* `allocated_outbound_ports` - The number of outbound ports used for NAT.

* `backend_address_pool_id` - The ID of the Backend Address Pool. Outbound traffic is randomly load balanced across IPs in the backend IPs.

* `frontend_ip_configuration` - A `frontend_ip_configuration` block as defined below.

* `idle_timeout_in_minutes` - The timeout for the TCP idle connection.

* `protocol` - The transport protocol for the external endpoint.

* `tcp_reset_enabled` - Is the bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination enabled? This value is useful when the protocol is set to TCP.

---

A `frontend_ip_configuration` block exports the following:

* `id` - The ID of the Frontend IP Configuration.

* `name` - The name of the Frontend IP Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Load Balancer Outbound Rule.
