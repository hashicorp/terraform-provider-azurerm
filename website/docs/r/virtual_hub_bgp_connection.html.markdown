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

resource "azurerm_virtual_hub" "example" {
  name                = "example-vhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Standard"
}

resource "azurerm_public_ip" "example" {
  name                = "example-pip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.5.1.0/24"
}

resource "azurerm_virtual_hub_ip" "example" {
  name                         = "example-vhubip"
  virtual_hub_id               = azurerm_virtual_hub.example.id
  private_ip_address           = "10.5.1.18"
  private_ip_allocation_method = "Static"
  public_ip_address_id         = azurerm_public_ip.example.id
  subnet_id                    = azurerm_subnet.example.id
}

resource "azurerm_virtual_hub_bgp_connection" "example" {
  name           = "example-vhub-bgpconnection"
  virtual_hub_id = azurerm_virtual_hub.example.id
  peer_asn       = 65514
  peer_ip        = "169.254.21.5"

  depends_on = [azurerm_virtual_hub_ip.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Virtual Hub Bgp Connection. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub within which this Bgp connection should be created. Changing this forces a new resource to be created.

* `peer_asn` - (Optional) The peer autonomous system number for the Virtual Hub Bgp Connection. Changing this forces a new resource to be created.

* `peer_ip` - (Optional) The peer ip address for the Virtual Hub Bgp Connection. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub Bgp Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Hub Bgp Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub Bgp Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Hub Bgp Connection.

## Import

Virtual Hub Bgp Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_hub_bgp_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/virtualHub1/bgpConnections/connection1
```
