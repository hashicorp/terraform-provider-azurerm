---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_servicebus_queue_authorization_rule"
description: |-
  Gets information about an existing ServiceBus Queue Authorisation Rule within a ServiceBus Queue.
---

# Data Source: azurerm_servicebus_queue_authorization_rule

Use this data source to access information about an existing ServiceBus Queue Authorisation Rule within a ServiceBus Queue.

## Example Usage

```hcl
data "azurerm_servicebus_queue_authorization_rule" "example" {
  name                = "example-tfex_name"
  resource_group_name = "example-resources"
  queue_name          = "example-servicebus_queue"
  namespace_name      = "example-namespace"
}

output "id" {
  value = data.azurerm_servicebus_queue_authorization_rule.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this ServiceBus Queue Authorisation Rule.

* `namespace_name` - (Required) The name of the ServiceBus Namespace.

* `queue_name` - (Required) The name of the ServiceBus Queue.

* `resource_group_name` - (Required) The name of the Resource Group where the ServiceBus Queue Authorisation Rule exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the ServiceBus Queue Authorisation Rule.

* `primary_key` - The Primary Key for the ServiceBus Queue authorization Rule.

* `primary_connection_string` - The Primary Connection String for the ServiceBus Queue authorization Rule.

* `secondary_key` - The Secondary Key for the ServiceBus Queue authorization Rule.

* `secondary_connection_string` - The Secondary Connection String for the ServiceBus Queue authorization Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Queue Authorisation Rule.
