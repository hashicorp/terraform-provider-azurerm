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
  name         = "tfex_servicebus_topic"
  namespace_id = azurerm_servicebus_namespace.example.id

  partitioning_enabled = true
}

resource "azurerm_servicebus_subscription" "example" {
  name               = "tfex_servicebus_subscription"
  topic_id           = azurerm_servicebus_topic.example.id
  max_delivery_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Subscription resource. Changing this forces a new resource to be created.

* `topic_id` - (Required) The ID of the ServiceBus Topic to create this Subscription in. Changing this forces a new resource to be created.

* `max_delivery_count` - (Required) The maximum number of deliveries.

* `auto_delete_on_idle` - (Optional) The idle interval after which the topic is automatically deleted as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). The minimum duration is `5` minutes or `PT5M`. Defaults to `P10675199DT2H48M5.4775807S`.

* `default_message_ttl` - (Optional) The Default message timespan to live as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the value used when TimeToLive is not set on a message itself. Defaults to `P10675199DT2H48M5.4775807S`.

* `lock_duration` - (Optional) The lock duration for the subscription as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). The default value is `1` minute or `P0DT0H1M0S` . The maximum value is `5` minutes or `P0DT0H5M0S` . Defaults to `PT1M`.

* `dead_lettering_on_message_expiration` - (Optional) Boolean flag which controls whether the Subscription has dead letter support when a message expires.

* `dead_lettering_on_filter_evaluation_error` - (Optional) Boolean flag which controls whether the Subscription has dead letter support on filter evaluation exceptions. Defaults to `true`.

* `batched_operations_enabled` - (Optional) Boolean flag which controls whether the Subscription supports batched operations.

* `requires_session` - (Optional) Boolean flag which controls whether this Subscription supports the concept of a session. Changing this forces a new resource to be created.

* `forward_to` - (Optional) The name of a Queue or Topic to automatically forward messages to.

* `forward_dead_lettered_messages_to` - (Optional) The name of a Queue or Topic to automatically forward Dead Letter messages to.

* `status` - (Optional) The status of the Subscription. Possible values are `Active`,`ReceiveDisabled`, or `Disabled`. Defaults to `Active`.

* `client_scoped_subscription_enabled` - (Optional) whether the subscription is scoped to a client id. Defaults to `false`.

~> **Note:** Client Scoped Subscription can only be used for JMS subscription (Java Message Service).

* `client_scoped_subscription` - (Optional) A `client_scoped_subscription` block as defined below.

---

A `client_scoped_subscription` block supports the following:

* `client_id` - (Optional) Specifies the Client ID of the application that created the client-scoped subscription. Changing this forces a new resource to be created.

~> **Note:** Client ID can be null or empty, but it must match the client ID set on the JMS client application. From the Azure Service Bus perspective, a null client ID and an empty client id have the same behavior. If the client ID is set to null or empty, it is only accessible to client applications whose client ID is also set to null or empty.

* `is_client_scoped_subscription_shareable` - (Optional) Whether the client scoped subscription is shareable. Defaults to `true` Changing this forces a new resource to be created.

* `is_client_scoped_subscription_durable` - (Optional) Whether the client scoped subscription is durable. This property can only be controlled from the application side.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ServiceBus Subscription ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Subscription.

## Import

Service Bus Subscriptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_subscription.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceBus/namespaces/sbns1/topics/sntopic1/subscriptions/sbsub1
```
