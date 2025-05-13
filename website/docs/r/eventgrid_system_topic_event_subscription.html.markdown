---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_system_topic_event_subscription"
description: |-
  Manages an EventGrid System Topic Event Subscription.
---

# azurerm_eventgrid_system_topic_event_subscription

Manages an EventGrid System Topic Event Subscription.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_queue" "example" {
  name                 = "examplestoragequeue"
  storage_account_name = azurerm_storage_account.example.name
}

resource "azurerm_eventgrid_system_topic" "example" {
  name                   = "example-system-topic"
  location               = "Global"
  resource_group_name    = azurerm_resource_group.example.name
  source_arm_resource_id = azurerm_resource_group.example.id
  topic_type             = "Microsoft.Resources.ResourceGroups"
}

resource "azurerm_eventgrid_system_topic_event_subscription" "example" {
  name                = "example-event-subscription"
  system_topic        = azurerm_eventgrid_system_topic.example.name
  resource_group_name = azurerm_resource_group.example.name

  storage_queue_endpoint {
    storage_account_id = azurerm_storage_account.example.id
    queue_name         = azurerm_storage_queue.example.name
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Event Subscription. Changing this forces a new Event Subscription to be created.

* `system_topic` - (Required) The System Topic where the Event Subscription should be created in. Changing this forces a new Event Subscription to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the System Topic exists. Changing this forces a new Event Subscription to be created.

* `expiration_time_utc` - (Optional) Specifies the expiration time of the event subscription (Datetime Format `RFC 3339`).

* `event_delivery_schema` - (Optional) Specifies the event delivery schema for the event subscription. Possible values include: `EventGridSchema`, `CloudEventSchemaV1_0`, `CustomInputSchema`. Defaults to `EventGridSchema`. Changing this forces a new resource to be created.

* `azure_function_endpoint` - (Optional) An `azure_function_endpoint` block as defined below.

* `eventhub_endpoint_id` - (Optional) Specifies the id where the Event Hub is located.

* `hybrid_connection_endpoint_id` - (Optional) Specifies the id where the Hybrid Connection is located.

* `service_bus_queue_endpoint_id` - (Optional) Specifies the id where the Service Bus Queue is located.

* `service_bus_topic_endpoint_id` - (Optional) Specifies the id where the Service Bus Topic is located.

* `storage_queue_endpoint` - (Optional) A `storage_queue_endpoint` block as defined below.

* `webhook_endpoint` - (Optional) A `webhook_endpoint` block as defined below.

~> **Note:** One of `azure_function_endpoint`, `eventhub_endpoint_id`, `hybrid_connection_endpoint`, `hybrid_connection_endpoint_id`, `service_bus_queue_endpoint_id`, `service_bus_topic_endpoint_id`, `storage_queue_endpoint` or `webhook_endpoint` must be specified.

* `included_event_types` - (Optional) A list of applicable event types that need to be part of the event subscription.

* `subject_filter` - (Optional) A `subject_filter` block as defined below.

* `advanced_filter` - (Optional) A `advanced_filter` block as defined below.

* `delivery_identity` - (Optional) A `delivery_identity` block as defined below.

* `delivery_property` - (Optional) One or more `delivery_property` blocks as defined below.

* `dead_letter_identity` - (Optional) A `dead_letter_identity` block as defined below.

-> **Note:** `storage_blob_dead_letter_destination` must be specified when a `dead_letter_identity` is specified

* `storage_blob_dead_letter_destination` - (Optional) A `storage_blob_dead_letter_destination` block as defined below.

* `retry_policy` - (Optional) A `retry_policy` block as defined below.

* `labels` - (Optional) A list of labels to assign to the event subscription.

* `advanced_filtering_on_arrays_enabled` - (Optional) Specifies whether advanced filters should be evaluated against an array of values instead of expecting a singular value. Defaults to `false`.

---

A `storage_queue_endpoint` block supports the following:

* `storage_account_id` - (Required) Specifies the id of the storage account id where the storage queue is located.

* `queue_name` - (Required) Specifies the name of the storage queue where the Event Subscription will receive events.

* `queue_message_time_to_live_in_seconds` - (Optional) Storage queue message time to live in seconds.

---

An `azure_function_endpoint` block supports the following:

* `function_id` - (Required) Specifies the ID of the Function where the Event Subscription will receive events. This must be the functions ID in format {function_app.id}/functions/{name}.

* `max_events_per_batch` - (Optional) Maximum number of events per batch.

* `preferred_batch_size_in_kilobytes` - (Optional) Preferred batch size in Kilobytes.

---

A `webhook_endpoint` block supports the following:

* `url` - (Required) Specifies the url of the webhook where the Event Subscription will receive events.

* `base_url` - (Computed) The base url of the webhook where the Event Subscription will receive events.

* `max_events_per_batch` - (Optional) Maximum number of events per batch.

* `preferred_batch_size_in_kilobytes` - (Optional) Preferred batch size in Kilobytes.

* `active_directory_tenant_id` - (Optional) The Azure Active Directory Tenant ID to get the access token that will be included as the bearer token in delivery requests.

* `active_directory_app_id_or_uri` - (Optional) The Azure Active Directory Application ID or URI to get the access token that will be included as the bearer token in delivery requests.

---

A `subject_filter` block supports the following:

* `subject_begins_with` - (Optional) A string to filter events for an event subscription based on a resource path prefix.

* `subject_ends_with` - (Optional) A string to filter events for an event subscription based on a resource path suffix.

* `case_sensitive` - (Optional) Specifies if `subject_begins_with` and `subject_ends_with` case sensitive. This value 

---

A `advanced_filter` supports the following nested blocks:

* `bool_equals` - (Optional) Compares a value of an event using a single boolean value.
* `number_greater_than` - (Optional) Compares a value of an event using a single floating point number.
* `number_greater_than_or_equals` - (Optional) Compares a value of an event using a single floating point number.
* `number_less_than` - (Optional) Compares a value of an event using a single floating point number.
* `number_less_than_or_equals` - (Optional) Compares a value of an event using a single floating point number.
* `number_in` - (Optional) Compares a value of an event using multiple floating point numbers.
* `number_not_in` - (Optional) Compares a value of an event using multiple floating point numbers.
* `number_in_range` - (Optional) Compares a value of an event using multiple floating point number ranges.
* `number_not_in_range` - (Optional) Compares a value of an event using multiple floating point number ranges.
* `string_begins_with` - (Optional) Compares a value of an event using multiple string values.
* `string_not_begins_with` - (Optional) Compares a value of an event using multiple string values.
* `string_ends_with` - (Optional) Compares a value of an event using multiple string values.
* `string_not_ends_with` - (Optional) Compares a value of an event using multiple string values.
* `string_contains` - (Optional) Compares a value of an event using multiple string values.
* `string_not_contains` - (Optional) Compares a value of an event using multiple string values.
* `string_in` - (Optional) Compares a value of an event using multiple string values.
* `string_not_in` - (Optional) Compares a value of an event using multiple string values.
* `is_not_null` - (Optional) Evaluates if a value of an event isn't NULL or undefined.
* `is_null_or_undefined` - (Optional) Evaluates if a value of an event is NULL or undefined.

Each nested block consists of a key and a value(s) element.

* `key` - (Required) Specifies the field within the event data that you want to use for filtering. Type of the field can be a number, boolean, or string.

* `value` - (Required) Specifies a single value to compare to when using a single value operator.

OR

* `values` - (Required) Specifies an array of values to compare to when using a multiple values operator.

~> **Note:** A maximum of total number of advanced filter values allowed on event subscription is 25.

---

A `delivery_identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that is used for event delivery. Allowed value is `SystemAssigned`, `UserAssigned`.

* `user_assigned_identity` - (Optional) The user identity associated with the resource.

---

A `delivery_property` block supports the following:

~> **Note:** `delivery_property` blocks are only effective when using an `azure_function_endpoint`, `eventhub_endpoint_id`, `hybrid_connection_endpoint_id`, `service_bus_topic_endpoint_id`, or `webhook_endpoint` endpoint specification.

* `header_name` - (Required) The name of the header to send on to the destination.

* `type` - (Required) Either `Static` or `Dynamic`.

* `value` - (Optional) If the `type` is `Static`, then provide the value to use.

* `source_field` - (Optional) If the `type` is `Dynamic`, then provide the payload field to be used as the value. Valid source fields differ by subscription type.

* `secret` - (Optional) Set to `true` if the `value` is a secret and should be protected, otherwise `false`. If `true` then this value won't be returned from Azure API calls.

---

A `dead_letter_identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that is used for dead lettering. Allowed value is `SystemAssigned`, `UserAssigned`.

* `user_assigned_identity` - (Optional) The user identity associated with the resource.

---

A `storage_blob_dead_letter_destination` block supports the following:

* `storage_account_id` - (Required) Specifies the id of the storage account id where the storage blob is located.

* `storage_blob_container_name` - (Required) Specifies the name of the Storage blob container that is the destination of the deadletter events.

---

A `retry_policy` block supports the following:

* `max_delivery_attempts` - (Required) Specifies the maximum number of delivery retry attempts for events.

* `event_time_to_live` - (Required) Specifies the time to live (in minutes) for events. Supported range is `1` to `1440`. See [official documentation](https://docs.microsoft.com/azure/event-grid/manage-event-delivery#set-retry-policy) for more details.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the EventGrid System Topic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Messaging.
* `read` - (Defaults to 5 minutes) Used when retrieving the Messaging.
* `update` - (Defaults to 30 minutes) Used when updating the Messaging.
* `delete` - (Defaults to 30 minutes) Used when deleting the Messaging.

## Import

EventGrid System Topic Event Subscriptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_system_topic_event_subscription.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/systemTopics/topic1/eventSubscriptions/subscription1
```
