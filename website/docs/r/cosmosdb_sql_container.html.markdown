---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_sql_container"
description: |-
  Manages a SQL Container within a Cosmos DB Account.
---

# azurerm_cosmosdb_sql_container

Manages a SQL Container within a Cosmos DB Account.

## Example Usage

```hcl
resource "azurerm_cosmosdb_sql_container" "example" {
  name                = "example-container"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_sql_database.example.name
  partition_key_path  = "/definition/id"
  throughput          = 400

  unique_key {
    paths = ["/definition/idlong", "/definition/idshort"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB SQL Database. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB SQL Database is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Account to create the container within. Changing this forces a new resource to be created.

* `database_name` - (Required) The name of the Cosmos DB SQL Database to create the container within. Changing this forces a new resource to be created.

* `partition_key_path` - (Optional) Define a partition key. Changing this forces a new resource to be created.

* `unique_key` - (Optional) One or more `unique_key` blocks as defined below. Changing this forces a new resource to be created.

* `throughput` - (Optional) The throughput of SQL container (RU/s). Must be set in increments of `100`. The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.

* `default_ttl` - (Optional) The default time to live of SQL container. If missing, items are not expired automatically. If present and the value is set to `-1`, it is equal to infinity, and items don’t expire by default. If present and the value is set to some number `n` – items will expire `n` seconds after their last modified time.

---
A `unique_key` block supports the following:

* `paths` - (Required) A list of paths to use for this unique key.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CosmosDB SQL Container.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CosmosDB SQL Container.
* `update` - (Defaults to 30 minutes) Used when updating the CosmosDB SQL Container.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB SQL Container.
* `delete` - (Defaults to 30 minutes) Used when deleting the CosmosDB SQL Container.

## Import

Cosmos SQL Containers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_sql_container.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/databaseAccounts/account1/apis/sql/databases/database1/containers/container1
```
