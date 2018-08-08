---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_route"
sidebar_current: "docs-azurerm-resource-network-route-x"
description: |-
  Manages a Route within a Route Table.
---

# azurerm_route

Manages a Route within a Route Table.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_route_table" "example" {
  # ...
}

resource "azurerm_route" "example" {
  name                = "example-route"
  resource_group_name = "${azurerm_resource_group.example.name}"
  route_table_name    = "${azurerm_route_table.example.name}"
  address_prefix      = "10.1.0.0/16"
  next_hop_type       = "vnetlocal"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the route. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the route. Changing this forces a new resource to be created.

* `route_table_name` - (Required) The name of the route table within which create the route. Changing this forces a new resource to be created.

* `address_prefix` - (Required) The destination CIDR to which the route applies, such as `10.1.0.0/16`

* `next_hop_type` - (Required) The type of Azure hop the packet should be sent to. Possible values are `VirtualNetworkGateway`, `VnetLocal`, `Internet`, `VirtualAppliance` and `None`

* `next_hop_in_ip_address` - (Optional) Contains the IP address packets should be forwarded to. This field can only be set when `next_hop_type` is set to `VirtualAppliance`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Route.

## Import

Routes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_route.testRoute /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/routeTables/mytable1/routes/myroute1
```
