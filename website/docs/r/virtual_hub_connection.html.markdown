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
  address_space       = ["172.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_wan" "test" {
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

* `internet_security_enabled` - (Optional) Should Internet Security be enabled to secure internet traffic? Changing this forces a new resource to be created. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Virtual Hub Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub Connection.
* `delete` - (Defaults to 60 minutes) Used when deleting the Virtual Hub Connection.

## Import

Virtual Hub Connection's can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_virtual_hub_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/hub1/hubVirtualNetworkConnections/connection1
```
