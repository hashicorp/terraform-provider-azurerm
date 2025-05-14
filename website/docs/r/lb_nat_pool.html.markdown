---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_nat_pool"
description: |-
  Manages a Load Balancer NAT Pool.
---

# azurerm_lb_nat_pool

Manages a Load Balancer NAT pool.

-> **Note:** This resource cannot be used with with virtual machines, instead use the `azurerm_lb_nat_rule` resource.

~> **Note:** When using this resource, the Load Balancer needs to have a FrontEnd IP Configuration Attached

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

resource "azurerm_lb_nat_pool" "example" {
  resource_group_name            = azurerm_resource_group.example.name
  loadbalancer_id                = azurerm_lb.example.id
  name                           = "SampleApplicationPool"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 81
  backend_port                   = 8080
  frontend_ip_configuration_name = "PublicIPAddress"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the NAT pool. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the resource. Changing this forces a new resource to be created.
* `loadbalancer_id` - (Required) The ID of the Load Balancer in which to create the NAT pool. Changing this forces a new resource to be created.
* `frontend_ip_configuration_name` - (Required) The name of the frontend IP configuration exposing this rule.
* `protocol` - (Required) The transport protocol for the external endpoint. Possible values are `All`, `Tcp` and `Udp`.
* `frontend_port_start` - (Required) The first port number in the range of external ports that will be used to provide Inbound NAT to NICs associated with this Load Balancer. Possible values range between 1 and 65534, inclusive.
* `frontend_port_end` - (Required) The last port number in the range of external ports that will be used to provide Inbound NAT to NICs associated with this Load Balancer. Possible values range between 1 and 65534, inclusive.
* `backend_port` - (Required) The port used for the internal endpoint. Possible values range between 1 and 65535, inclusive.
* `idle_timeout_in_minutes` - (Optional) Specifies the idle timeout in minutes for TCP connections. Valid values are between `4` and `30`. Defaults to `4`.
* `floating_ip_enabled` - (Optional) Are the floating IPs enabled for this Load Balancer Rule? A floating IP is reassigned to a secondary server in case the primary server fails. Required to configure a SQL AlwaysOn Availability Group.
* `tcp_reset_enabled` - (Optional) Is TCP Reset enabled for this Load Balancer Rule? 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Load Balancer NAT pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Load Balancer NAT Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Load Balancer NAT Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Load Balancer NAT Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Load Balancer NAT Pool.

## Import

Load Balancer NAT Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb_nat_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/inboundNatPools/pool1
```
