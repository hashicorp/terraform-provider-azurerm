---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_admin_rule_collection"
description: |-
  Manages a Network Manager Admin Rule Collection.
---

# azurerm_network_manager_admin_rule_collection

Manages a Network Manager Admin Rule Collection.

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager Admin Rule Collection. Changing this forces a new Network Manager Admin Rule Collection to be created.

* `security_admin_configuration_id` - (Required) Specifies the ID of the Network Manager Security Admin Configuration. Changing this forces a new Network Manager Admin Rule Collection to be created.

* `network_group_ids` - (Required) A list of Network Group ID which this Network Manager Admin Rule Collection applies to.

* `description` - (Optional) A description of the Network Manager Admin Rule Collection.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Admin Rule Collection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Admin Rule Collection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Admin Rule Collection.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Admin Rule Collection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Admin Rule Collection.

## Import

Network Manager Admin Rule Collection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_admin_rule_collection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/securityAdminConfigurations/configuration1/ruleCollections/ruleCollection1
```
