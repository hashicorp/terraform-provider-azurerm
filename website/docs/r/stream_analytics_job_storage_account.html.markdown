---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_job_storage_account"
description: |-
  Manages a Stream Analytics Job Storage Account.
---

# azurerm_stream_analytics_job_storage_account

Manages a Stream Analytics Job Storage Account. Use this resource for managing the Job Storage Account using `Msi` authentication with a `SystemAssigned` identity.

~> **Note:** The Job Storage Account for a Stream Analytics Job can be managed on the `azurerm_stream_analytics_job` resource with the `job_storage_account` block, or with this resource. We do not recommend managing the Job Storage Account through both means as this can lead to conflicts.

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
  compatibility_level                      = "1.2"
  data_locale                              = "en-GB"
  events_late_arrival_max_delay_in_seconds = 60
  events_out_of_order_max_delay_in_seconds = 50
  events_out_of_order_policy               = "Adjust"
  output_error_policy                      = "Drop"
  streaming_units                          = 3
  sku_name                                 = "StandardV2"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "Example"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

  lifecycle {
    ignore_changes = [job_storage_account]
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "exampleaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_stream_analytics_job_storage_account" "example" {
  stream_analytics_job_id = azurerm_stream_analytics_job.example.id
  storage_account_name    = azurerm_storage_account.example.name
  authentication_mode     = "Msi"
}
```

## Argument Reference

The following arguments are supported:

* `stream_analytics_job_id` - (Required) The ID of the Stream Analytics Job. Changing this forces a new resource to be created.

* `authentication_mode` - (Required) The authentication mode for the Stream Analytics Job's Storage Account. Possible values are `ConnectionString`, and `Msi`.

-> **Note:** The parent Stream Analytics Job must have the `identity` block set when using `Msi` as the `authentication_mode`.

* `storage_account_name` - (Required) The Storage Account name for the Stream Analytics Job.

* `storage_account_key` - (Optional) The Storage Account Key for accessing the Storage Account of the Stream Analytics Job.

-> **Note:** `storage_account_key` must be specified when `authentication_mode` is `ConnectionString` and must be absent if `authentication_mode` is `Msi`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Job.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Job Storage Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Job Storage Account.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Job Storage Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Job Storage Account.

## Import

Stream Analytics Job Storage Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_job_storage_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1
```
