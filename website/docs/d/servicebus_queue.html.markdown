---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_servicebus_queue"
description: |-
  Gets information about an existing Service Bus Queue.
---

# Data Source: azurerm_servicebus_queue

Use this data source to access information about an existing Service Bus Queue.

## Example Usage

```hcl
data "azurerm_servicebus_queue" "example" {
  name                = "existing"
  resource_group_name = "existing"
  namespace_name      = "existing"
}

output "id" {
  value = data.azurerm_servicebus_queue.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Service Bus Queue.

* `namespace_name` - (Required) The name of the ServiceBus Namespace.

* `resource_group_name` - (Required) The name of the Resource Group where the Service Bus Queue exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Service Bus Queue.

* `auto_delete_on_idle` - The ISO 8601 timespan duration of the idle interval after which the Queue is automatically deleted, minimum of 5 minutes.

* `dead_lettering_on_message_expiration` - Boolean flag which controls whether the Queue has dead letter support when a message expires.

* `default_message_ttl` - The ISO 8601 timespan duration of the TTL of messages sent to this queue. This is the default value used when TTL is not set on a message itself.

* `duplicate_detection_history_time_window` - The ISO 8601 timespan duration during which duplicates can be detected.

* `enable_batched_operations` - Boolean flag which controls whether server-side batched operations are enabled.

* `enable_express` - Boolean flag which controls whether Express Entities are enabled. An express queue holds a message in memory temporarily before writing it to persistent storage.

* `enable_partitioning` - Boolean flag which controls whether to enable the queue to be partitioned across multiple message brokers. 

* `forward_dead_lettered_messages_to` - The name of a Queue or Topic to automatically forward dead lettered messages to.

* `forward_to` - The name of a Queue or Topic to automatically forward messages to. Please [see the documentation](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-auto-forwarding) for more information.

* `lock_duration` - The ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers.

* `max_delivery_count` - Integer value which controls when a message is automatically dead lettered.

* `max_size_in_megabytes` - Integer value which controls the size of memory allocated for the queue. For supported values see the "Queue or topic size" section of [Service Bus Quotas](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quotas).

* `requires_duplicate_detection` - Boolean flag which controls whether the Queue requires duplicate detection.

* `requires_session` - Boolean flag which controls whether the Queue requires sessions. This will allow ordered handling of unbounded sequences of related messages. With sessions enabled a queue can guarantee first-in-first-out delivery of messages.

* `status` -  The status of the Queue. Possible values are `Active`, `Creating`, `Deleting`, `Disabled`, `ReceiveDisabled`, `Renaming`, `SendDisabled`, `Unknown`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Service Bus Queue.
