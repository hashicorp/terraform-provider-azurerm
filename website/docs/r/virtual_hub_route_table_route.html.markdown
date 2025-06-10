---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_route_table_route"
description: |-
  Manages a Route in a Virtual Hub Route Table.
---

# azurerm_virtual_hub_route_table_route

Manages a Route in a Virtual Hub Route Table.

~> **Note:** Route table routes can managed with this resource, or in-line with the [virtual_hub_route_table](virtual_hub_route_table.html) resource. Using both is not supported.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_security_group" "example" {
  name                = "example-nsg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "examplesubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-vhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.2.0/24"
}

resource "azurerm_virtual_hub_connection" "example" {
  name                      = "example-vhubconn"
  virtual_hub_id            = azurerm_virtual_hub.example.id
  remote_virtual_network_id = azurerm_virtual_network.example.id

  routing {
    associated_route_table_id = azurerm_virtual_hub_route_table.example.id
  }
}

resource "azurerm_virtual_hub_route_table" "example" {
  name           = "example-vhubroutetable"
  virtual_hub_id = azurerm_virtual_hub.example.id
  labels         = ["label1"]
}

resource "azurerm_virtual_hub_route_table_route" "example" {
  route_table_id = azurerm_virtual_hub_route_table.example.id

  name              = "example-route"
  destinations_type = "CIDR"
  destinations      = ["10.0.0.0/16"]
  next_hop_type     = "ResourceId"
  next_hop          = azurerm_virtual_hub_connection.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `route_table_id` - (Required) The ID of the Virtual Hub Route Table to link this route to. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this route. Changing this forces a new resource to be created.

* `destinations` - (Required) A list of destination addresses for this route.

* `destinations_type` - (Required) The type of destinations. Possible values are `CIDR`, `ResourceId` and `Service`.

* `next_hop` - (Required) The next hop's resource ID.

* `next_hop_type` - (Optional) The type of next hop. Currently the only possible value is `ResourceId`. Defaults to `ResourceId`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Hub Route Table.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Hub Route Table.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub Route Table.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Hub Route Table.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Hub Route Table.

## Import

Virtual Hub Route Table Routes can be imported using `<Route Table Resource Id>/routes/<Route Name>`, e.g.

```shell
terraform import azurerm_virtual_hub_route_table_route.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/routeTable1/routes/routeName
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network`: 2024-05-01
