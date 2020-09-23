---
subcategory: "Load Balancer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lb"
description: |-
  Manages a Load Balancer Resource.
---

# azurerm_lb

Manages a Load Balancer Resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "LoadBalancerRG"
  location = "West US"
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
* `private_ip_address_version` - The version of IP that the Private IP Address is. Possible values are `IPv4` or `IPv6`.
* `public_ip_address_id` - (Optional) The ID of a Public IP Address which should be associated with the Load Balancer.
* `public_ip_prefix_id` - (Optional) The ID of a Public IP Prefix which should be associated with the Load Balancer. Public IP Prefix can only be used with outbound rules.
* `zones` - (Optional) A list of Availability Zones which the Load Balancer's IP Addresses should be created in.

-> **Please Note**: Availability Zones are only supported with a [Standard SKU](https://docs.microsoft.com/en-us/azure/load-balancer/load-balancer-standard-availability-zones) and [in select regions](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview) at this time. Standard SKU Load Balancer that do not specify a zone are zone redundant by default.

## Attributes Reference

The following attributes are exported:

* `id` - The Load Balancer ID.
* `frontend_ip_configuration` - A `frontend_ip_configuration` block as documented below.
* `private_ip_address` - The first private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.
* `private_ip_addresses` - The list of private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.

---

A `frontend_ip_configuration` block exports the following:

* `id` - The id of the Frontend IP Configuration.
* `inbound_nat_rules` - The list of IDs of inbound rules that use this frontend IP.
* `load_balancer_rules` - The list of IDs of load balancing rules that use this frontend IP.
* `outbound_rules` - The list of IDs outbound rules that use this frontend IP.
* `private_ip_address` - Private IP Address to assign to the Load Balancer.
* `private_ip_address_allocation` - The allocation method for the Private IP Address used by this Load Balancer.
* `public_ip_address_id` - The ID of a  Public IP Address which is associated with this Load Balancer.
* `public_ip_prefix_id` - The ID of a Public IP Prefix which is associated with the Load Balancer.
* `subnet_id` - The ID of the Subnet which is associated with the IP Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Load Balancer.
* `update` - (Defaults to 30 minutes) Used when updating the Load Balancer.
* `read` - (Defaults to 5 minutes) Used when retrieving the Load Balancer.
* `delete` - (Defaults to 30 minutes) Used when deleting the Load Balancer.

## Import

Load Balancers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1
```
