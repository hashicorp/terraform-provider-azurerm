---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_table"
description: |-
  Manages a Stream Analytics Output Table.
---

# azurerm_stream_analytics_output_table

Manages a Stream Analytics Output Table.

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

resource "azurerm_storage_table" "example" {
  name                 = "exampletable"
  storage_account_name = azurerm_storage_account.example.name
}

resource "azurerm_stream_analytics_output_table" "example" {
  name                      = "output-to-storage-table"
  stream_analytics_job_name = data.azurerm_stream_analytics_job.example.name
  resource_group_name       = data.azurerm_stream_analytics_job.example.resource_group_name
  storage_account_name      = azurerm_storage_account.example.name
  storage_account_key       = azurerm_storage_account.example.primary_access_key
  table                     = azurerm_storage_table.example.name
  partition_key             = "foo"
  row_key                   = "bar"
  batch_size                = 100
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Output. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) The name of the Storage Account.

* `storage_account_key` - (Required) The Access Key which should be used to connect to this Storage Account.

* `table` - (Required) The name of the table where the stream should be output to.

* `partition_key` - (Required) The name of the output column that contains the partition key.

* `row_key` - (Required) The name of the output column that contains the row key.

* `batch_size` - (Required) The number of records for a batch operation. Must be between `1` and `100`.

* `columns_to_remove` - (Optional) A list of the column names to be removed from output event entities.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Stream Analytics Output Table.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics.

## Import

Stream Analytics Output to Table can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_table.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingjobs/job1/outputs/output1
```
