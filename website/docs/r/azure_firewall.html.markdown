---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azure_firewall"
sidebar_current: "docs-azurerm-resource-azurefirewall-x"
description: |-
  Manages an Azure Firewall.
---

# azurerm_azure_firewall

Manages an Azure Firewall.

~> **NOTE** This resource is currently in public preview.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "afwrg"
  location = "northeurope"
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "testpip"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Static"
  sku                          = "Standard"
}

resource "azurerm_azure_firewall" "test" {
  name                = "testfirewall"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "configuration"
    subnet_id                     = "${azurerm_subnet.test.id}"
    internal_public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the firewall.
* `resource_group_name` - (Required) The name of the resource group in which to create the resource.
* `location` - (Required) Specifies the supported Azure location where the resource exists.
* `ip_configuration` - (Required) An ip configuration block as documented below.
* `tags` - (Optional) A mapping of tags to assign to the resource.

`ip_configuration` supports the following:

* `name` - (Required) Specifies the name of the ip configuration.
* `subnet_id` - (Required) Reference to the subnet associated with the ip configuration.

~> **NOTE** The firewall subnet must be called `AzureFirewallSubnet` and the subnet mask must be at least `/25`

* `internal_public_ip_address_id` - (Required) Reference to the public IP address associated with the firewall.

~> **NOTE** The public IP must have a `Static` allocation and `Standard` sku