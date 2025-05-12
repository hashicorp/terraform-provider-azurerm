---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_mssql"
description: |-
  Manages a Stream Analytics Output to Microsoft SQL Server Database.
---

# azurerm_stream_analytics_output_mssql

Manages a Stream Analytics Output to Microsoft SQL Server Database.

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

resource "azurerm_sql_server" "example" {
  name                         = "example-server"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "dbadmin"
  administrator_login_password = "example-password"
}

resource "azurerm_sql_database" "example" {
  name                             = "exampledb"
  resource_group_name              = azurerm_resource_group.example.name
  location                         = azurerm_resource_group.example.location
  server_name                      = azurerm_sql_server.example.name
  requested_service_objective_name = "S0"
  collation                        = "SQL_LATIN1_GENERAL_CP1_CI_AS"
  max_size_bytes                   = "268435456000"
  create_mode                      = "Default"
}

resource "azurerm_stream_analytics_output_mssql" "example" {
  name                      = "example-output-sql"
  stream_analytics_job_name = data.azurerm_stream_analytics_job.example.name
  resource_group_name       = data.azurerm_stream_analytics_job.example.resource_group_name

  server   = azurerm_sql_server.example.fully_qualified_domain_name
  user     = azurerm_sql_server.example.administrator_login
  password = azurerm_sql_server.example.administrator_login_password
  database = azurerm_sql_database.example.name
  table    = "ExampleTable"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Output. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `server` - (Required) The SQL server url. Changing this forces a new resource to be created.

* `user` - (Optional) Username used to login to the Microsoft SQL Server. Changing this forces a new resource to be created. Required if `authentication_mode` is `ConnectionString`. 

* `database` - (Required) The MS SQL database name where the reference table exists. Changing this forces a new resource to be created.

* `password` - (Optional) Password used together with username, to login to the Microsoft SQL Server. Required if `authentication_mode` is `ConnectionString`.

* `table` - (Required) Table in the database that the output points to. Changing this forces a new resource to be created.

* `max_batch_count` - (Optional) The max batch count to write to the SQL Database. Defaults to `10000`. Possible values are between `1` and `1073741824`.

* `max_writer_count` - (Optional) The max writer count for the SQL Database. Defaults to `1`. Possible values are `0` which bases the writer count on the query partition and `1` which corresponds to a single writer.

* `authentication_mode` - (Optional) The authentication mode for the Stream Output. Possible values are `Msi` and `ConnectionString`. Defaults to `ConnectionString`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Output Microsoft SQL Server Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Output Microsoft SQL Server Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Output Microsoft SQL Server Database.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Output Microsoft SQL Server Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Output Microsoft SQL Server Database.

## Import

Stream Analytics Outputs to Microsoft SQL Server Database can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_mssql.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/outputs/output1
```
