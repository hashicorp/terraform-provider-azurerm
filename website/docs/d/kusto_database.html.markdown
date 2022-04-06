---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_database"
description: |-
  Manages Kusto / Data Explorer Database
---

# Data Source: azurerm_kusto_database

Use this data source to access information about an existing Kusto Database

## Example Usage

```hcl
data "azurerm_kusto_database" "example" {
  name                = "my-kusto-database"
  resource_group_name = "test_resource_group"
  cluster_name        = "test_cluster"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Kusto Database to create.

* `resource_group_name` - Specifies the Resource Group where the Kusto Database should exist.

* `cluster_name` - Specifies the name of the Kusto Cluster this database will be added to.

## Attributes Reference

The following attributes are exported:

* `id` - The Kusto Cluster ID.

* `hot_cache_period` - The time the data that should be kept in cache for fast queries as ISO 8601 timespan.

* `soft_delete_period` - The time the data should be kept before it stops being accessible to queries as ISO 8601 timespan.

* `size` - The size of the database in bytes.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Database.
