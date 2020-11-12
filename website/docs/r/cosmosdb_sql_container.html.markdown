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

  indexing_policy {
    indexing_mode = "Consistent"

    included_path {
      path = "/*"
    }

    included_path {
      path = "/included/?"
    }

    excluded_path {
      path = "/excluded/?"
    }
  }

  unique_key {
    paths = ["/definition/idlong", "/definition/idshort"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB SQL Container. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB SQL Container is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Account to create the container within. Changing this forces a new resource to be created.

* `database_name` - (Required) The name of the Cosmos DB SQL Database to create the container within. Changing this forces a new resource to be created.

* `partition_key_path` - (Optional) Define a partition key. Changing this forces a new resource to be created.

* `unique_key` - (Optional) One or more `unique_key` blocks as defined below. Changing this forces a new resource to be created.

* `throughput` - (Optional) The throughput of SQL container (RU/s). Must be set in increments of `100`. The minimum value is `400`. This must be set upon container creation otherwise it cannot be updated without a manual terraform destroy-apply.

* `autoscale_settings` - (Optional) An `autoscale_settings` block as defined below. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply. Requires `partition_key_path` to be set.

~> **Note:** Switching between autoscale and manual throughput is not supported via Terraform and must be completed via the Azure Portal and refreshed. 

* `indexing_policy` - (Optional) An `indexing_policy` block as defined below.

* `default_ttl` - (Optional) The default time to live of SQL container. If missing, items are not expired automatically. If present and the value is set to `-1`, it is equal to infinity, and items don’t expire by default. If present and the value is set to some number `n` – items will expire `n` seconds after their last modified time.

---

An `autoscale_settings` block supports the following:

* `max_throughput` - (Optional) The maximum throughput of the SQL container (RU/s). Must be between `4,000` and `1,000,000`. Must be set in increments of `1,000`. Conflicts with `throughput`.

---
A `unique_key` block supports the following:

* `paths` - (Required) A list of paths to use for this unique key.

An `indexing_policy` block supports the following:

* `indexing_mode` - (Optional) Indicates the indexing mode. Possible values include: `Consistent` and `None`. Defaults to `Consistent`.

* `included_path` - (Optional) One or more `included_path` blocks as defined below. Either `included_path` or `excluded_path` must contain the `path` `/*`

* `excluded_path` - (Optional) One or more `excluded_path` blocks as defined below. Either `included_path` or `excluded_path` must contain the `path` `/*`

* `composite_index` - (Optional) One or more `composite_index` blocks as defined below.

An `included_path` block supports the following:

* `path` - Path for which the indexing behaviour applies to.

An `excluded_path` block supports the following:

* `path` - Path that is excluded from indexing.

A `composite_index` block supports the following:

* `index` - One or more `index` blocks as defined below.

An `index` block supports the following:

* `path` - Path for which the indexing behaviour applies to.

* `order` - Order of the index. Possible values are `Ascending` or `Descending`.

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
terraform import azurerm_cosmosdb_sql_container.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlDatabases/database1/containers/container1
```
