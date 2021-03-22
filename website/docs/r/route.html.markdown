---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_route"
description: |-
  Manages a Route within a Route Table.
---

# azurerm_route

Manages a Route within a Route Table.

~> **NOTE on Route Tables and Routes:** Terraform currently
provides both a standalone [Route resource](route.html), and allows for Routes to be defined in-line within the [Route Table resource](route_table.html).
At this time you cannot use a Route Table with in-line Routes in conjunction with any Route resources. Doing so will cause a conflict of Route configurations and will overwrite Routes.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_route_table" "example" {
  name                = "acceptanceTestRouteTable1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_route" "example" {
  name                = "acceptanceTestRoute1"
  resource_group_name = azurerm_resource_group.example.name
  route_table_name    = azurerm_route_table.example.name
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Route.
* `update` - (Defaults to 30 minutes) Used when updating the Route.
* `read` - (Defaults to 5 minutes) Used when retrieving the Route.
* `delete` - (Defaults to 30 minutes) Used when deleting the Route.

## Import

Routes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_route.exampleRoute /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/routeTables/mytable1/routes/myroute1
```
