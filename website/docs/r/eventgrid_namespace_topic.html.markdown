---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_namespace_topic"
description: |-
  Manages a Event Grid Namespace Topic.
---

# azurerm_eventgrid_namespace_topic

Manages a Event Grid Namespace Topic.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventgrid_namespace" "example" {
  name                = "my-eventgrid-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_eventgrid_namespace_topic" "example" {
  name                    = "topic-namespace-example"
  namespace_id            = azurerm_eventgrid_namespace.test.id
  event_retention_in_days = 1
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Event Grid Namespace Topic. Changing this forces a new Event Grid Namespace Topic to be created.

* `namespace_id` - (Required) The ID of the Event Grid Namespace Topic. Changing this forces a new Event Grid Namespace Topic to be created.

---

* `event_retention_in_days` - (Optional) Event retention for the namespace topic expressed in days. The property default value is 1 day. Min event retention duration value is 1 day and max event retention duration value is 1 day. Defaults to `7`.

* `input_schema` - (Optional) This determines the format that is expected for incoming events published to the topic. Defaults to `CloudEventSchemaV1_0`.

* `publisher_type` - (Optional) Publisher type of the namespace topic. Defaults to `Custom`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Event Grid Namespace Topic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Event Grid Namespace Topic.
* `read` - (Defaults to 5 minutes) Used when retrieving the Event Grid Namespace Topic.
* `update` - (Defaults to 30 minutes) Used when updating the Event Grid Namespace Topic.
* `delete` - (Defaults to 30 minutes) Used when deleting the Event Grid Namespace Topic.

## Import

Event Grid Namespace Topics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_namespace_topic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/namespaces/eventgrid1/topics/topic1
```