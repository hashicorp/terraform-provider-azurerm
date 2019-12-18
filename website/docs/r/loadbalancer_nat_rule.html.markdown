---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb_nat_rule"
sidebar_current: "docs-azurerm-resource-loadbalancer-nat-rule"
description: |-
  Manages a Load Balancer NAT Rule.
---

# azurerm_lb_nat_rule

Manages a Load Balancer NAT Rule. 

-> **NOTE:** This resource cannot be used with with virtual machine scale sets, instead use the `azurerm_lb_nat_pool` resource.

~> **NOTE** When using this resource, the Load Balancer needs to have a FrontEnd IP Configuration Attached

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "LoadBalancerRG"
  location = "West US"
}

resource "azurerm_public_ip" "example" {
  name                = "PublicIPForLB"
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.example.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "TestLoadBalancer"
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.example.name}"

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = "${azurerm_public_ip.example.id}"
  }
}

resource "azurerm_lb_nat_rule" "example" {
  resource_group_name            = "${azurerm_resource_group.example.name}"
  loadbalancer_id                = "${azurerm_lb.example.id}"
  name                           = "RDPAccess"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "PublicIPAddress"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the NAT Rule.
* `resource_group_name` - (Required) The name of the resource group in which to create the resource.
* `loadbalancer_id` - (Required) The ID of the Load Balancer in which to create the NAT Rule.
* `frontend_ip_configuration_name` - (Required) The name of the frontend IP configuration exposing this rule.
* `protocol` - (Required) The transport protocol for the external endpoint. Possible values are `Udp`, `Tcp` or `All`.
* `frontend_port` - (Required) The port for the external endpoint. Port numbers for each Rule must be unique within the Load Balancer. Possible values range between 1 and 65534, inclusive.
* `backend_port` - (Required) The port used for internal connections on the endpoint. Possible values range between 1 and 65535, inclusive.
* `enable_floating_ip` - (Optional) Enables the Floating IP Capacity, required to configure a SQL AlwaysOn Availability Group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Load Balancer NAT Rule.

## Import

Load Balancer NAT Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb_nat_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/inboundNatRules/rule1
```
