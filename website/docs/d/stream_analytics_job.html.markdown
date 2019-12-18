---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_job"
sidebar_current: "docs-azurerm-datasource-stream-analytics-job"
description: |-
  Gets information about an existing Stream Analytics Job.

---

# Data Source: azurerm_stream_analytics_job

Use this data source to access information about an existing Stream Analytics Job.

## Example Usage

```hcl
data "azurerm_stream_analytics_job" "example" {
  name                = "example-job"
  resource_group_name = "example-resources"
}

output "job_id" {
  value = "${data.azurerm_stream_analytics_job.example.job_id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Stream Analytics Job.

* `resource_group_name` - (Required) Specifies the name of the resource group the Stream Analytics Job is located in.

## Attributes Reference

* `id` - The ID of the Stream Analytics Job.

* `compatibility_level` - The compatibility level for this job.

* `data_locale` - The Data Locale of the Job.

* `events_late_arrival_max_delay_in_seconds` - The maximum tolerable delay in seconds where events arriving late could be included.

* `events_out_of_order_max_delay_in_seconds` - The maximum tolerable delay in seconds where out-of-order events can be adjusted to be back in order.

* `events_out_of_order_policy` - The policy which should be applied to events which arrive out of order in the input event stream.

* `job_id` - The Job ID assigned by the Stream Analytics Job.

* `location` - The Azure location where the Stream Analytics Job exists.

* `output_error_policy` - The policy which should be applied to events which arrive at the output and cannot be written to the external storage due to being malformed (such as missing column values, column values of wrong type or size). 

* `streaming_units` - The number of streaming units that the streaming job uses.

* `transformation_query` - The query that will be run in the streaming job, [written in Stream Analytics Query Language (SAQL)](https://msdn.microsoft.com/library/azure/dn834998).

