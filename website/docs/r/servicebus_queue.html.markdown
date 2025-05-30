---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_queue"
description: |-
  Manages a ServiceBus Queue.
---

# azurerm_servicebus_queue

Manages a ServiceBus Queue.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "terraform-servicebus"
  location = "West Europe"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "tfex-servicebus-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"

  tags = {
    source = "terraform"
  }
}

resource "azurerm_servicebus_queue" "example" {
  name         = "tfex_servicebus_queue"
  namespace_id = azurerm_servicebus_namespace.example.id

  partitioning_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Queue resource. Changing this forces a new resource to be created.

* `namespace_id` - (Required) The ID of the ServiceBus Namespace to create this queue in. Changing this forces a new resource to be created.

* `lock_duration` - (Optional) The ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. Maximum value is 5 minutes. Defaults to `PT1M` (1 Minute).

* `max_message_size_in_kilobytes` - (Optional) Integer value which controls the maximum size of a message allowed on the queue for Premium SKU. For supported values see the "Large messages support" section of [this document](https://docs.microsoft.com/azure/service-bus-messaging/service-bus-premium-messaging#large-messages-support-preview).

* `max_size_in_megabytes` - (Optional) Integer value which controls the size of memory allocated for the queue. For supported values see the "Queue or topic size" section of [Service Bus Quotas](https://docs.microsoft.com/azure/service-bus-messaging/service-bus-quotas).

* `requires_duplicate_detection` - (Optional) Boolean flag which controls whether the Queue requires duplicate detection. Changing this forces a new resource to be created. Defaults to `false`.

* `requires_session` - (Optional) Boolean flag which controls whether the Queue requires sessions. This will allow ordered handling of unbounded sequences of related messages. With sessions enabled a queue can guarantee first-in-first-out delivery of messages. Changing this forces a new resource to be created. Defaults to `false`.

* `default_message_ttl` - (Optional) The ISO 8601 timespan duration of the TTL of messages sent to this queue. This is the default value used when TTL is not set on message itself.

* `dead_lettering_on_message_expiration` - (Optional) Boolean flag which controls whether the Queue has dead letter support when a message expires. Defaults to `false`.

* `duplicate_detection_history_time_window` - (Optional) The ISO 8601 timespan duration during which duplicates can be detected. Defaults to `PT10M` (10 Minutes).

* `max_delivery_count` - (Optional) Integer value which controls when a message is automatically dead lettered. Defaults to `10`.

* `status` - (Optional) The status of the Queue. Possible values are `Active`, `Creating`, `Deleting`, `Disabled`, `ReceiveDisabled`, `Renaming`, `SendDisabled`, `Unknown`. Note that `Restoring` is not accepted. Defaults to `Active`.

* `batched_operations_enabled` - (Optional) Boolean flag which controls whether server-side batched operations are enabled. Defaults to `true`.

* `auto_delete_on_idle` - (Optional) The ISO 8601 timespan duration of the idle interval after which the Queue is automatically deleted, minimum of 5 minutes.

* `partitioning_enabled` - (Optional) Boolean flag which controls whether to enable the queue to be partitioned across multiple message brokers. Changing this forces a new resource to be created. Defaults to `false` for Basic and Standard.

-> **Note:** Partitioning is available at entity creation for all queues and topics in Basic or Standard SKUs. For premium namespace, partitioning is available at namespace creation, and all queues and topics in the partitioned namespace will be partitioned, for the premium namespace that has `premium_messaging_partitions` sets to `1`, the namespace is not partitioned.

* `express_enabled` - (Optional) Boolean flag which controls whether Express Entities are enabled. An express queue holds a message in memory temporarily before writing it to persistent storage. Defaults to `false` for Basic and Standard. For Premium, it MUST be set to `false`.

~> **Note:** Service Bus Premium namespaces do not support Express Entities, so `express_enabled` MUST be set to `false`.

* `forward_to` - (Optional) The name of a Queue or Topic to automatically forward messages to. Please [see the documentation](https://docs.microsoft.com/azure/service-bus-messaging/service-bus-auto-forwarding) for more information.

* `forward_dead_lettered_messages_to` - (Optional) The name of a Queue or Topic to automatically forward dead lettered messages to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ServiceBus Queue ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Queue.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Queue.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Queue.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Queue.

## Import

Service Bus Queue can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_queue.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceBus/namespaces/sbns1/queues/snqueue1
```
