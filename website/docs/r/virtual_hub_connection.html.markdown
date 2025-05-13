---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_connection"
description: |-
  Manages a Connection for a Virtual Hub.
---

# azurerm_virtual_hub_connection

Manages a Connection for a Virtual Hub.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["172.16.0.0/12"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-hub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_virtual_hub_connection" "example" {
  name                      = "example-vhub"
  virtual_hub_id            = azurerm_virtual_hub.example.id
  remote_virtual_network_id = azurerm_virtual_network.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this Connection, which must be unique within the Virtual Hub. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub within which this connection should be created. Changing this forces a new resource to be created.

* `remote_virtual_network_id` - (Required) The ID of the Virtual Network which the Virtual Hub should be connected to. Changing this forces a new resource to be created.

---

* `internet_security_enabled` - (Optional) Should Internet Security be enabled to secure internet traffic? Defaults to `false`.

* `routing` - (Optional) A `routing` block as defined below.

---

A `routing` block supports the following:

* `associated_route_table_id` - (Optional) The ID of the route table associated with this Virtual Hub connection.

* `inbound_route_map_id` - (Optional) The resource ID of the Route Map associated with this Routing Configuration for inbound learned routes.

* `outbound_route_map_id` - (Optional) The resource ID of the Route Map associated with this Routing Configuration for outbound advertised routes.

* `propagated_route_table` - (Optional) A `propagated_route_table` block as defined below.

* `static_vnet_local_route_override_criteria` - (Optional) The static VNet local route override criteria that is used to determine whether NVA in spoke VNet is bypassed for traffic with destination in spoke VNet. Possible values are `Contains` and `Equal`. Defaults to `Contains`. Changing this forces a new resource to be created.

* `static_vnet_propagate_static_routes_enabled` - (Optional) Whether the static routes should be propagated to the Virtual Hub. Defaults to `true`.

* `static_vnet_route` - (Optional) A `static_vnet_route` block as defined below.

---

A `propagated_route_table` block supports the following:

* `labels` - (Optional) The list of labels to assign to this route table.

* `route_table_ids` - (Optional) A list of Route Table IDs to associated with this Virtual Hub Connection.

---

A `static_vnet_route` block supports the following:

* `name` - (Optional) The name which should be used for this Static Route.

* `address_prefixes` - (Optional) A list of CIDR Ranges which should be used as Address Prefixes.

* `next_hop_ip_address` - (Optional) The IP Address which should be used for the Next Hop.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Hub Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Virtual Hub Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub Connection.
* `update` - (Defaults to 1 hour) Used when updating the Virtual Hub Connection.
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Hub Connection.

## Import

Virtual Hub Connection's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_hub_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/hub1/hubVirtualNetworkConnections/connection1
```
