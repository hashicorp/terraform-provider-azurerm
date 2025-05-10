---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_profile"
description: |-
  Manages a Network Profile.

---

# azurerm_network_profile

Manages a Network Profile.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "examplegroup"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "examplevnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "examplesubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.1.0.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_network_profile" "example" {
  name                = "examplenetprofile"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  container_network_interface {
    name = "examplecnic"

    ip_configuration {
      name      = "exampleipconfig"
      subnet_id = azurerm_subnet.example.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Network Profile. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the resource. Changing this forces a new resource to be created.

* `container_network_interface` - (Required) A `container_network_interface` block as documented below.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

A `container_network_interface` block supports the following:

* `name` - (Required) Specifies the name of the IP Configuration.

* `ip_configuration` - (Required) One or more `ip_configuration` blocks as documented below.

---

A `ip_configuration` block supports the following:

* `name` - (Required) Specifies the name of the IP Configuration.

* `subnet_id` - (Required) Reference to the subnet associated with the IP Configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Profile.

* `container_network_interface_ids` - A list of Container Network Interface IDs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Network Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Profile.

## Import

Network Profile can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkProfiles/examplenetprofile
```
