---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_queue"
sidebar_current: "docs-azurerm-resource-messaging-servicebus-queue-x"
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
  name                = "tfex_sevicebus_namespace"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "Standard"

  tags = {
    source = "terraform"
  }
}

resource "azurerm_servicebus_queue" "example" {
  name                = "tfex_servicebus_queue"
  resource_group_name = "${azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"

  enable_partitioning = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Queue resource. Changing this forces a
    new resource to be created.

* `namespace_name` - (Required) The name of the ServiceBus Namespace to create
    this queue in. Changing this forces a new resource to be created.

* `location` - (Optional / **Deprecated**) Specifies the supported Azure location where the resource exists.
    Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the namespace. Changing this forces a new resource to be created.

* `auto_delete_on_idle` - (Optional) The ISO 8601 timespan duration of the idle interval after which the
    Queue is automatically deleted, minimum of 5 minutes.

* `default_message_ttl` - (Optional) The ISO 8601 timespan duration of the TTL of messages sent to this
    queue. This is the default value used when TTL is not set on message itself.

* `duplicate_detection_history_time_window` - (Optional) The ISO 8601 timespan duration during which
    duplicates can be detected. Default value is 10 minutes. (`PT10M`)

* `enable_express` - (Optional) Boolean flag which controls whether Express Entities
    are enabled. An express queue holds a message in memory temporarily before writing
    it to persistent storage. Defaults to `false` for Basic and Standard. For Premium, it MUST
    be set to `false`.

~> **NOTE:** Service Bus Premium namespaces do not support Express Entities, so `enable_express` MUST be set to `false`.

* `enable_partitioning` - (Optional) Boolean flag which controls whether to enable
    the queue to be partitioned across multiple message brokers. Changing this forces
    a new resource to be created. Defaults to `false` for Basic and Standard. For Premium, it MUST
    be set to `true`.

-> **NOTE:** Partitioning is available at entity creation for all queues and topics in Basic or Standard SKUs. It is not available for the Premium messaging SKU, but any previously existing partitioned entities in Premium namespaces continue to work as expected. Please [see the documentation](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-partitioning) for more information.

* `lock_duration` - (Optional) The ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. Maximum value is 5 minutes. Defaults to 1 minute. (`PT1M`)

* `max_size_in_megabytes` - (Optional) Integer value which controls the size of
    memory allocated for the queue. For supported values see the "Queue/topic size"
    section of [this document](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quotas).

* `requires_duplicate_detection` - (Optional) Boolean flag which controls whether
    the Queue requires duplicate detection. Changing this forces
    a new resource to be created. Defaults to `false`.

* `requires_session` - (Optional) Boolean flag which controls whether the Queue requires sessions.
    This will allow ordered handling of unbounded sequences of related messages. With sessions enabled
    a queue can guarantee first-in-first-out delivery of messages.
    Changing this forces a new resource to be created. Defaults to `false`.

* `dead_lettering_on_message_expiration` - (Optional) Boolean flag which controls whether the Queue has dead letter support when a message expires. Defaults to `false`.

* `max_delivery_count` - (Optional) Integer value which controls when a message is automatically deadlettered. Defaults to `10`.

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Queue ID.

## Import

Service Bus Queue can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_queue.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.servicebus/namespaces/sbns1/queues/snqueue1
```
