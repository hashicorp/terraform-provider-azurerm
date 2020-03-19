---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_topic_authorization_rule"
description: |-
  Gets information about a ServiceBus Topic authorization Rule within a ServiceBus Topic.
---

# Data Source: azurerm_servicebus_topic_authorization_rule

Use this data source to access information about a ServiceBus Topic Authorization Rule within a ServiceBus Topic.

## Example Usage

```hcl
data "azurerm_servicebus_topic_authorization_rule" "example" {
  name                = "example-tfex_name"
  resource_group_name = "example-resources"
  namespace_name      = "example-namespace"
  topic_name          = "example-servicebus_topic"
}

output "servicebus_authorization_rule_id" {
  value = "${data.azurem_servicebus_topic_authorization_rule.example.id}"

}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the ServiceBus Topic Authorization Rule resource.

* `resource_group_name` - The name of the resource group in which the ServiceBus Namespace exists.

* `namespace_name` - The name of the ServiceBus Namespace.

* `topic_name` - The name of the ServiceBus Topic.

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Topic ID.

* `primary_key` - The Primary Key for the ServiceBus Topic authorization Rule.

* `primary_connection_string` - The Primary Connection String for the ServiceBus Topic authorization Rule.

* `secondary_key` - The Secondary Key for the ServiceBus Topic authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the ServiceBus Topic authorization Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Topic Authorization Rule.
