---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bastion_host"
description: |-
  Manages a Bastion Host.

---

# azurerm_bastion_host

Manages a Bastion Host.

## Example Usage

This example deploys an Azure Bastion Host Instance to a target virtual network.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "examplevnet"
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["192.168.1.224/27"]
}

resource "azurerm_public_ip" "example" {
  name                = "examplepip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "example" {
  name                = "examplebastion"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.example.id
    public_ip_address_id = azurerm_public_ip.example.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Bastion Host. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Bastion Host.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.  Review [Azure Bastion Host FAQ](https://docs.microsoft.com/en-us/azure/bastion/bastion-faq) for supported locations.

* `ip_configuration` - (Required) A `ip_configuration` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `ip_configuration` block supports the following:

* `name` - (Required) The name of the IP configuration.

* `subnet_id` - (Required) Reference to a subnet in which this Bastion Host has been created.

* `public_ip_address_id` (Required)  Reference to a Public IP Address to associate with this Bastion Host.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Bastion Host.

* `dns_name` - The FQDN for the Bastion Host.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Bastion Host.
* `update` - (Defaults to 30 minutes) Used when updating the Bastion Host.
* `read` - (Defaults to 5 minutes) Used when retrieving the Bastion Host.
* `delete` - (Defaults to 30 minutes) Used when deleting the Bastion Host.

## Import

Bastion Hosts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bastion_host.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/bastionHosts/instance1
```
