---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_outbound_rule"
description: |-
  Manages a Load Balancer Outbound Rule.
---

# azurerm_lb_outbound_rule

Manages a Load Balancer Outbound Rule.

~> **Note:** When using this resource, the Load Balancer needs to have a FrontEnd IP Configuration and a Backend Address Pool Attached.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "LoadBalancerRG"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "PublicIPForLB"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "TestLoadBalancer"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_lb_backend_address_pool" "example" {
  name            = "example"
  loadbalancer_id = azurerm_lb.example.id
}

resource "azurerm_lb_outbound_rule" "example" {
  name                    = "OutboundRule"
  loadbalancer_id         = azurerm_lb.example.id
  protocol                = "Tcp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.example.id

  frontend_ip_configuration {
    name = "PublicIPAddress"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Outbound Rule. Changing this forces a new resource to be created.
* `loadbalancer_id` - (Required) The ID of the Load Balancer in which to create the Outbound Rule. Changing this forces a new resource to be created.
* `frontend_ip_configuration` - (Optional) One or more `frontend_ip_configuration` blocks as defined below.
* `backend_address_pool_id` - (Required) The ID of the Backend Address Pool. Outbound traffic is randomly load balanced across IPs in the backend IPs.
* `protocol` - (Required) The transport protocol for the external endpoint. Possible values are `Udp`, `Tcp` or `All`.
* `enable_tcp_reset` - (Optional) Receive bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination. This element is only used when the protocol is set to TCP.
* `allocated_outbound_ports` - (Optional) The number of outbound ports to be used for NAT. Defaults to `1024`.
* `idle_timeout_in_minutes` - (Optional) The timeout for the TCP idle connection Defaults to `4`.

---

A `frontend_ip_configuration` block supports the following:

* `name` - (Required) The name of the Frontend IP Configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Load Balancer Outbound Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Load Balancer Outbound Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Load Balancer Outbound Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Load Balancer Outbound Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Load Balancer Outbound Rule.

## Import

Load Balancer Outbound Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb_outbound_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/outboundRules/rule1
```
