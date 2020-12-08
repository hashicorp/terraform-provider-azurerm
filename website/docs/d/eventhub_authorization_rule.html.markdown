---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_authorization_rule"
description: |-
  Gets information about an Event Hubs Authorization Rule within an Event Hub.
---

# azurerm_eventhub_authorization_rule

Use this data source to access information about an existing Event Hubs Authorization Rule within an Event Hub.

## Example Usage

```hcl
data "azurerm_eventhub_authorization_rule" "test" {
  name                = "test"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  eventhub_name       = "${azurerm_eventhub.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
```

## Argument Reference

* `name` - Specifies the name of the EventHub Authorization Rule resource. be created.

* `namespace_name` - Specifies the name of the grandparent EventHub Namespace.

* `eventhub_name` - Specifies the name of the EventHub.

* `resource_group_name` - The name of the resource group in which the EventHub Authorization Rule's grandparent Namespace exists.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub ID.

* `primary_connection_string_alias` - The alias of the Primary Connection String for the Event Hubs Authorization Rule.

* `secondary_connection_string_alias` - The alias of the Secondary Connection String for the Event Hubs Authorization Rule.

* `primary_connection_string` - The Primary Connection String for the Event Hubs Authorization Rule.

* `primary_key` - The Primary Key for the Event Hubs Authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the Event Hubs Authorization Rule.

* `secondary_key` - The Secondary Key for the Event Hubs Authorization Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Authorization Rule.
