---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_subscription"
description: |-
  Gets information about an existing ServiceBus Subscription.
---

# Data Source: azurerm_servicebus_subscription

Use this data source to access information about an existing ServiceBus Subscription.

## Example Usage

```hcl
data "azurerm_servicebus_subscription" "example" {
  name     = "examplesubscription"
  topic_id = "exampletopic"
}

output "servicebus_subscription" {
  value = data.azurerm_servicebus_namespace.example
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the ServiceBus Subscription.

* `topic_id` - (Required) The ID of the ServiceBus Topic where the Service Bus Subscription exists.

## Attributes Reference

* `max_delivery_count` - The maximum number of deliveries.

* `auto_delete_on_idle` - The idle interval after which the topic is automatically deleted.

* `default_message_ttl` - The Default message timespan to live. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself.

* `lock_duration` - The lock duration for the subscription.

* `dead_lettering_on_message_expiration` - Does the Service Bus Subscription have dead letter support when a message expires?

* `dead_lettering_on_filter_evaluation_error` - Does the ServiceBus Subscription have dead letter support on filter evaluation exceptions?

* `enable_batched_operations` - Are batched operations enabled on this ServiceBus Subscription?

* `requires_session` - Whether or not this ServiceBus Subscription supports session.

* `forward_to` - The name of a ServiceBus Queue or ServiceBus Topic where messages are automatically forwarded.

* `forward_dead_lettered_messages_to` - The name of a Queue or Topic to automatically forward Dead Letter messages to.

* `client_scoped_subscription_enabled` - Does the subscription scoped to a client id or not.

* `client_scoped_subscription` - (Optional)  A `client_scoped_subscription` block as defined below.

---

* `client_id` - The Client ID of the application that created the client-scoped subscription.

* `is_client_scoped_subscription_shareable` - The client scoped subscription is shareable or not.

* `is_client_scoped_subscription_durable` - The client scoped subscription is durable or not.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Subscription.
