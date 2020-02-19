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
  name     = "defaultResourceGroup"
  location = "West US 2"
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
  resource_group_name  = azurerm_resource_group.default.name
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

* `event_delivery_schema` - (Optional) Specifies the event delivery schema for the event subscription. Possible values include: `EventGridSchema`, `CloudEventV01Schema`, `CustomInputSchema`.

* `topic_name` - (Optional) Specifies the name of the topic to associate with the event subscription.

* `storage_queue_endpoint` - (Optional) A `storage_queue_endpoint` block as defined below.

* `eventhub_endpoint` - (Optional) A `eventhub_endpoint` block as defined below.

* `hybrid_connection_endpoint` - (Optional) A `hybrid_connection_endpoint` block as defined below.

* `webhook_endpoint` - (Optional) A `webhook_endpoint` block as defined below.

~> **NOTE:** One of `storage_queue_endpoint`, `eventhub_endpoint`, `hybrid_connection_endpoint` or `webhook_endpoint` must be specified.

* `included_event_types` - (Optional) A list of applicable event types that need to be part of the event subscription.

* `subject_filter` - (Optional) A `subject_filter` block as defined below.

* `storage_blob_dead_letter_destination` - (Optional) A `storage_blob_dead_letter_destination` block as defined below.

* `retry_policy` - (Optional) A `retry_policy` block as defined below.

* `labels` - (Optional) A list of labels to assign to the event subscription.

---

A `storage_queue_endpoint` supports the following:

* `storage_account_id` - (Required) Specifies the id of the storage account id where the storage queue is located.

* `queue_name` - (Required) Specifies the name of the storage queue where the Event Subscriptio will receive events.

---

A `eventhub_endpoint` supports the following:

* `eventhub_id` - (Required) Specifies the id of the eventhub where the Event Subscription will receive events.

---

A `hybrid_connection_endpoint` supports the following:

* `hybrid_connection_id` - (Required) Specifies the id of the hybrid connection where the Event Subscription will receive events.

A `webhook_endpoint` supports the following:

* `url` - (Required) Specifies the url of the webhook where the Event Subscription will receive events.

---

A `subject_filter` supports the following:

* `subject_begins_with` - (Optional) A string to filter events for an event subscription based on a resource path prefix.

* `subject_ends_with` - (Optional) A string to filter events for an event subscription based on a resource path suffix.

* `case_sensitive` - (Optional) Specifies if `subject_begins_with` and `subject_ends_with` case sensitive. This value defaults to `false`.

---

A `storage_blob_dead_letter_destination` supports the following:

* `storage_account_id` - (Required) Specifies the id of the storage account id where the storage blob is located.

* `storage_blob_container_name` - (Required) Specifies the name of the Storage blob container that is the destination of the deadletter events

---

A `retry_policy` supports the following:

* `max_delivery_attempts` - (Required) Specifies the maximum number of delivery retry attempts for events.

* `event_time_to_live` - (Required) Specifies the time to live (in minutes) for events.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the EventGrid Event Subscription.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventGrid Event Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the EventGrid Event Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Event Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid Event Subscription.

## Import

EventGrid Domain's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_event_subscription.eventSubscription1
/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/eventSubscriptions/eventSubscription1
```
