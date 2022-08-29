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

* `name` - (Required) The name of the Kusto Database.

* `resource_group_name` - (Required) The Resource Group where the Kusto Database exists.

* `cluster_name` - (Required) The name of the Kusto Cluster this database is added to.

## Attributes Reference

The following attributes are exported:

* `id` - The Kusto Cluster ID.

* `location` - The Azure Region in which the managed Kusto Database exists.

* `hot_cache_period` - The time the data that should be kept in cache for fast queries as ISO 8601 timespan.

* `soft_delete_period` - The time the data should be kept before it stops being accessible to queries as ISO 8601 timespan.

* `size` - The size of the database in bytes.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Database.
