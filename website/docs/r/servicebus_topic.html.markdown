---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_topic"
description: |-
  Manages a ServiceBus Topic.
---

# azurerm_servicebus_topic

Manages a ServiceBus Topic.

~> **Note:** Topics can only be created in Namespaces with an SKU of `Standard` or higher.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-servicebus-topic"
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

resource "azurerm_servicebus_topic" "example" {
  name         = "tfex_servicebus_topic"
  namespace_id = azurerm_servicebus_namespace.example.id

  partitioning_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Topic resource. Changing this forces a new resource to be created.

* `namespace_id` - (Required) The ID of the ServiceBus Namespace to create this topic in. Changing this forces a new resource to be created.

* `status` - (Optional) The Status of the Service Bus Topic. Acceptable values are `Active` or `Disabled`. Defaults to `Active`.

* `auto_delete_on_idle` - (Optional) The ISO 8601 timespan duration of the idle interval after which the Topic is automatically deleted, minimum of 5 minutes. Defaults to `P10675199DT2H48M5.4775807S`.

* `default_message_ttl` - (Optional) The ISO 8601 timespan duration of TTL of messages sent to this topic if no TTL value is set on the message itself. Defaults to `P10675199DT2H48M5.4775807S`.

* `duplicate_detection_history_time_window` - (Optional) The ISO 8601 timespan duration during which duplicates can be detected. Defaults to `PT10M` (10 Minutes).

* `batched_operations_enabled` - (Optional) Boolean flag which controls if server-side batched operations are enabled.

* `express_enabled` - (Optional) Boolean flag which controls whether Express Entities are enabled. An express topic holds a message in memory temporarily before writing it to persistent storage.

* `partitioning_enabled` - (Optional) Boolean flag which controls whether to enable the topic to be partitioned across multiple message brokers. Changing this forces a new resource to be created.

-> **Note:** Partitioning is available at entity creation for all queues and topics in Basic or Standard SKUs. It is not available for the Premium messaging SKU, but any previously existing partitioned entities in Premium namespaces continue to work as expected. For premium namespaces, partitioning is available at namespace creation and all queues and topics in the partitioned namespace will be partitioned. Premium namespaces that have `premium_messaging_partitions` set to `1` are not partitioned. Please [see the documentation](https://docs.microsoft.com/azure/service-bus-messaging/service-bus-partitioning) for more information.

* `max_message_size_in_kilobytes` - (Optional) Integer value which controls the maximum size of a message allowed on the topic for Premium SKU. For supported values see the "Large messages support" section of [this document](https://docs.microsoft.com/azure/service-bus-messaging/service-bus-premium-messaging#large-messages-support-preview). Defaults to `256`.

* `max_size_in_megabytes` - (Optional) Integer value which controls the size of memory allocated for the topic. For supported values see the "Queue/topic size" section of [this document](https://docs.microsoft.com/azure/service-bus-messaging/service-bus-quotas). Defaults to `5120`.

* `requires_duplicate_detection` - (Optional) Boolean flag which controls whether the Topic requires duplicate detection. Defaults to `false`. Changing this forces a new resource to be created.

* `support_ordering` - (Optional) Boolean flag which controls whether the Topic supports ordering.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ServiceBus Topic ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Topic.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Topic.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Topic.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Topic.

## Import

Service Bus Topics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_topic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceBus/namespaces/sbns1/topics/sntopic1
```
