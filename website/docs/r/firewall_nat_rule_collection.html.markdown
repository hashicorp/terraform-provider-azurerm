---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall_nat_rule_collection"
description: |-
  Manages a NAT Rule Collection within an Azure Firewall.

---

# azurerm_firewall_nat_rule_collection

Manages a NAT Rule Collection within an Azure Firewall.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "North Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "testvnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "testpip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "example" {
  name                = "testfirewall"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.example.id
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_firewall_nat_rule_collection" "example" {
  name                = "testcollection"
  azure_firewall_name = azurerm_firewall.example.name
  resource_group_name = azurerm_resource_group.example.name
  priority            = 100
  action              = "Dnat"

  rule {
    name = "testrule"

    source_addresses = [
      "10.0.0.0/16",
    ]

    destination_ports = [
      "53",
    ]

    destination_addresses = [
      azurerm_public_ip.example.ip_address
    ]

    translated_port = 53

    translated_address = "8.8.8.8"

    protocols = [
      "TCP",
      "UDP",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the NAT Rule Collection which must be unique within the Firewall. Changing this forces a new resource to be created.

* `azure_firewall_name` - (Required) Specifies the name of the Firewall in which the NAT Rule Collection should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Firewall exists. Changing this forces a new resource to be created.

* `priority` - (Required) Specifies the priority of the rule collection. Possible values are between `100` - `65000`.

* `action` - (Required) Specifies the action the rule will apply to matching traffic. Possible values are `Dnat` and `Snat`.

* `rule` - (Required) One or more `rule` blocks as defined below.

---

A `rule` block supports the following:

* `name` - (Required) Specifies the name of the rule.

* `description` - (Optional) Specifies a description for the rule.

* `destination_addresses` - (Required) A list of destination IP addresses and/or IP ranges.

* `destination_ports` - (Required) A list of destination ports.

* `protocols` - (Required) A list of protocols. Possible values are `Any`, `ICMP`, `TCP` and `UDP`.  If `action` is `Dnat`, protocols can only be `TCP` and `UDP`.

* `source_addresses` - (Optional) A list of source IP addresses and/or IP ranges.

* `source_ip_groups` - (Optional) A list of source IP Group IDs for the rule.

-> **NOTE** At least one of `source_addresses` and `source_ip_groups` must be specified for a rule.

* `translated_address` - (Required) The address of the service behind the Firewall.

* `translated_port` - (Required) The port of the service behind the Firewall.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Firewall NAT Rule Collection.
* `update` - (Defaults to 30 minutes) Used when updating the Firewall NAT Rule Collection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall NAT Rule Collection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Firewall NAT Rule Collection.

## Import

Azure Firewall NAT Rule Collections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_firewall_nat_rule_collection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/natRuleCollections/mycollection
```
