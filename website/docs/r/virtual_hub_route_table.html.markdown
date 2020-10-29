---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_route_table"
description: |-
  Manages a Virtual Hub Route Table.
---

# azurerm_virtual_hub_route_table

Manages a Virtual Hub Route Table.

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
}

resource "azurerm_virtual_hub_route_table" "example" {
  name           = "example-vhubroutetable"
  virtual_hub_id = azurerm_virtual_hub.example.id
  labels         = ["label1"]

  route {
    name              = "example-route"
    destinations_type = "CIDR"
    destinations      = ["10.0.0.0/16"]
    next_hop_type     = "ResourceId"
    next_hop          = azurerm_virtual_hub_connection.example.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for Virtual Hub Route Table. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub within which this route table should be created. Changing this forces a new resource to be created.

* `labels` - (Optional) List of labels associated with this route table.

* `route` - (Optional)  A `route` block as defined below.

---

An `route` block exports the following:

* `name` - (Required) The name which should be used for this route.

* `destinations` - (Required) A list of destination addresses for this route.

* `destinations_type` - (Required) The type of destinations. Possible values are `CIDR`, `ResourceId` and `Service`.

* `next_hop` - (Required) The next hop's resource ID.

* `next_hop_type` - (Optional) The type of next hop. Currently the only possible value is `ResourceId`. Defaults to `ResourceId`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub Route Table.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Hub Route Table.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub Route Table.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Hub Route Table.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Hub Route Table.

## Import

Virtual Hub Route Tables can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_hub_route_table.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/routeTable1
```
