---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_gremlin_graph"
description: |-
  Manages a Gremlin Graph within a Cosmos DB Account.
---

# azurerm_cosmosdb_gremlin_graph

Manages a Gremlin Graph within a Cosmos DB Account.

## Example Usage

```hcl
data "azurerm_cosmosdb_account" "example" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = "tfex-cosmosdb-account-rg"
}

resource "azurerm_cosmosdb_gremlin_database" "example" {
  name                = "tfex-cosmos-gremlin-db"
  resource_group_name = data.azurerm_cosmosdb_account.example.resource_group_name
  account_name        = data.azurerm_cosmosdb_account.example.name
}

resource "azurerm_cosmosdb_gremlin_graph" "example" {
  name                = "tfex-cosmos-gremlin-graph"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_gremlin_database.example.name
  partition_key_path  = "/Example"
  throughput          = 400

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }

  conflict_resolution_policy {
    mode                     = "LastWriterWins"
    conflict_resolution_path = "/_ts"
  }

  unique_key {
    paths = ["/definition/id1", "/definition/id2"]
  }
}
```

-> **NOTE:** The CosmosDB Account needs to have the `EnableGremlin` capability enabled to use this resource - which can be done by adding this to the `capabilities` list within the `azurerm_cosmosdb_account` resource.

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB Gremlin Graph. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB Gremlin Graph is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the CosmosDB Account to create the Gremlin Graph within. Changing this forces a new resource to be created.

* `database_name` - (Required) The name of the Cosmos DB Graph Database in which the Cosmos DB Gremlin Graph is created. Changing this forces a new resource to be created.

* `partition_key_path` - (Required) Define a partition key. Changing this forces a new resource to be created.

* `partition_key_version` - (Optional) Define a partition key version. Changing this forces a new resource to be created. Possible values are `1 `and `2`. This should be set to `2` in order to use large partition keys.

* `throughput` - (Optional) The throughput of the Gremlin graph (RU/s). Must be set in increments of `100`. The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.

* `default_ttl` - (Optional) The default time to live (TTL) of the Gremlin graph. If the value is missing or set to "-1", items don’t expire.

* `autoscale_settings` - (Optional) An `autoscale_settings` block as defined below. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply. Requires `partition_key_path` to be set.

~> **Note:** Switching between autoscale and manual throughput is not supported via Terraform and must be completed via the Azure Portal and refreshed. 

* `index_policy` - (Required) The configuration of the indexing policy. One or more `index_policy` blocks as defined below. Changing this forces a new resource to be created.

* `conflict_resolution_policy` - (Optional)  A `conflict_resolution_policy` blocks as defined below.

* `unique_key` (Optional) One or more `unique_key` blocks as defined below. Changing this forces a new resource to be created.

---

An `autoscale_settings` block supports the following:

* `max_throughput` - (Optional) The maximum throughput of the Gremlin graph (RU/s). Must be between `4,000` and `1,000,000`. Must be set in increments of `1,000`. Conflicts with `throughput`.

---

An `index_policy` block supports the following:

* `automatic` - (Optional) Indicates if the indexing policy is automatic. Defaults to `true`.

* `indexing_mode` - (Required) Indicates the indexing mode. Possible values include: `Consistent`, `Lazy`, `None`.

* `included_paths` - (Optional) List of paths to include in the indexing. Required if `indexing_mode` is `Consistent` or `Lazy`.

* `excluded_paths` - (Optional) List of paths to exclude from indexing. Required if `indexing_mode` is `Consistent` or `Lazy`.

* `composite_index` - (Optional) One or more `composite_index` blocks as defined below.

* `spatial_index` - (Optional) One or more `spatial_index` blocks as defined below.

---

A `spatial_index` block supports the following:

* `path` - (Required) Path for which the indexing behaviour applies to. According to the service design, all spatial types including `LineString`, `MultiPolygon`, `Point`, and `Polygon` will be applied to the path. 

---

An `conflict_resolution_policy` block supports the following:

* `mode` - (Required) Indicates the conflict resolution mode. Possible values include: `LastWriterWins`, `Custom`.

* `conflict_resolution_path` - (Optional) The conflict resolution path in the case of LastWriterWins mode.

* `conflict_resolution_procedure` - (Optional) The procedure to resolve conflicts in the case of custom mode.

---

An `unique_key` block supports the following:

* `paths` - (Required) A list of paths to use for this unique key.

---

A `composite_index` block supports the following:

* `index` - One or more `index` blocks as defined below.

---

An `index` block supports the following:

* `path` - Path for which the indexing behaviour applies to.

* `order` - Order of the index. Possible values are `Ascending` or `Descending`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CosmosDB Gremlin Graph.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CosmosDB Gremlin Graph.
* `update` - (Defaults to 30 minutes) Used when updating the CosmosDB Gremlin Graph.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB Gremlin Graph.
* `delete` - (Defaults to 30 minutes) Used when deleting the CosmosDB Gremlin Graph.

## Import

Cosmos Gremlin Graphs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_gremlin_graph.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/gremlinDatabases/db1/graphs/graphs1
```
