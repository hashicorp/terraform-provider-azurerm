---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_sql_stored_procedure"
description: |-
  Manages a SQL Stored Procedure within a Cosmos DB Account SQL Database.
---

# azurerm_cosmosdb_sql_stored_procedure

Manages a SQL Stored Procedure within a Cosmos DB Account SQL Database.

## Example Usage

```hcl
data "azurerm_cosmosdb_account" "example" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = "tfex-cosmosdb-account-rg"
}

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "tfex-cosmos-db"
  resource_group_name = data.azurerm_cosmosdb_account.example.resource_group_name
  account_name        = data.azurerm_cosmosdb_account.example.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "example" {
  name                = "example-container"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_sql_database.example.name
  partition_key_path  = "/id"
}

resource "azurerm_cosmosdb_sql_stored_procedure" "example" {
  name                = "test-stored-proc"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_sql_database.example.name
  container_name      = azurerm_cosmosdb_sql_container.example.name

  body = <<BODY
  	function () { var context = getContext(); var response = context.getResponse(); response.setBody('Hello, World'); }
BODY
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB SQL Stored Procedure. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB SQL Database is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Account to create the stored procedure within. Changing this forces a new resource to be created.

* `database_name` - (Required) The name of the Cosmos DB SQL Database to create the stored procedure within. Changing this forces a new resource to be created.

* `container_name` - (Required) The name of the Cosmos DB SQL Container to create the stored procedure within. Changing this forces a new resource to be created.

* `body` - (Required) The body of the stored procedure.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cosmos DB SQL Stored Procedure.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CosmosDB SQL Stored Procedure.
* `update` - (Defaults to 30 minutes) Used when updating the CosmosDB SQL Stored Procedure.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB SQL Stored Procedure.
* `delete` - (Defaults to 30 minutes) Used when deleting the CosmosDB SQL Stored Procedure.

## Import

CosmosDB SQL Stored Procedures can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_sql_stored_procedure.db1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlDatabases/db1/containers/c1/storedProcedures/sp1
```
