---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_profile"
sidebar_current: "docs-azurerm-resource-network-profile-x"
description: |-
  Manages an Azure Network Profile.

---

# azurerm_network_profile

Manages an Azure Network Profile.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "testgroup"
  location = "West Europe"
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.1.0.0/24"

  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_network_profile" "test" {
  name                                      = "testnetprofile"
  location                                  = "${azurerm_resource_group.test.location}"
  resource_group_name                       = "${azurerm_resource_group.test.name}"

  container_network_interface_configuration {
    name             = "testcnic"
    ip_configuration {
      name      = "testipconfig"
      subnet_id = "${azurerm_subnet.test.id}"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Network Profile. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the resource. Changing this forces a new resource to be created.

* `container_network_interface_configuration` - (Required) A `container_network_interface_configuration` block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `container_network_interface_configuration` block supports the following:

* `name` - (Required) Specifies the name of the IP Configuration.

* `ip_configuration` - (Required) A `ip_configuration` block as documented below.

---
A `ip_configuration` block supports the following:

* `name` - (Required) Specifies the name of the IP Configuration.

* `subnet_id` - (Required) Reference to the subnet associated with the IP Configuration.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the Azure Network Profile.

## Import

Azure Network Profile can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_profile.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/networkProfiles/testnetprofile
```
