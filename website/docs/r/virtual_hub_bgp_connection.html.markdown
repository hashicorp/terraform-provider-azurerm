---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_bgp_connection"
description: |-
  Manages a Bgp Connection for a Virtual Hub.
---

# azurerm_virtual_hub_bgp_connection

Manages a Bgp Connection for a Virtual Hub.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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

resource "azurerm_virtual_hub_bgp_connection" "example" {
  name           = "example-vhubbgpconnection"
  virtual_hub_id = azurerm_virtual_hub.example.id
  peer_ip        = "192.168.1.5"
  pee_asn        = 20000
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Virtual Hub Bgp Connection. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub within which this Bgp connection should be created. Changing this forces a new resource to be created.

* `peer_asn` - (Optional) The peer autonomous system number.

* `peer_ip` - (Optional) The peer ip address.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network Virtual Hub Bgp Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the network Virtual Hub Bgp Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the network Virtual Hub Bgp Connection.
* `update` - (Defaults to 30 minutes) Used when updating the network Virtual Hub Bgp Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the network Virtual Hub Bgp Connection.

## Import

Virtual Hub Bgp Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_hub_bgp_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/virtualHub1/bgpConnections/connection1
```