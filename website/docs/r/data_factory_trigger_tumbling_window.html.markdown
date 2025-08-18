---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_trigger_tumbling_window"
description: |-
  Manages a Tumbling Window Trigger inside an Azure Data Factory.
---

# azurerm_data_factory_trigger_tumbling_window

Manages a Tumbling Window Trigger inside an Azure Data Factory.

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

resource "azurerm_data_factory_trigger_tumbling_window" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id
  start_time      = "2022-09-21T00:00:00Z"
  end_time        = "2022-09-21T08:00:00Z"
  frequency       = "Minute"
  interval        = 15
  delay           = "16:00:00"

  annotations = ["example1", "example2", "example3"]
  description = "example description"

  retry {
    count    = 1
    interval = 30
  }

  pipeline {
    name = azurerm_data_factory_pipeline.example.name
    parameters = {
      Env = "Prod"
    }
  }

  // Self dependency
  trigger_dependency {
    size   = "24:00:00"
    offset = "-24:00:00"
  }

  additional_properties = {
    foo = "value1"
    bar = "value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Tumbling Window Trigger. Changing this forces a new resource to be created.

* `data_factory_id` - (Required) The ID of Data Factory in which to associate the Trigger with. Changing this forces a new resource.

* `frequency` - (Required) Specifies the frequency of Tumbling Window. Possible values are `Hour`, `Minute` and `Month`. Changing this forces a new resource.

* `interval` - (Required) Specifies the interval of Tumbling Window. Changing this forces a new resource.

* `pipeline` - (Required) A `pipeline` block as defined below.

* `start_time` - (Required) Specifies the start time of Tumbling Window, formatted as an RFC3339 string. Changing this forces a new resource.

* `activated` - (Optional) Specifies if the Data Factory Tumbling Window Trigger is activated. Defaults to `true`.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Tumbling Window Trigger.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Tumbling Window Trigger.

* `delay` - (Optional) Specifies how long the trigger waits before triggering new run. formatted as an `D.HH:MM:SS`.

* `description` - (Optional) The description for the Data Factory Tumbling Window Trigger.

* `end_time` - (Optional) Specifies the end time of Tumbling Window, formatted as an RFC3339 string.

* `max_concurrency` - (Optional) The max number for simultaneous trigger run fired by Tumbling Window. Possible values are between `1` and `50`. Defaults to `50`.

* `retry` - (Optional) A `retry` block as defined below.

* `trigger_dependency` - (Optional) One or more `trigger_dependency` block as defined below.

---

A `pipeline` block supports the following:

* `name` - (Required) The Data Factory Pipeline name that the trigger will act on.

* `parameters` - (Optional) The Data Factory Pipeline parameters that the trigger will act on.

---

A `retry` block supports the following:

* `count` - (Required) The maximum retry attempts if the pipeline run failed.

* `interval` - (Optional) The Interval in seconds between each retry if the pipeline run failed. Defaults to `30`.

---

A `trigger_dependency` block supports the following:

* `offset` - (Optional) The offset of the dependency trigger. Must be in Timespan format (±hh:mm:ss) and must be a negative offset for a self dependency.
  
* `size` - (Optional) The size of the dependency tumbling window. Must be in Timespan format (hh:mm:ss).

* `trigger_name` - (Optional) The dependency trigger name. If not specified, it will use self dependency.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Tumbling Window Trigger.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Tumbling Window Trigger.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Tumbling Window Trigger.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Tumbling Window Trigger.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Tumbling Window Trigger.

## Import

Data Factory Tumbling Window Trigger can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_trigger_tumbling_window.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/triggers/example
```
