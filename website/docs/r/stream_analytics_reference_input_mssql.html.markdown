---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_reference_input_mssql"
description: |-
  Manages a Stream Analytics Reference Input from MS SQL.
---

# azurerm_stream_analytics_reference_input_mssql

Manages a Stream Analytics Reference Input from MS SQL. Reference data (also known as a lookup table) is a finite data set that is static or slowly changing in nature, used to perform a lookup or to correlate with your data stream. Learn more [here](https://docs.microsoft.com/azure/stream-analytics/stream-analytics-use-reference-data#azure-sql-database).

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

resource "azurerm_mssql_server" "example" {
  name                         = "example-sqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "admin"
  administrator_login_password = "password"
}

resource "azurerm_mssql_database" "example" {
  name      = "example-db"
  server_id = azurerm_mssql_server.example.id
}

resource "azurerm_stream_analytics_reference_input_mssql" "example" {
  name                      = "example-reference-input"
  resource_group_name       = data.azurerm_stream_analytics_job.example.resource_group_name
  stream_analytics_job_name = data.azurerm_stream_analytics_job.example.name
  server                    = azurerm_mssql_server.example.fully_qualified_domain_name
  database                  = azurerm_mssql_database.example.name
  username                  = "exampleuser"
  password                  = "examplepassword"
  refresh_type              = "RefreshPeriodicallyWithFull"
  refresh_interval_duration = "00:20:00"
  full_snapshot_query       = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Reference Input MS SQL data. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job should exist. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `server` - (Required) The fully qualified domain name of the MS SQL server.

* `database` - (Required) The MS SQL database name where the reference data exists.

* `username` - (Required) The username to connect to the MS SQL database.

* `password` - (Required) The username to connect to the MS SQL database.

* `refresh_type` - (Required) Defines whether and how the reference data should be refreshed. Accepted values are `Static`, `RefreshPeriodicallyWithFull` and `RefreshPeriodicallyWithDelta`.

* `refresh_interval_duration` - (Optional) The frequency in `hh:mm:ss` with which the reference data should be retrieved from the MS SQL database e.g. `00:20:00` for every 20 minutes. Must be set when `refresh_type` is `RefreshPeriodicallyWithFull` or `RefreshPeriodicallyWithDelta`.

* `full_snapshot_query` - (Required) The query used to retrieve the reference data from the MS SQL database.

* `delta_snapshot_query` - (Optional) The query used to retrieve incremental changes in the reference data from the MS SQL database. Cannot be set when `refresh_type` is `Static`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Stream Analytics.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics.

## Import

Stream Analytics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_reference_input_mssql.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingjobs/job1/inputs/input1
```
