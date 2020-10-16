---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub_ip_configuration"
description: |-
  Manages a Virtual Hub IP Configuration.
---

# azurerm_virtual_hub_ip_configuration

Manages a Virtual Hub IP Configuration.

~> **NOTE** Virtual Hub IP Configuration only supports Standard Virtual Hub without Virtual Wan.

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

resource "azurerm_virtual_hub_ip_configuration" "example" {
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

* `name` - (Required) The name which should be used for this Virtual Hub IP Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Hub IP Configuration should exist. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub within which this ip configuration should be created. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet. Changing this forces a new resource to be created.

* `private_ip_address` - (Optional) The private IP address of the IP configuration.

* `private_ip_allocation_method` - (Optional) The private IP address allocation method. Possible values are `Static` and `Dynamic` is allowed. Defaults to `Dynamic`.

* `public_ip_address_id` - (Optional) The ID of the Public IP Address.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub IP Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Virtual Hub IP Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub IP Configuration.
* `update` - (Defaults to 60 minutes) Used when updating the Virtual Hub IP Configuration.
* `delete` - (Defaults to 60 minutes) Used when deleting the Virtual Hub IP Configuration.

## Import

Virtual Hub IP Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_hub_ip_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/virtualHub1/ipConfigurations/ipConfig1
```
