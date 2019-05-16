---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmos_mongo_database"
sidebar_current: "docs-azurerm-resource-cosmos-mongo-database"
description: |-
  Manages a Cosmos Mongo Database.
---

# azurerm_cosmos_mongo_database

Manages a Cosmos Mongo Database.

## Example Usage

```hcl
resource "azurerm_cosmos_mongo_database" "db" {
  name                = "tfex-cosmos-mongo-db"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  account_name        = "${azurerm_cosmosdb_account.account.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos Mongo Database. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos Mongo Database is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos Mongo Database to create the table within. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The Cosmos Mongo Database ID.

## Import

Cosmos Mongo KeySpace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmos_mongo_database.db1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/apis/mongodb/databases/db1
```
