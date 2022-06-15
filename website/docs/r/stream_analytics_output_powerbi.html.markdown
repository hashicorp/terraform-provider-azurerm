---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_powerbi"
description: |-
  Manages a Stream Analytics Output powerBI.
---

# azurerm_stream_analytics_output_powerbi

Manages a Stream Analytics Output powerBI.

## Example Usage

```hcl
data "azurerm_resource_group" "example" {
  name = "example-resources"
}

data "azurerm_stream_analytics_job" "example" {
  name                = "example-job"
  resource_group_name = data.azurerm_resource_group.example.name
}

resource "azurerm_stream_analytics_output_powerbi" "example" {
  name                    = "output-to-powerbi"
  stream_analytics_job_id = data.azurerm_stream_analytics_job.example.id
  dataset                 = "example-dataset"
  table                   = "example-table"
  group_id                = "00000000-0000-0000-0000-000000000000"
  group_name              = "some-group-name"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Output. Changing this forces a new resource to be created.

* `stream_analytics_job_id` - (Required) The ID of the Stream Analytics Job. Changing this forces a new resource to be created.

* `dataset` - (Required) The name of the Power BI dataset.

* `table` - (Required) The name of the Power BI table under the specified dataset.

* `group_id` - (Required) The ID of the Power BI group, this must be a valid UUID.

* `group_name` - (Required) The name of the Power BI group. Use this property to help remember which specific Power BI group id was used.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Output for PowerBI.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Output for PowerBI.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Output for PowerBI.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Output for PowerBI.

## Import

Stream Analytics Output to Power BI can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_powerbi.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingjobs/job1/outputs/output1
```
