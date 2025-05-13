---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_ip"
description: |-
  Manages a Virtual Hub IP. This resource is also known as a Route Server.
---

# azurerm_virtual_hub_ip

Manages a Virtual Hub IP. This resource is also known as a Route Server.

~> **Note:** Virtual Hub IP only supports Standard Virtual Hub without Virtual Wan.

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
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "RouteServerSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_virtual_hub_ip" "example" {
  name                         = "example-vhubipconfig"
  virtual_hub_id               = azurerm_virtual_hub.example.id
  private_ip_address           = "10.5.1.18"
  private_ip_allocation_method = "Static"
  public_ip_address_id         = azurerm_public_ip.example.id
  subnet_id                    = azurerm_subnet.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Virtual Hub IP. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub within which this IP configuration should be created. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet that the IP will reside. Changing this forces a new resource to be created.

* `private_ip_address` - (Optional) The private IP address of the IP configuration.

* `private_ip_allocation_method` - (Optional) The private IP address allocation method. Possible values are `Static` and `Dynamic` is allowed. Defaults to `Dynamic`.

* `public_ip_address_id` - (Required) The ID of the Public IP Address. This option is required since September 1st 2021. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Hub IP.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Virtual Hub IP.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub IP.
* `update` - (Defaults to 1 hour) Used when updating the Virtual Hub IP.
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Hub IP.

## Import

Virtual Hub IPs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_hub_ip.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/virtualHub1/ipConfigurations/ipConfig1
```
