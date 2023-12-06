---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_port_authorization"
description: |-
  Manages an ExpressRoute Port Authorization.
---

# azurerm_express_route_port_authorization

Manages an ExpressRoute Port Authorization.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exprtTest"
  location = "West Europe"
}

resource "azurerm_express_route_port" "example" {
  name                = "port1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  peering_location    = "Airtel-Chennai-CLS"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_port_authorization" "example" {
  name                    = "exampleERCAuth"
  express_route_port_name = azurerm_express_route_port.example.name
  resource_group_name     = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ExpressRoute Port. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the ExpressRoute Port. Changing this forces a new resource to be created. 

* `express_route_port_name` - (Required) The name of the Express Route Port in which to create the Authorization. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the ExpressRoute Port Authorization.

* `authorization_key` - The Authorization Key.

* `authorization_use_status` - The authorization use status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ExpressRoute Port Authorization.
* `read` - (Defaults to 5 minutes) Used when retrieving the ExpressRoute Port Authorization.
* `delete` - (Defaults to 30 minutes) Used when deleting the ExpressRoute Port Authorization.

## Import

ExpressRoute Port Authorizations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_port_authorization.auth1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/expressRoutePorts/myExpressPort/authorizations/auth1
```
