---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_topic"
description: |-
  Manages a ServiceBus Topic.
---

# azurerm_servicebus_topic

Manages a ServiceBus Topic.

**Note** Topics can only be created in Namespaces with an SKU of `standard` or higher.

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
  name                = "tfex_servicebus_topic"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_servicebus_namespace.example.name

  enable_partitioning = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Topic resource. Changing this forces a
    new resource to be created.

* `namespace_name` - (Required) The name of the ServiceBus Namespace to create
    this topic in. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the namespace. Changing this forces a new resource to be created.

* `status` - (Optional) The Status of the Service Bus Topic. Acceptable values are `Active` or `Disabled`. Defaults to `Active`.

* `auto_delete_on_idle` - (Optional) The ISO 8601 timespan duration of the idle interval after which the
    Topic is automatically deleted, minimum of 5 minutes.

* `default_message_ttl` - (Optional) The ISO 8601 timespan duration of TTL of messages sent to this topic if no
    TTL value is set on the message itself.

* `duplicate_detection_history_time_window` - (Optional) The ISO 8601 timespan duration during which
    duplicates can be detected. Defaults to 10 minutes. (`PT10M`)

* `enable_batched_operations` - (Optional) Boolean flag which controls if server-side
    batched operations are enabled. Defaults to false.

* `enable_express` - (Optional) Boolean flag which controls whether Express Entities
    are enabled. An express topic holds a message in memory temporarily before writing
    it to persistent storage. Defaults to false.

* `enable_partitioning` - (Optional) Boolean flag which controls whether to enable
    the topic to be partitioned across multiple message brokers. Defaults to false.
    Changing this forces a new resource to be created.

-> **NOTE:** Partitioning is available at entity creation for all queues and topics in Basic or Standard SKUs. It is not available for the Premium messaging SKU, but any previously existing partitioned entities in Premium namespaces continue to work as expected. Please [see the documentation](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-partitioning) for more information.

* `max_size_in_megabytes` - (Optional) Integer value which controls the size of
    memory allocated for the topic. For supported values see the "Queue/topic size"
    section of [this document](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quotas).

* `requires_duplicate_detection` - (Optional) Boolean flag which controls whether
    the Topic requires duplicate detection. Defaults to false. Changing this forces
    a new resource to be created.

* `support_ordering` - (Optional) Boolean flag which controls whether the Topic
    supports ordering. Defaults to false.

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Topic ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Topic.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Topic.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Topic.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Topic.

## Import

Service Bus Topics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_topic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.servicebus/namespaces/sbns1/topics/sntopic1
```
