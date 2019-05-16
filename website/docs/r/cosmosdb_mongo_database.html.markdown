---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_mongo_database"
sidebar_current: "docs-azurerm-resource-cosmosdb-mongo-database"
description: |-
  Manages a Mongo Database within a Cosmos DB Account.
---

# azurerm_cosmosdb_mongo_database

Manages a Mongo Database within a Cosmos DB Account.

## Example Usage

```hcl
resource "azurerm_cosmosdb_mongo_database" "db" {
  name                = "tfex-cosmos-mongo-db"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  account_name        = "${azurerm_cosmosdb_account.account.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB Mongo Database. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB Mongo Database is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Mongo Database to create the table within. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - the Cosmos DB Mongo Database ID.

## Import

Cosmos Mongo KeySpace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_mongo_database.db1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/apis/mongodb/databases/db1
```
