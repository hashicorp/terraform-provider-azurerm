---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_route_server_bgp_connection"
description: |-
  Manages a BGP Connection for a Route Server.
---

# azurerm_route_server_bgp_connection

Manages a Bgp Connection for a Route Server

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vn"
  address_space       = ["10.0.0.0/16"]
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  tags = {
    environment = "Production"
  }
}

resource "azurerm_subnet" "example" {
  name                 = "RouteServerSubnet"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "example-pip"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_route_server" "example" {
  name                             = "example-routerserver"
  resource_group_name              = azurerm_resource_group.example.name
  location                         = azurerm_resource_group.example.location
  sku                              = "Standard"
  public_ip_address_id             = azurerm_public_ip.example.id
  subnet_id                        = azurerm_subnet.example.id
  branch_to_branch_traffic_enabled = true
}

resource "azurerm_route_server_bgp_connection" "example" {
  name            = "example-rs-bgpconnection"
  route_server_id = azurerm_route_server.example.id
  peer_asn        = 65501
  peer_ip         = "169.254.21.5"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Route Server Bgp Connection. Changing this forces a new resource to be created.

* `route_server_id` - (Required) The ID of the Route Server within which this Bgp connection should be created. Changing this forces a new resource to be created.

* `peer_asn` - (Required) The peer autonomous system number for the Route Server Bgp Connection. Changing this forces a new resource to be created.

* `peer_ip` - (Required) The peer ip address for the Route Server Bgp Connection. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Route Server Bgp Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Route Server Bgp Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Route Server Bgp Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Route Server Bgp Connection.

## Import

Route Server Bgp Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_route_server_bgp_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/routeServer1/bgpConnections/connection1
```
