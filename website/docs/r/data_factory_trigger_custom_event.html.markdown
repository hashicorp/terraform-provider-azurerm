---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_trigger_custom_event"
description: |-
  Manages a Custom Event Trigger inside an Azure Data Factory.
---

# azurerm_data_factory_trigger_custom_event

Manages a Custom Event Trigger inside an Azure Data Factory.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_pipeline" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id
}

resource "azurerm_eventgrid_topic" "example" {
  name                = "example-topic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_trigger_custom_event" "example" {
  name                = "example"
  data_factory_id     = azurerm_data_factory.example.id
  eventgrid_topic_id  = azurerm_eventgrid_topic.example.id
  events              = ["event1", "event2"]
  subject_begins_with = "abc"
  subject_ends_with   = "xyz"

  annotations = ["example1", "example2", "example3"]
  description = "example description"

  pipeline {
    name = azurerm_data_factory_pipeline.example.name
    parameters = {
      Env = "Prod"
    }
  }

  additional_properties = {
    foo = "foo1"
    bar = "bar2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Custom Event Trigger. Changing this forces a new resource to be created.

* `data_factory_id` - (Required) The ID of Data Factory in which to associate the Trigger with. Changing this forces a new resource.

* `eventgrid_topic_id` - (Required) The ID of Event Grid Topic in which event will be listened. Changing this forces a new resource.

* `events` - (Required) List of events that will fire this trigger. At least one event must be specified.

* `pipeline` - (Required) One or more `pipeline` blocks as defined below.

* `activated` - (Optional) Specifies if the Data Factory Custom Event Trigger is activated. Defaults to `true`.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Custom Event Trigger.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Custom Event Trigger.

* `description` - (Optional) The description for the Data Factory Custom Event Trigger.

* `subject_begins_with` - (Optional) The pattern that event subject starts with for trigger to fire.

* `subject_ends_with` - (Optional) The pattern that event subject ends with for trigger to fire.

---

A `pipeline` block supports the following:

* `name` - (Required) The Data Factory Pipeline name that the trigger will act on.

* `parameters` - (Optional) The Data Factory Pipeline parameters that the trigger will act on.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Custom Event Trigger.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Custom Event Trigger.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Custom Event Trigger.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Custom Event Trigger.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Custom Event Trigger.

## Import

Data Factory Custom Event Trigger can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_trigger_custom_event.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/triggers/example
```
