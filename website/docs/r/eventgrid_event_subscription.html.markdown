---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_event_subscription"
description: |-
  Manages an EventGrid Event Subscription

---

# azurerm_eventgrid_event_subscription

Manages an EventGrid Event Subscription

## Example Usage

```hcl
resource "azurerm_resource_group" "default" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "default" {
  name                     = "defaultStorageAccount"
  resource_group_name      = azurerm_resource_group.default.name
  location                 = azurerm_resource_group.default.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "default" {
  name                 = "defaultStorageQueue"
  storage_account_name = azurerm_storage_account.default.name
}

resource "azurerm_eventgrid_event_subscription" "default" {
  name  = "defaultEventSubscription"
  scope = azurerm_resource_group.default.id

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.default.id
    queue_name         = azurerm_storage_queue.default.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventGrid Event Subscription resource. Changing this forces a new resource to be created.

* `scope` - (Required) Specifies the scope at which the EventGrid Event Subscription should be created. Changing this forces a new resource to be created.

* `expiration_time_utc` - (Optional) Specifies the expiration time of the event subscription (Datetime Format `RFC 3339`).

* `event_delivery_schema` - (Optional) Specifies the event delivery schema for the event subscription. Possible values include: `EventGridSchema`, `CloudEventSchemaV1_0`, `CustomInputSchema`. Defaults to `EventGridSchema`. Changing this forces a new resource to be created.

* `azure_function_endpoint` - (Optional) An `azure_function_endpoint` block as defined below.

* `eventhub_endpoint` - (Optional / **Deprecated in favour of `eventhub_endpoint_id`**) A `eventhub_endpoint` block as defined below.

* `eventhub_endpoint_id` - (Optional) Specifies the id where the Event Hub is located.

* `hybrid_connection_endpoint` - (Optional / **Deprecated in favour of `hybrid_connection_endpoint_id`**) A `hybrid_connection_endpoint` block as defined below.

* `hybrid_connection_endpoint_id` - (Optional) Specifies the id where the Hybrid Connection is located.

* `service_bus_queue_endpoint_id` - (Optional) Specifies the id where the Service Bus Queue is located.

* `service_bus_topic_endpoint_id` - (Optional) Specifies the id where the Service Bus Topic is located.

* `storage_queue_endpoint` - (Optional) A `storage_queue_endpoint` block as defined below.

* `webhook_endpoint` - (Optional) A `webhook_endpoint` block as defined below.

~> **NOTE:** One of `eventhub_endpoint`, `eventhub_endpoint_id`, `hybrid_connection_endpoint`, `hybrid_connection_endpoint_id`, `service_bus_queue_endpoint_id`, `service_bus_topic_endpoint_id`, `storage_queue_endpoint` or `webhook_endpoint` must be specified.

* `included_event_types` - (Optional) A list of applicable event types that need to be part of the event subscription.

* `subject_filter` - (Optional) A `subject_filter` block as defined below.

* `advanced_filter` - (Optional) A `advanced_filter` block as defined below.

* `storage_blob_dead_letter_destination` - (Optional) A `storage_blob_dead_letter_destination` block as defined below.

* `retry_policy` - (Optional) A `retry_policy` block as defined below.

* `labels` - (Optional) A list of labels to assign to the event subscription.

---

A `storage_queue_endpoint` supports the following:

* `storage_account_id` - (Required) Specifies the id of the storage account id where the storage queue is located.

* `queue_name` - (Required) Specifies the name of the storage queue where the Event Subscription will receive events.

---

An `azure_function_endpoint` supports the following:

* `function_id` - (Required) Specifies the ID of the Function where the Event Subscription will receive events. This must be the functions ID in format {function_app.id}/functions/{name}.

* `max_events_per_batch` - (Optional) Maximum number of events per batch.

* `preferred_batch_size_in_kilobytes` - (Optional) Preferred batch size in Kilobytes.

---

A `eventhub_endpoint` supports the following:

* `eventhub_id` - (Required) Specifies the id of the eventhub where the Event Subscription will receive events.

---

A `hybrid_connection_endpoint` supports the following:

* `hybrid_connection_id` - (Required) Specifies the id of the hybrid connection where the Event Subscription will receive events.

---

A `webhook_endpoint` supports the following:

* `url` - (Required) Specifies the url of the webhook where the Event Subscription will receive events.

* `base_url` - (Computed) The base url of the webhook where the Event Subscription will receive events.

* `max_events_per_batch` - (Optional) Maximum number of events per batch.

* `preferred_batch_size_in_kilobytes` - (Optional) Preferred batch size in Kilobytes.

* `active_directory_tenant_id` - (Optional) The Azure Active Directory Tenant ID to get the access token that will be included as the bearer token in delivery requests.

* `active_directory_app_id_or_uri` - (Optional) The Azure Active Directory Application ID or URI to get the access token that will be included as the bearer token in delivery requests.

---

A `subject_filter` supports the following:

* `subject_begins_with` - (Optional) A string to filter events for an event subscription based on a resource path prefix.

* `subject_ends_with` - (Optional) A string to filter events for an event subscription based on a resource path suffix.

* `case_sensitive` - (Optional) Specifies if `subject_begins_with` and `subject_ends_with` case sensitive. This value defaults to `false`.

---

A `advanced_filter` supports the following nested blocks:

* `bool_equals` - Compares a value of an event using a single boolean value.
* `number_greater_than` - Compares a value of an event using a single floating point number.
* `number_greater_than_or_equals` - Compares a value of an event using a single floating point number.
* `number_less_than` - Compares a value of an event using a single floating point number.
* `number_less_than_or_equals` - Compares a value of an event using a single floating point number.
* `number_in` - Compares a value of an event using multiple floating point numbers.
* `number_not_in` - Compares a value of an event using multiple floating point numbers.
* `string_begins_with` - Compares a value of an event using multiple string values.
* `string_ends_with` - Compares a value of an event using multiple string values.
* `string_contains` - Compares a value of an event using multiple string values.
* `string_in` - Compares a value of an event using multiple string values.
* `string_not_in` - Compares a value of an event using multiple string values.

Each nested block consists of a key and a value(s) element.

* `key` - (Required) Specifies the field within the event data that you want to use for filtering. Type of the field can be a number, boolean, or string.

* `value` - (Required) Specifies a single value to compare to when using a single value operator. 

**OR** 

* `values` - (Required) Specifies an array of values to compare to when using a multiple values operator.

~> **NOTE:** A maximum of total number of advanced filter values allowed on event subscription is 25.

---

A `storage_blob_dead_letter_destination` supports the following:

* `storage_account_id` - (Required) Specifies the id of the storage account id where the storage blob is located.

* `storage_blob_container_name` - (Required) Specifies the name of the Storage blob container that is the destination of the deadletter events.

---

A `retry_policy` supports the following:

* `max_delivery_attempts` - (Required) Specifies the maximum number of delivery retry attempts for events.

* `event_time_to_live` - (Required) Specifies the time to live (in minutes) for events. Supported range is `1` to `1440`. Defaults to `1440`. See [official documentation](https://docs.microsoft.com/en-us/azure/event-grid/manage-event-delivery#set-retry-policy) for more details.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the EventGrid Event Subscription.

* `topic_name` - (Optional/ **Deprecated) Specifies the name of the topic to associate with the event subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventGrid Event Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the EventGrid Event Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Event Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid Event Subscription.

## Import

EventGrid Event Subscription's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_event_subscription.eventSubscription1
/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/topics/topic1/providers/Microsoft.EventGrid/eventSubscriptions/eventSubscription1
```
