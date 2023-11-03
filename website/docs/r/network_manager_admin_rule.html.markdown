---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_admin_rule"
description: |-
  Manages a Network Manager Admin Rule.
---

# azurerm_network_manager_admin_rule

Manages a Network Manager Admin Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin"]
  description    = "example network manager"
}

resource "azurerm_network_manager_network_group" "example" {
  name               = "example-network-group"
  network_manager_id = azurerm_network_manager.example.id
}

resource "azurerm_network_manager_security_admin_configuration" "example" {
  name               = "example-admin-conf"
  network_manager_id = azurerm_network_manager.example.id
}

resource "azurerm_network_manager_admin_rule_collection" "example" {
  name                            = "example-admin-rule-collection"
  security_admin_configuration_id = azurerm_network_manager_security_admin_configuration.example.id
  network_group_ids               = [azurerm_network_manager_network_group.example.id]
}

resource "azurerm_network_manager_admin_rule" "example" {
  name                     = "example-admin-rule"
  admin_rule_collection_id = azurerm_network_manager_admin_rule_collection.example.id
  action                   = "Deny"
  direction                = "Outbound"
  priority                 = 1
  protocol                 = "Tcp"
  source_port_ranges       = ["80", "1024-65535"]
  destination_port_ranges  = ["80"]
  source {
    address_prefix_type = "ServiceTag"
    address_prefix      = "Internet"
  }
  destination {
    address_prefix_type = "IPPrefix"
    address_prefix      = "10.1.0.1"
  }
  destination {
    address_prefix_type = "IPPrefix"
    address_prefix      = "10.0.0.0/24"
  }
  description = "example admin rule"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager Admin Rule. Changing this forces a new Network Manager Admin Rule to be created.

* `admin_rule_collection_id` - (Required) Specifies the ID of the Network Manager Admin Rule Collection. Changing this forces a new Network Manager Admin Rule to be created.

* `action` - (Required) Specifies the action allowed for this Network Manager Admin Rule. Possible values are `Allow`, `AlwaysAllow`, and `Deny`.

* `direction` - (Required) Indicates if the traffic matched against the rule in inbound or outbound. Possible values are `Inbound` and `Outbound`.

* `priority` - (Required) The priority of the rule. Possible values are integers between `1` and `4096`. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.

* `protocol` - (Required) Specifies which network protocol this Network Manager Admin Rule applies to. Possible values are `Ah`, `Any`, `Esp`, `Icmp`, `Tcp`, and `Udp`.

* `description` - (Optional) A description of the Network Manager Admin Rule.

* `destination_port_ranges` - (Optional) A list of string specifies the destination port ranges. Specify one or more single port number or port ranges such as `1024-65535`. Use `*` to specify any port.

* `destination` - (Optional) One or more `destination` blocks as defined below.

* `source_port_ranges` - (Optional) A list of string specifies the source port ranges. Specify one or more single port number or port ranges such as `1024-65535`. Use `*` to specify any port.

* `source` - (Optional) One or more `source` blocks as defined below.

---

A `destination` block supports the following:

* `address_prefix` - (Required) Specifies the address prefix. 

* `address_prefix_type` - (Required) Specifies the address prefix type. Possible values are `IPPrefix` and `ServiceTag`. For more information, please see [this document](https://learn.microsoft.com/en-us/azure/virtual-network-manager/concept-security-admins#source-and-destination-types).

---

A `source` block supports the following:

* `address_prefix` - (Required) Specifies the address prefix.

* `address_prefix_type` - (Required) Specifies the address prefix type. Possible values are `IPPrefix` and `ServiceTag`. For more information, please see [this document](https://learn.microsoft.com/en-us/azure/virtual-network-manager/concept-security-admins#source-and-destination-types).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Admin Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Admin Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Admin Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Admin Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Admin Rule.

## Import

Network Manager Admin Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_admin_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/securityAdminConfigurations/configuration1/ruleCollections/ruleCollection1/rules/rule1
```
