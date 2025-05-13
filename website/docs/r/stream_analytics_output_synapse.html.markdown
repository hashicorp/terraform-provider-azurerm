---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_synapse"
description: |-
  Manages a Stream Analytics Output to an Azure Synapse Analytics Workspace.
---

# azurerm_stream_analytics_output_synapse

Manages a Stream Analytics Output to an Azure Synapse Analytics Workspace.

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
  name                     = "examplestorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "example"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_stream_analytics_output_synapse" "example" {
  name                      = "example-output-synapse"
  stream_analytics_job_name = data.azurerm_stream_analytics_job.example.name
  resource_group_name       = data.azurerm_stream_analytics_job.example.resource_group_name

  server   = azurerm_synapse_workspace.example.connectivity_endpoints["sqlOnDemand"]
  user     = azurerm_synapse_workspace.example.sql_administrator_login
  password = azurerm_synapse_workspace.example.sql_administrator_login_password
  database = "master"
  table    = "ExampleTable"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Output. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `server` - (Required) The name of the SQL server containing the Azure SQL database. Changing this forces a new resource to be created.

* `database` - (Required) The name of the Azure SQL database. Changing this forces a new resource to be created.

* `user` - (Required) The user name that will be used to connect to the Azure SQL database. Changing this forces a new resource to be created.

* `password` - (Required) The password that will be used to connect to the Azure SQL database. 

* `table` - (Required) The name of the table in the Azure SQL database. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Output to an Azure Synapse Analytics Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Output to an Azure Synapse Analytics Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Output to an Azure Synapse Analytics Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Output to an Azure Synapse Analytics Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Output to an Azure Synapse Analytics Workspace.

## Import

A Stream Analytics Output to an Azure Synapse Analytics Workspace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_synapse.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/outputs/output1
```
