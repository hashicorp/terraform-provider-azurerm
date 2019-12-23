---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_sql_database"
sidebar_current: "docs-azurerm-resource-cosmosdb-sql-database"
description: |-
  Manages a SQL Database within a Cosmos DB Account.
---

# azurerm_cosmosdb_sql_database

Manages a SQL Database within a Cosmos DB Account.

## Example Usage

```hcl
data "azurerm_cosmosdb_account" "example" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = "tfex-cosmosdb-account-rg"
}

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "tfex-cosmos-mongo-db"
  resource_group_name = "${data.azurerm_cosmosdb_account.example.resource_group_name}"
  account_name        = "${data.azurerm_cosmosdb_account.example.name}"
  throughput          = 400
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB SQL Database. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB SQL Database is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB SQL Database to create the table within. Changing this forces a new resource to be created.

* `throughput` - (Optional) The throughput of SQL database (RU/s). Must be set in increments of `100`. The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.


## Attributes Reference

The following attributes are exported:

* `id` - the Cosmos DB SQL Database ID.

## Import

Cosmos SQL Database can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_sql_database.db1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/apis/sql/databases/db1
```
