---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_database"
description: |-
  Manages Kusto / Data Explorer Database
---

# azurerm_kusto_database

Manages a Kusto (also known as Azure Data Explorer) Database

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "my-kusto-rg"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "kustocluster"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }
}

resource "azurerm_kusto_database" "database" {
  name                = "my-kusto-database"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name

  hot_cache_period   = "P7D"
  soft_delete_period = "P31D"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Kusto Database to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Kusto Database should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Database should exist. Changing this forces a new resource to be created.

* `cluster_name` - (Required) Specifies the name of the Kusto Cluster this database will be added to. Changing this forces a new resource to be created.

* `hot_cache_period` - (Optional) The time the data that should be kept in cache for fast queries as ISO 8601 timespan. Default is unlimited. For more information see: [ISO 8601 Timespan](https://en.wikipedia.org/wiki/ISO_8601#Durations)

* `soft_delete_period` - (Optional) The time the data should be kept before it stops being accessible to queries as ISO 8601 timespan. Default is unlimited. For more information see: [ISO 8601 Timespan](https://en.wikipedia.org/wiki/ISO_8601#Durations)


## Attributes Reference

The following attributes are exported:

* `id` - The Kusto Cluster ID.

* `size` - The size of the database in bytes.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Kusto Database.
* `update` - (Defaults to 60 minutes) Used when updating the Kusto Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Database.
* `delete` - (Defaults to 60 minutes) Used when deleting the Kusto Database.

## Import

Kusto Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1/Databases/database1
```
