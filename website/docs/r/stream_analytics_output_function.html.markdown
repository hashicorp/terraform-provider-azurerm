---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_function"
description: |-
  Manages a Stream Analytics Output Function.
---

# azurerm_stream_analytics_output_function

Manages a Stream Analytics Output Function.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "example" {
  name                = "exampleappserviceplan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "FunctionApp"
  reserved            = true

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "example" {
  name                       = "examplefunctionapp"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  app_service_plan_id        = azurerm_app_service_plan.example.id
  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
  os_type                    = "linux"
  version                    = "~3"
}

resource "azurerm_stream_analytics_job" "example" {
  name                 = "examplestreamanalyticsjob"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  streaming_units      = 3
  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY
}

resource "azurerm_stream_analytics_output_function" "example" {
  name                      = "exampleoutput"
  resource_group_name       = azurerm_stream_analytics_job.example.resource_group_name
  stream_analytics_job_name = azurerm_stream_analytics_job.example.name
  function_app              = azurerm_function_app.example.name
  function_name             = "examplefunctionname"
  api_key                   = "exampleapikey"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Stream Analytics Output. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Output should exist. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `api_key` - (Required) The API key for the Function.

* `function_app` - (Required) The name of the Function App.

* `function_name` - (Required) The name of the function in the Function App.

---

* `batch_max_count` - (Optional) The maximum number of events in each batch that's sent to the function. Defaults to `100`.

* `batch_max_in_bytes` - (Optional) The maximum batch size in bytes that's sent to the function. Defaults to `262144` (256 kB).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Stream Analytics Output Function.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Output Function.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Output Function.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Output Function.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Output Function.

## Import

Stream Analytics Output Functions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_function.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingjobs/job1/outputs/output1
```
