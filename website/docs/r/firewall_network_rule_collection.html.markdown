---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall_network_rule_collection"
sidebar_current: "docs-azurerm-resource-firewall-networkrulecollection"
description: |-
  Manages an Azure Firewall network rule collection.
---

# azurerm_firewall

Manages an Azure Firewall network rule collection.

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

resource "azurerm_firewall" "test" {
  name                = "testfirewall"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "configuration"
    subnet_id                     = "${azurerm_subnet.test.id}"
    internal_public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_firewall_network_rule_collection" "test" {
  name                = "testcollection"
  azure_firewall_name = "${azurerm_firewall.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  priority            = 100
  action              = "Allow"

  rule {
    name = "testrule"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      "8.8.8.8",
      "8.8.4.4",
    ]

    protocols = [
      "TCP",
      "UDP",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the rule collection.
* `azure_firewall_name` - (Required) Specifies the name of the Azure Firewall in which to create the collection.
* `resource_group_name` - (Required) Specifies the name of the resource group containing the Azure Firewall.
* `priority` - (Required) Specifies the priority of the rule collection.
* `action` - (Required) Specifies the action the rule will apply to matching traffic. Accepted values are `Allow` and `Deny`.
* `rule` - (Required) A rule block as described below. At least one rule must be configured.

`rule` supports the following:

* `name` - (Required) Specifies the name of the rule.
* `description` - (Optional) Specifies a description for the rule.
* `source_addresses` - (Required) A list of source IP addresses and/or IP ranges.
* `destination_addresses` - (Required) A list of destination IP addresses and/or IP ranges.
* `destination_ports` - (Required) A list of destination ports.
* `protocols` - (Required) A list of protocols. Accepted values are `Any`, `ICMP`, `TCP` or `UDP`.