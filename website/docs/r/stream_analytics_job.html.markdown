---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_job"
description: |-
  Manages a Stream Analytics Job.
---

# azurerm_stream_analytics_job

Manages a Stream Analytics Job.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_stream_analytics_job" "example" {
  name                                     = "example-job"
  resource_group_name                      = azurerm_resource_group.example.name
  location                                 = azurerm_resource_group.example.location
  compatibility_level                      = "1.1"
  data_locale                              = "en-GB"
  events_late_arrival_max_delay_in_seconds = 60
  events_out_of_order_max_delay_in_seconds = 50
  events_out_of_order_policy               = "Adjust"
  output_error_policy                      = "Drop"
  streaming_units                          = 3

  tags = {
    environment = "Example"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region in which the Resource Group exists. Changing this forces a new resource to be created.

* `compatibility_level` - (Required) Specifies the compatibility level for this job - which controls certain runtime behaviours of the streaming job. Possible values are `1.0` and `1.1`.

-> **NOTE:** Support for Compatibility Level 1.2 is dependent on a new version of the Stream Analytics API, which [being tracked in this issue](https://github.com/Azure/azure-rest-api-specs/issues/5604).

* `data_locale` - (Optional) Specifies the Data Locale of the Job, which [should be a supported .NET Culture](https://msdn.microsoft.com/en-us/library/system.globalization.culturetypes(v=vs.110).aspx).

* `events_late_arrival_max_delay_in_seconds` - (Optional) Specifies the maximum tolerable delay in seconds where events arriving late could be included. Supported range is `-1` (indefinite) to `1814399` (20d 23h 59m 59s).  Default is `0`.

* `events_out_of_order_max_delay_in_seconds` - (Optional) Specifies the maximum tolerable delay in seconds where out-of-order events can be adjusted to be back in order. Supported range is `0` to `599` (9m 59s). Default is `5`.

* `events_out_of_order_policy` - (Optional) Specifies the policy which should be applied to events which arrive out of order in the input event stream. Possible values are `Adjust` and `Drop`.  Default is `Adjust`.

* `output_error_policy` - (Optional) Specifies the policy which should be applied to events which arrive at the output and cannot be written to the external storage due to being malformed (such as missing column values, column values of wrong type or size). Possible values are `Drop` and `Stop`.  Default is `Drop`.

* `streaming_units` - (Required) Specifies the number of streaming units that the streaming job uses. Supported values are `1`, `3`, `6` and multiples of `6` up to `120`.

* `transformation_query` - (Required) Specifies the query that will be run in the streaming job, [written in Stream Analytics Query Language (SAQL)](https://msdn.microsoft.com/library/azure/dn834998).

* `tags` - A mapping of tags assigned to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Stream Analytics Job.

* `job_id` - The Job ID assigned by the Stream Analytics Job.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Job.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Job.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Job.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Job.

## Import

Stream Analytics Job's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_job.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.StreamAnalytics/streamingjobs/job1
```
