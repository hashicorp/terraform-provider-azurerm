---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_circuit_authorization"
sidebar_current: "docs-azurerm-resource-network-express-route-circuit-authorization"
description: |-
  Manages an ExpressRoute Circuit Authorization.
---

# azurerm_express_route_circuit_authorization

Manages an ExpressRoute Circuit Authorization.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_express_route_circuit" "example" {
  # ...
}

resource "azurerm_express_route_circuit_authorization" "example" {
  name                       = "example-auth-rule"
  express_route_circuit_name = "${azurerm_express_route_circuit.example.name}"
  resource_group_name        = "${azurerm_resource_group.example.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ExpressRoute circuit. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the ExpressRoute circuit. Changing this forces a new resource to be created.

* `express_route_circuit_name` - (Required) The name of the Express Route Circuit in which to create the Authorization.


## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the ExpressRoute Circuit Authorization.

* `authorization_key` - The Authorization Key.

* `authorization_use_status` - The authorization use status.

## Import

ExpressRoute Circuit Authorizations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_circuit_authorization.auth1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/expressRouteCircuits/myExpressRoute/authorizations/auth1
```
