---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_trigger_tumbling_window"
description: |-
  Manages a Tumbling Window Trigger.
---

# azurerm_data_factory_trigger_tumbling_window

Manages a Tumbling Window Trigger inside an Azure Data Factory.

## Example Usage

```hcl
resource "azurerm_data_factory_trigger_tumbling_window" "example" {
  name                = "example"
  resource_group_name = "example"
  start_time          = "TODO"
  pipeline_name       = "example"
  data_factory_name   = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Trigger. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the Resource Group where the Tumbling Window Trigger should exist. Changing this forces a new Tumbling Window Trigger to be created.

* `data_factory_name` - (Required) The name of the Azure Data Factory to associate the Trigger with. Changing this forces a new Tumbling Window Trigger to be created.

* `pipeline_name` - (Required) The name of the pipeline to be triggered, there is a one to one relationship between pipelines and Tumbling Window Triggers.

* `start_time` - (Required) The first occurrence, which can be in the past. The first trigger interval is (startTime, startTime + interval). Changing this forces a new Tumbling Window Trigger to be created. Must be in RFC3339 format eg "2020-09-21T00:00:00Z"

* `frequency` - (Required) TODO. Changing this forces a new Tumbling Window Trigger to be created.

* `interval` - (Required) TODO. Changing this forces a new Tumbling Window Trigger to be created.

---

* `annotations` - (Optional) List of tags that can be used for describing the Trigger.

* `delay` - (Optional) The amount of time to delay the start of data processing for the window. The pipeline run is started after the expected execution time plus the amount of delay. The delay defines how long the trigger waits past the due time before triggering a new run. The delay doesn’t alter the window startTime. Must be in Timespan format (hh:mm:ss).

* `end_time` - (Optional) The last occurrence, which can be in the past. Must be in RFC3339 format eg "2020-09-21T00:00:00Z".

* `max_concurrency` - (Optional) TODO.

* `pipeline_parameters` - (Optional) Specifies a list of TODO.

* `retry` - (Optional) A `retry` block as defined below.

* `trigger_dependency` - (Optional) One or more `trigger_dependency` blocks as defined below.

---

A `retry` block supports the following:

* `count` - (Required) TODO.

* `interval` - (Optional) TODO.

---

A `trigger_dependency` block supports the following:

* `offset` - (Optional) TODO.

* `size` - (Optional) TODO.

* `trigger` - (Optional) TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Tumbling Window Trigger.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Tumbling Window Trigger.
* `read` - (Defaults to 5 minutes) Used when retrieving the Tumbling Window Trigger.
* `update` - (Defaults to 30 minutes) Used when updating the Tumbling Window Trigger.
* `delete` - (Defaults to 30 minutes) Used when deleting the Tumbling Window Trigger.

## Import

Tumbling Window Triggers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_trigger_tumbling_window.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ResourceGroupName/providers/Microsoft.DataFactory/factories/DataFactoryName/triggers/TriggerName/
```
