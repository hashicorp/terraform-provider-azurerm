---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_cosmosdb"
description: |-
  Manages a Stream Analytics Output to CosmosDB.
---

# azurerm_stream_analytics_output_cosmosdb

Manages a Stream Analytics Output to CosmosDB.

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

resource "azurerm_cosmosdb_account" "example" {
  name                = "exampledb"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "cosmos-sql-db"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "example" {
  name                = "examplecontainer"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_sql_database.example.name
  partition_key_path  = "foo"
}

resource "azurerm_stream_analytics_output_cosmosdb" "example" {
  name                     = "output-to-cosmosdb"
  stream_analytics_job_id  = data.azurerm_stream_analytics_job.example.id
  cosmosdb_account_key     = azurerm_cosmosdb_account.example.primary_key
  cosmosdb_sql_database_id = azurerm_cosmosdb_sql_database.example.id
  container_name           = azurerm_cosmosdb_sql_container.example.name
  document_id              = "exampledocumentid"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Analytics Output. Changing this forces a new resource to be created.

* `stream_analytics_job_id` - (Required) The ID of the Stream Analytics Job. Changing this forces a new resource to be created.

* `cosmosdb_account_key` - (Required) The account key for the CosmosDB database.

* `cosmosdb_sql_database_id` - (Required) The ID of the CosmosDB database.

* `container_name` - (Required) The name of the CosmosDB container.

* `document_id` - (Optional) The name of the field in output events used to specify the primary key which insert or update operations are based on.

* `partition_key` - (Optional) The name of the field in output events used to specify the key for partitioning output across collections. If `container_name` contains `{partition}` token, this property is required to be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Output for CosmosDB.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Output for CosmosDB.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Output for CosmosDB.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Output for CosmosDB.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Output for CosmosDB.

## Import

Stream Analytics Outputs for CosmosDB can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_cosmosdb.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/outputs/output1
```
