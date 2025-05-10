---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_blob"
description: |-
  Manages a Stream Analytics Output to Blob Storage.
---

# azurerm_stream_analytics_output_blob

Manages a Stream Analytics Output to Blob Storage.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

data "azurerm_stream_analytics_job" "example" {
  name                = "example-job"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
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

resource "azurerm_stream_analytics_output_blob" "example" {
  name                      = "output-to-blob-storage"
  stream_analytics_job_name = data.azurerm_stream_analytics_job.example.name
  resource_group_name       = data.azurerm_stream_analytics_job.example.resource_group_name
  storage_account_name      = azurerm_storage_account.example.name
  storage_account_key       = azurerm_storage_account.example.primary_access_key
  storage_container_name    = azurerm_storage_container.example.name
  path_pattern              = "some-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Output. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `date_format` - (Required) The date format. Wherever `{date}` appears in `path_pattern`, the value of this property is used as the date format instead.

* `path_pattern` - (Required) The blob path pattern. Not a regular expression. It represents a pattern against which blob names will be matched to determine whether or not they should be included as input or output to the job.

* `storage_account_name` - (Required) The name of the Storage Account.

* `storage_container_name` - (Required) The name of the Container within the Storage Account.

* `time_format` - (Required) The time format. Wherever `{time}` appears in `path_pattern`, the value of this property is used as the time format instead.

* `serialization` - (Required) A `serialization` block as defined below.

* `authentication_mode` - (Optional) The authentication mode for the Stream Output. Possible values are `Msi` and `ConnectionString`. Defaults to `ConnectionString`.

* `batch_max_wait_time` - (Optional) The maximum wait time per batch in `hh:mm:ss` e.g. `00:02:00` for two minutes.

* `batch_min_rows` - (Optional) The minimum number of rows per batch (must be between `0` and `1000000`).

* `storage_account_key` - (Optional) The Access Key which should be used to connect to this Storage Account.

* `blob_write_mode` - (Optional) Determines whether blob blocks are either committed automatically or appended. Possible values are `Append` and `Once`. Defaults to `Append`.

---

A `serialization` block supports the following:

* `type` - (Required) The serialization format used for outgoing data streams. Possible values are `Avro`, `Csv`, `Json` and `Parquet`.

-> **Note:** `batch_max_wait_time` and `batch_min_rows` are required when `type` is set to `Parquet`

* `encoding` - (Optional) The encoding of the incoming data in the case of input and the encoding of outgoing data in the case of output. It currently can only be set to `UTF8`.

-> **Note:** This is required when `type` is set to `Csv` or `Json`.

* `field_delimiter` - (Optional) The delimiter that will be used to separate comma-separated value (CSV) records. Possible values are ` ` (space), `,` (comma), `	` (tab), `|` (pipe) and `;`.

-> **Note:** This is required when `type` is set to `Csv`.

* `format` - (Optional) Specifies the format of the JSON the output will be written in. Possible values are `Array` and `LineSeparated`.

-> **Note:** This is Required and can only be specified when `type` is set to `Json`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Output Blob Storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Output Blob Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Output Blob Storage.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Output Blob Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Output Blob Storage.

## Import

Stream Analytics Outputs to Blob Storage can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_blob.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/outputs/output1
```
