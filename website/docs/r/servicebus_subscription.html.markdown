---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_subscription"
description: |-
  Manages a ServiceBus Subscription.
---

# azurerm_servicebus_subscription

Manages a ServiceBus Subscription.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-servicebus-subscription"
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

resource "azurerm_servicebus_subscription" "example" {
  name                = "tfex_servicebus_subscription"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_servicebus_namespace.example.name
  topic_name          = azurerm_servicebus_topic.example.name
  max_delivery_count  = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Subscription resource. Changing this forces a new resource to be created.

* `namespace_name` - (Required) The name of the ServiceBus Namespace to create this Subscription in. Changing this forces a new resource to be created.

* `topic_name` - (Required) The name of the ServiceBus Topic to create this Subscription in. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the namespace. Changing this forces a new resource to be created.

* `max_delivery_count` - (Required) The maximum number of deliveries.

* `auto_delete_on_idle` - (Optional) The idle interval after which the topic is automatically deleted as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). The minimum duration is `5` minutes or `P5M`.

* `default_message_ttl` - (Optional) The Default message timespan to live as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself.

* `lock_duration` - (Optional) The lock duration for the subscription as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). The default value is `1` minute or `P1M`.

* `dead_lettering_on_message_expiration` - (Optional) Boolean flag which controls whether the Subscription has dead letter support when a message expires. Defaults to `false`.

* `dead_lettering_on_filter_evaluation_error` - (Optional) Boolean flag which controls whether the Subscription has dead letter support on filter evaluation exceptions. Defaults to `true`.

* `enable_batched_operations` - (Optional) Boolean flag which controls whether the Subscription supports batched operations. Defaults to `false`.

* `requires_session` - (Optional) Boolean flag which controls whether this Subscription supports the concept of a session. Defaults to `false`. Changing this forces a new resource to be created.

* `forward_to` - (Optional) The name of a Queue or Topic to automatically forward messages to.

* `forward_dead_lettered_messages_to` - (Optional) The name of a Queue or Topic to automatically forward Dead Letter messages to.

* `status` - (Optional) The status of the Subscription. Possible values are `Active`,`ReceiveDisabled`, or `Disabled`. Defaults to `Active`.

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Subscription ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Subscription.

## Import

Service Bus Subscriptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_subscription.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.servicebus/namespaces/sbns1/topics/sntopic1/subscriptions/sbsub1
```
