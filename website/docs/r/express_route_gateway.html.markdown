---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_gateway"
description: |-
  Manages an ExpressRoute gateway.
---

# azurerm_express_route_gateway

Manages an ExpressRoute gateway.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-virtualwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-virtualhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_express_route_gateway" "example" {
  name                = "expressRoute1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_hub_id      = azurerm_virtual_hub.example.id
  scale_units         = 1

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ExpressRoute gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the ExpressRoute gateway. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The ID of a Virtual HUB within which the ExpressRoute gateway should be created. Changing this forces a new resource to be created.

* `scale_units` - (Required) The number of scale units with which to provision the ExpressRoute gateway. Each scale unit is equal to 2Gbps, with support for up to 10 scale units (20Gbps).

* `allow_non_virtual_wan_traffic` - (Optional) Specified whether this gateway accept traffic from non-Virtual WAN networks. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the ExpressRoute gateway.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the ExpressRoute Gateway.

* `read` - (Defaults to 5 minutes) Used when retrieving the ExpressRoute Gateway.

* `update` - (Defaults to 90 minutes) Used when updating the ExpressRoute Gateway.

* `delete` - (Defaults to 90 minutes) Used when deleting the ExpressRoute Gateway.

## Import

ExpressRoute Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/expressRouteGateways/myExpressRouteGateway
```
