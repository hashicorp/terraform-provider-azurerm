---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_servicebus_topic"
description: |-
  Gets information about an existing Service Bus Topic.
---

# Data Source: azurerm_servicebus_topic

Use this data source to access information about an existing Service Bus Topic.

## Example Usage

```hcl
data "azurerm_servicebus_topic" "example" {
  name         = "existing"
  namespace_id = "existing"
}

output "id" {
  value = data.azurerm_servicebus_topic.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Service Bus Topic.

* `namespace_id` - (Required) The ID of the ServiceBus Namespace where the Service Bus Topic exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Bus Topic.

* `auto_delete_on_idle` - The ISO 8601 timespan duration of the idle interval after which the Topic is automatically deleted, minimum of 5 minutes.

* `default_message_ttl` - The ISO 8601 timespan duration of TTL of messages sent to this topic if no TTL value is set on the message itself.

* `duplicate_detection_history_time_window` - The ISO 8601 timespan duration during which duplicates can be detected.

* `batched_operations_enabled` - Boolean flag which controls if server-side batched operations are enabled.

* `express_enabled` - Boolean flag which controls whether Express Entities are enabled. An express topic holds a message in memory temporarily before writing it to persistent storage.

* `partitioning_enabled` - Boolean flag which controls whether to enable the topic to be partitioned across multiple message brokers.

* `max_size_in_megabytes` - Integer value which controls the size of memory allocated for the topic. For supported values see the "Queue/topic size" section of [this document](https://docs.microsoft.com/azure/service-bus-messaging/service-bus-quotas).

* `requires_duplicate_detection` - Boolean flag which controls whether the Topic requires duplicate detection.

* `status` - The Status of the Service Bus Topic. Acceptable values are Active or Disabled.

* `support_ordering` - Boolean flag which controls whether the Topic supports ordering.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Service Bus Topic.
