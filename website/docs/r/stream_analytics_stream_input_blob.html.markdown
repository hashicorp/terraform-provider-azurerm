---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_stream_input_blob"
description: |-
  Manages a Stream Analytics Stream Input Blob.
---

# azurerm_stream_analytics_stream_input_blob

Manages a Stream Analytics Stream Input Blob.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_stream_analytics_job" "example" {
  name                = "example-job"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
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

resource "azurerm_stream_analytics_stream_input_blob" "example" {
  name                      = "blob-stream-input"
  stream_analytics_job_name = data.azurerm_stream_analytics_job.example.name
  resource_group_name       = data.azurerm_stream_analytics_job.example.resource_group_name
  storage_account_name      = azurerm_storage_account.example.name
  storage_account_key       = azurerm_storage_account.example.primary_access_key
  storage_container_name    = azurerm_storage_container.example.name
  path_pattern              = "some-random-pattern"
  date_format               = "yyyy/MM/dd"
  time_format               = "HH"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Input Blob. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `date_format` - (Required) The date format. Wherever `{date}` appears in `path_pattern`, the value of this property is used as the date format instead.

* `path_pattern` - (Required) The blob path pattern. Not a regular expression. It represents a pattern against which blob names will be matched to determine whether or not they should be included as input or output to the job.

* `storage_account_name` - (Required) The name of the Storage Account.

* `storage_account_key` - (Required) The Access Key which should be used to connect to this Storage Account.

* `storage_container_name` - (Required) The name of the Container within the Storage Account.

* `time_format` - (Required) The time format. Wherever `{time}` appears in `path_pattern`, the value of this property is used as the time format instead.

* `serialization` - (Required) A `serialization` block as defined below.

* `authentication_mode` - (Optional) The authentication mode for the Stream Analytics Input. Possible values are `Msi` and `ConnectionString`. Defaults to `ConnectionString`.

---

A `serialization` block supports the following:

* `type` - (Required) The serialization format used for incoming data streams. Possible values are `Avro`, `Csv` and `Json`.

* `encoding` - (Optional) The encoding of the incoming data in the case of input and the encoding of outgoing data in the case of output. It currently can only be set to `UTF8`.

-> **Note:** This is required when `type` is set to `Csv` or `Json`.

* `field_delimiter` - (Optional) The delimiter that will be used to separate comma-separated value (CSV) records. Possible values are ` ` (space), `,` (comma), `	` (tab), `|` (pipe) and `;`.

-> **Note:** This is required when `type` is set to `Csv`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Stream Input Blob.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Stream Input Blob.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Stream Input Blob.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Stream Input Blob.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Stream Input Blob.

## Import

Stream Analytics Stream Input Blob's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_stream_input_blob.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/inputs/input1
```
