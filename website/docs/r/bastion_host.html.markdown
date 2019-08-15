---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bastion_host"
sidebar_current: "docs-azurerm-resource-bastion-host-x"
description: |-
  Manages a Bastion Host Instance.

---

# azurerm_bastion_host

Manages a Bastion Host Instance.

## Example Usage

This example deploys an Azure Bastion Host Instance to a target virtual network.

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources-2"
  location = "West Europe"
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  address_space       = ["192.168.1.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "192.168.1.224/27"
}

resource "azurerm_public_ip" "test" {
  name                = "testpip"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                = "testbastion"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = "${azurerm_subnet.test.id}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the App Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the App Service.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `ip_configuration` - (Required) A `ip_configuration` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `ip_configuration` block supports the following:

* `name` - (Required) The Client ID of this relying party application. Enables OpenIDConnection authentication with Azure Active Directory.

* `subnet_id` - (Required) The Client Secret of this relying party application. If no secret is provided, implicit flow will be used.

* `public_ip_address_id` (Required) Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

## Import

App Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bastion_host.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1
```
