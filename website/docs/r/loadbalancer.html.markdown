---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb"
sidebar_current: "docs-azurerm-resource-loadbalancer-x"
description: |-
  Manages a Load Balancer Resource.
---

# azurerm_lb

Manages a Load Balancer Resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "LoadBalancerRG"
  location = "West US"
}

resource "azurerm_public_ip" "test" {
  name                = "PublicIPForLB"
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "TestLoadBalancer"
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Load Balancer.
* `resource_group_name` - (Required) The name of the Resource Group in which to create the Load Balancer.
* `location` - (Required) Specifies the supported Azure Region where the Load Balancer should be created.
* `frontend_ip_configuration` - (Optional) One or multiple `frontend_ip_configuration` blocks as documented below.
* `sku` - (Optional) The SKU of the Azure Load Balancer. Accepted values are `Basic` and `Standard`. Defaults to `Basic`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

`frontend_ip_configuration` supports the following:

* `name` - (Required) Specifies the name of the frontend ip configuration.
* `subnet_id` - The ID of the Subnet which should be associated with the IP Configuration.
* `private_ip_address` - (Optional) Private IP Address to assign to the Load Balancer. The last one and first four IPs in any range are reserved and cannot be manually assigned.
* `private_ip_address_allocation` - (Optional) The allocation method for the Private IP Address used by this Load Balancer. Possible values as `Dynamic` and `Static`.
* `public_ip_address_id` - (Optional) The ID of a Public IP Address which should be associated with the Load Balancer.
* `public_ip_prefix_id` - (Optional) The ID of a Public IP Prefix which should be associated with the Load Balancer. Public IP Prefix can only be used with outbound rules.
* `zones` - (Optional) A list of Availability Zones which the Load Balancer's IP Addresses should be created in.

-> **Please Note**: Availability Zones are [only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview).

## Attributes Reference

The following attributes are exported:

* `id` - The Load Balancer ID.
* `private_ip_address` - The first private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.
* `private_ip_addresses` - The list of private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.
* `id` - The id of the Frontend IP Configuration.

## Import

Load Balancers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1
```
