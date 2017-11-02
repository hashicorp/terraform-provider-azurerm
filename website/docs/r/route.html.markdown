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
resource "azurerm_resource_group" "test" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_route_table" "test" {
  name                = "acceptanceTestRouteTable1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_route" "test" {
  name                = "acceptanceTestRoute1"
  resource_group_name = "${azurerm_resource_group.test.name}"
  route_table_name    = "${azurerm_route_table.test.name}"
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

* `next_hop_in_ip_address` - (Optional) Contains the IP address packets should be forwarded to. Next hop values are only allowed in routes where the next hop type is `VirtualAppliance`.

## Attributes Reference

The following attributes are exported:

* `id` - The Route ID.

## Import

Routes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_route.testRoute /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/routeTables/mytable1/routes/myroute1
```
