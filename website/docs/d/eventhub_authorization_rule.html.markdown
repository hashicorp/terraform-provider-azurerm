---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_authorization_rule"
description: |-
  Gets information about an existing EventHubs authorization Rule within an EventHub.
---

# Data Source: azurerm_eventhub_authorization_rule

Use this data source to access information about an existing EventHub's authorization Rule within an EventHub.

## Example Usage

```hcl
data "azurerm_eventhub_namespace_authorization_rule" "example" {
  name                = "navi"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}" 
  resource_group_name = "${azurerm_eventhub_authorization_rule.example.resource_group_name}"

output "eventhub_authorization_rule_id" {
  value = "${data.azurem_eventhub_authorization_rule.example.id}"}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the EventHub Authorization Rule resource. 

* `resource_group_name` - (Required) The name of the resource group in which the EventHub Namespace exists.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub ID.

* `namespace_name` - The name of the grandparent EventHub Namespace. 

* `eventhub_name` - The name of the EventHub. 

* `listen` - Does this Authorization Rule have permissions to Listen to the Event Hub?

* `send` - Does this Authorization Rule have permissions to Send to the Event Hub?

* `manage` - Does this Authorization Rule have permissions to Manage to the Event Hub?

* `primary_key` - The Primary Key for the Event Hubs authorization Rule.

* `primary_connection_string` - The Primary Connection String for the Event Hubs authorization Rule.

* `secondary_key` - The Secondary Key for the Event Hubs authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the Event Hubs authorization Rule.

