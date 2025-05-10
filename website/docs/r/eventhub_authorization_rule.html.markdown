---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_authorization_rule"
description: |-
  Manages a Event Hubs authorization Rule within an Event Hub.
---

# azurerm_eventhub_authorization_rule

Manages a Event Hubs authorization Rule within an Event Hub.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "acceptanceTestEventHubNamespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Basic"
  capacity            = 2

  tags = {
    environment = "Production"
  }
}

resource "azurerm_eventhub" "example" {
  name                = "acceptanceTestEventHub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 2
}

resource "azurerm_eventhub_authorization_rule" "example" {
  name                = "navi"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name
  listen              = true
  send                = false
  manage              = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Authorization Rule resource. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the name of the grandparent EventHub Namespace. Changing this forces a new resource to be created.

* `eventhub_name` - (Required) Specifies the name of the EventHub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventHub Namespace exists. Changing this forces a new resource to be created.

~> **Note:** At least one of the 3 permissions below needs to be set.

* `listen` - (Optional) Does this Authorization Rule have permissions to Listen to the Event Hub? Defaults to `false`.

* `send` - (Optional) Does this Authorization Rule have permissions to Send to the Event Hub? Defaults to `false`.

* `manage` - (Optional) Does this Authorization Rule have permissions to Manage to the Event Hub? When this property is `true` - both `listen` and `send` must be too. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The EventHub ID.

* `primary_connection_string_alias` - The alias of the Primary Connection String for the Event Hubs authorization Rule, which is generated when disaster recovery is enabled.

* `secondary_connection_string_alias` - The alias of the Secondary Connection String for the Event Hubs Authorization Rule, which is generated when disaster recovery is enabled.

* `primary_connection_string` - The Primary Connection String for the Event Hubs authorization Rule.

* `primary_key` - The Primary Key for the Event Hubs authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the Event Hubs Authorization Rule.

* `secondary_key` - The Secondary Key for the Event Hubs Authorization Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Authorization Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Authorization Rule.
* `update` - (Defaults to 30 minutes) Used when updating the EventHub Authorization Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventHub Authorization Rule.

## Import

EventHub Authorization Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_authorization_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/authorizationRules/rule1
```
