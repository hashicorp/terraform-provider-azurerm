---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_namespace_authorization_rule"
description: |-
  Gets information about an existing Authorization Rule for an Event Hub Namespace.
---

# Data Source: azurerm_eventhub_namespace_authorization_rule

Use this data source to access information about an Authorization Rule for an Event Hub Namespace.

## Example Usage

```hcl
data "azurerm_eventhub_namespace_authorization_rule" "example" {
  name                = "navi"
  resource_group_name = "example-resources"
  namespace_name      = "example-ns"
}

output "eventhub_authorization_rule_id" {
  value = data.azurem_eventhub_namespace_authorization_rule.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the EventHub Authorization Rule resource. 

* `resource_group_name` - The name of the resource group in which the EventHub Namespace exists.

* `namespace_name` - Specifies the name of the EventHub Namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub ID.

* `primary_connection_string_alias` - The alias of the Primary Connection String for the Event Hubs authorization Rule.

* `secondary_connection_string_alias` - The alias of the Secondary Connection String for the Event Hubs authorization Rule.

* `listen` - Does this Authorization Rule have permissions to Listen to the Event Hub?

* `manage` - Does this Authorization Rule have permissions to Manage to the Event Hub?

* `primary_connection_string` - The Primary Connection String for the Event Hubs authorization Rule.

* `primary_key` - The Primary Key for the Event Hubs authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the Event Hubs authorization Rule.

* `secondary_key` - The Secondary Key for the Event Hubs authorization Rule.

* `send` - Does this Authorization Rule have permissions to Send to the Event Hub?

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Authorization Rule.
