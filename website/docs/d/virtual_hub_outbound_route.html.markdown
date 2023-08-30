---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_outbound_route"
description: |-
  Gets information about the outbound routes configured for the Virtual Hub on a particular connection. 
---

# Data Source: azurerm_virtual_hub_outbound_route

Uses this data source to access information about the outbound routes configured for the Virtual Hub on a particular connection.

## Virtual Hub outbound Route Usage

```hcl
data "azurerm_virtual_hub" "example" {
  name                = "example-hub"
  resource_group_name = "example-resources"
}

data "azurerm_virtual_hub_connection" "example" {
  name                = "example-connection"
  resource_group_name = data.azurerm_virtual_hub.example.resource_group_name
  virtual_hub_name    = data.azurerm_virtual_hub.example.name
}

data "azurerm_virtual_hub_outbound_route" "example" {
  virtual_hub_id     = data.azurerm_virtual_hub.example.id
  target_resource_id = data.azurerm_virtual_hub_connection.example.id
  connection_type    = "HubVirtualNetworkConnection"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_hub_id` - The ID of the Virtual Hub.

* `target_resource_id` - The ID of the connection resource whose outbound routes are being requested.
 
* `connection_type` - The type of the specified connection resource such as `ExpressRouteConnection`, `HubVirtualNetworkConnection`, `VpnConnection` and `P2SConnection`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub outbound Route.

* `route_map` - One or more `route_map` blocks as defined below.

---

An `route_map` block exports the following:

* `prefix` - The address prefix of the route.

* `bgp_communities` - BGP communities of the route. 

* `as_path` - The AS path of the route.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub route.
