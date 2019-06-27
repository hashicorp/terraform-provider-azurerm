---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall_application_rule_collection"
sidebar_current: "docs-azurerm-resource-network-firewall-application-rule-collection"
description: |-
  Manages an Application Rule Collection within an Azure Firewall.

---

# azurerm_firewall_application_rule_collection

Manages an Application Rule Collection within an Azure Firewall.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "North Europe"
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
  name                = "testpip"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "testfirewall"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = "${azurerm_subnet.test.id}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_firewall_application_rule_collection" "test" {
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

    target_fqdns = [
      "*.google.com",
    ]

    protocol {
      port = "443"
      type = "Https"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Rule Collection which must be unique within the Firewall. Changing this forces a new resource to be created.

* `azure_firewall_name` - (Required) Specifies the name of the Firewall in which the Application Rule Collection should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Firewall exists. Changing this forces a new resource to be created.

* `priority` - (Required) Specifies the priority of the rule collection. Possible values are between `100` - `65000`.

* `action` - (Required) Specifies the action the rule will apply to matching traffic. Possible values are `Allow` and `Deny`.

* `rule` - (Required) One or more `rule` blocks as defined below.

---

A `rule` block supports the following:

* `name` - (Required) Specifies the name of the rule.

* `description` - (Optional) Specifies a description for the rule.

* `source_addresses` - (Required) A list of source IP addresses and/or IP ranges.

* `fqdn_tags` - (Optional) A list of FQDN tags. Possible values are `AppServiceEnvironment`, `AzureBackup`, `MicrosoftActiveProtectionService`, `WindowsDiagnostics` and `WindowsUpdate`

* `target_fqdns` - (Optional) A list of FQDNs.

* `protocol` - (Optional) One or more `protocol` blocks as defined below.

---

A `protocol` block supports the following:

* `port` - (Optional) Specify a port for the connection.

* `type` - (Required) Specifies the type of connection. Possible values are `Http` or `Https`.

## Import

Azure Firewall Application Rule Collections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_firewall_application_rule_collection.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/applicationRuleCollections/mycollection
```
