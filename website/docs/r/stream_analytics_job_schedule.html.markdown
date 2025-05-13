---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_job_schedule"
description: |-
  Manages a Stream Analytics Job Schedule.
---

# azurerm_stream_analytics_job_schedule

Manages a Stream Analytics Job Schedule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "example"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "example" {
  name                   = "example"
  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.example.name
  type                   = "Block"
  source                 = "example.csv"
}

resource "azurerm_stream_analytics_job" "example" {
  name                                     = "example-job"
  resource_group_name                      = azurerm_resource_group.example.name
  location                                 = azurerm_resource_group.example.location
  compatibility_level                      = "1.2"
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
    INTO [exampleoutput]
    FROM [exampleinput]
QUERY
}

resource "azurerm_stream_analytics_stream_input_blob" "example" {
  name                      = "exampleinput"
  stream_analytics_job_name = azurerm_stream_analytics_job.example.name
  resource_group_name       = azurerm_stream_analytics_job.example.resource_group_name
  storage_account_name      = azurerm_storage_account.example.name
  storage_account_key       = azurerm_storage_account.example.primary_access_key
  storage_container_name    = azurerm_storage_container.example.name
  path_pattern              = "" // this checks for blobs in the root of the container
  date_format               = "yyyy/MM/dd"
  time_format               = "HH"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}

resource "azurerm_stream_analytics_output_blob" "example" {
  name                      = "exampleoutput"
  stream_analytics_job_name = azurerm_stream_analytics_job.example.name
  resource_group_name       = azurerm_stream_analytics_job.example.resource_group_name
  storage_account_name      = azurerm_storage_account.example.name
  storage_account_key       = azurerm_storage_account.example.primary_access_key
  storage_container_name    = azurerm_storage_container.example.name
  path_pattern              = "example-{date}-{time}"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type = "Avro"
  }
}

resource "azurerm_stream_analytics_job_schedule" "example" {
  stream_analytics_job_id = azurerm_stream_analytics_job.example.id
  start_mode              = "CustomTime"
  start_time              = "2022-09-21T00:00:00Z"

  depends_on = [
    azurerm_stream_analytics_job.example,
    azurerm_stream_analytics_stream_input_blob.example,
    azurerm_stream_analytics_output_blob.example,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `stream_analytics_job_id` - (Required) The ID of the Stream Analytics Job that should be scheduled or started. Changing this forces a new resource to be created.

* `start_mode` - (Required) The starting mode of the Stream Analytics Job. Possible values are `JobStartTime`, `CustomTime` and `LastOutputEventTime`.

-> **Note:** Setting `start_mode` to `LastOutputEventTime` is only possible if the job had been previously started and produced output.

* `start_time` - (Optional) The time in ISO8601 format at which the Stream Analytics Job should be started e.g. `2022-04-01T00:00:00Z`. This property can only be specified if `start_mode` is set to `CustomTime`

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Job.

* `last_output_time` - The time at which the Stream Analytics job last produced an output.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Job.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Job.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Job.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Job.

## Import

Stream Analytics Job's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_job_schedule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/schedule/default
```
