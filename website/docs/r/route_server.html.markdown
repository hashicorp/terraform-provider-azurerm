---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_route_server"
description: |-
  Manages an Azure Route Server 
---

# azurerm_route_server

Manages an Azure Route Server

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
  hub_routing_preference           = "ASPath"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Route Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Route Server should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Route Server should exist. Changing this forces a new resource to be created.

* `public_ip_address_id` - (Required) The ID of the Public IP Address. This option is required since September 1st 2021. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the Route Server. The only possible value is `Standard`. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet that the Route Server will reside. Changing this forces a new resource to be created.

-> **Note:** Azure Route Server requires a dedicated subnet named RouteServerSubnet. The subnet size has to be at least /27 or short prefix (such as /26 or /25) and cannot be attached to any security group, otherwise, you'll receive an error message when deploying the Route Server.

* `branch_to_branch_traffic_enabled` - (Optional) Whether to enable route exchange between Azure Route Server and the gateway(s).

* `hub_routing_preference` - (Optional) The hub routing preference. Valid values are `ASPath`, `ExpressRoute` or `VpnGateway`. Defaults to `ExpressRoute`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Route Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Route Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Route Server.
* `update` - (Defaults to 1 hour) Used when updating the Route Server.
* `delete` - (Defaults to 1 hour) Used when deleting the Route Server.

## Import

Route Server can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_route_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/routeServer1
```
