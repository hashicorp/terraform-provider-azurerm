---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_rule"
description: |-
  Manages a Load Balancer Rule.
---

# azurerm_lb_rule

Manages a Load Balancer Rule.

~> **NOTE** When using this resource, the Load Balancer needs to have a FrontEnd IP Configuration Attached

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "LoadBalancerRG"
  location = "West Europe"
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

resource "azurerm_lb_rule" "example" {
  resource_group_name            = azurerm_resource_group.example.name
  loadbalancer_id                = azurerm_lb.example.id
  name                           = "LBRule"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "PublicIPAddress"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the LB Rule.
* `resource_group_name` - (Required) The name of the resource group in which to create the resource.
* `loadbalancer_id` - (Required) The ID of the Load Balancer in which to create the Rule.
* `frontend_ip_configuration_name` - (Required) The name of the frontend IP configuration to which the rule is associated.
* `protocol` - (Required) The transport protocol for the external endpoint. Possible values are `Tcp`, `Udp` or `All`.
* `frontend_port` - (Required) The port for the external endpoint. Port numbers for each Rule must be unique within the Load Balancer. Possible values range between 0 and 65534, inclusive.
* `backend_port` - (Required) The port used for internal connections on the endpoint. Possible values range between 0 and 65535, inclusive.
* `backend_address_pool_id` - (Optional) A reference to a Backend Address Pool over which this Load Balancing Rule operates.
* `probe_id` - (Optional) A reference to a Probe used by this Load Balancing Rule.
* `enable_floating_ip` - (Optional) Are the Floating IPs enabled for this Load Balncer Rule? A "floating” IP is reassigned to a secondary server in case the primary server fails. Required to configure a SQL AlwaysOn Availability Group. Defaults to `false`.
* `idle_timeout_in_minutes` - (Optional) Specifies the idle timeout in minutes for TCP connections. Valid values are between `4` and `30` minutes. Defaults to `4` minutes.
* `load_distribution` - (Optional) Specifies the load balancing distribution type to be used by the Load Balancer. Possible values are: `Default` – The load balancer is configured to use a 5 tuple hash to map traffic to available servers. `SourceIP` – The load balancer is configured to use a 2 tuple hash to map traffic to available servers. `SourceIPProtocol` – The load balancer is configured to use a 3 tuple hash to map traffic to available servers. Also known as Session Persistence, where  the options are called `None`, `Client IP` and `Client IP and Protocol` respectively.
* `disable_outbound_snat` - (Optional) Is snat enabled for this Load Balancer Rule? Default `false`.
* `enable_tcp_reset` - (Optional) Is TCP Reset enabled for this Load Balancer Rule? Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Load Balancer Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Load Balancer Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Load Balancer Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Load Balancer Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Load Balancer Rule.

## Import

Load Balancer Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/loadBalancingRules/rule1
```
