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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Load Balancer. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the Load Balancer. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure Region where the Load Balancer should be created. Changing this forces a new resource to be created.

---

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Load Balancer should exist. Changing this forces a new Load Balancer to be created.

* `frontend_ip_configuration` - (Optional) One or more `frontend_ip_configuration` blocks as documented below.

* `sku` - (Optional) The SKU of the Azure Load Balancer. Accepted values are `Basic`, `Standard` and `Gateway`. Defaults to `Standard`. Changing this forces a new resource to be created.

-> **Note:** The `Microsoft.Network/AllowGatewayLoadBalancer` feature is required to be registered in order to use the `Gateway` SKU. The feature can only be registered by the Azure service team, please submit an [Azure support ticket](https://azure.microsoft.com/en-us/support/create-ticket/) for that.

* `sku_tier` - (Optional) `sku_tier` - (Optional) The SKU tier of this Load Balancer. Possible values are `Global` and `Regional`. Defaults to `Regional`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `frontend_ip_configuration` block supports the following:

* `name` - (Required) Specifies the name of the frontend IP configuration.

* `zones` - (Optional) Specifies a list of Availability Zones in which the IP Address for this Load Balancer should be located.

-> **Note:** Availability Zones are only supported with a [Standard SKU](https://docs.microsoft.com/azure/load-balancer/load-balancer-standard-availability-zones) and [in select regions](https://docs.microsoft.com/azure/availability-zones/az-overview) at this time.

* `subnet_id` - (Optional) The ID of the Subnet which should be associated with the IP Configuration.
* `gateway_load_balancer_frontend_ip_configuration_id` - (Optional) The Frontend IP Configuration ID of a Gateway SKU Load Balancer.
* `private_ip_address` - (Optional) Private IP Address to assign to the Load Balancer. The last one and first four IPs in any range are reserved and cannot be manually assigned.
* `private_ip_address_allocation` - (Optional) The allocation method for the Private IP Address used by this Load Balancer. Possible values as `Dynamic` and `Static`.
* `private_ip_address_version` - (Optional) The version of IP that the Private IP Address is. Possible values are `IPv4` or `IPv6`.
* `public_ip_address_id` - (Optional) The ID of a Public IP Address which should be associated with the Load Balancer.
* `public_ip_prefix_id` - (Optional) The ID of a Public IP Prefix which should be associated with the Load Balancer. Public IP Prefix can only be used with outbound rules.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Load Balancer ID.
* `frontend_ip_configuration` - A `frontend_ip_configuration` block as documented below.
* `private_ip_address` - The first private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.
* `private_ip_addresses` - The list of private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.

---

A `frontend_ip_configuration` block exports the following:

* `gateway_load_balancer_frontend_ip_configuration_id` - The id of the Frontend IP Configuration of a Gateway Load Balancer that this Load Balancer points to.
* `id` - The id of the Frontend IP Configuration.
* `inbound_nat_rules` - The list of IDs of inbound rules that use this frontend IP.
* `load_balancer_rules` - The list of IDs of load balancing rules that use this frontend IP.
* `outbound_rules` - The list of IDs outbound rules that use this frontend IP.
* `private_ip_address` - Private IP Address to assign to the Load Balancer.
* `private_ip_address_allocation` - The allocation method for the Private IP Address used by this Load Balancer. Possible values are `Dynamic` and `Static`.
* `public_ip_address_id` - The ID of a Public IP Address which is associated with this Load Balancer.
* `public_ip_prefix_id` - The ID of a Public IP Prefix which is associated with the Load Balancer.
* `subnet_id` - The ID of the Subnet which is associated with the IP Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Load Balancer.
* `read` - (Defaults to 5 minutes) Used when retrieving the Load Balancer.
* `update` - (Defaults to 30 minutes) Used when updating the Load Balancer.
* `delete` - (Defaults to 30 minutes) Used when deleting the Load Balancer.

## Import

Load Balancers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lb.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1
```
