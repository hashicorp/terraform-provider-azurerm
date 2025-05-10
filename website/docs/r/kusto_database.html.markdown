---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_database"
description: |-
  Manages Kusto / Data Explorer Database
---

# azurerm_kusto_database

Manages a Kusto (also known as Azure Data Explorer) Database

!> **Note:** To mitigate the possibility of accidental data loss it is highly recommended that you use the `prevent_destroy` lifecycle argument in your configuration file for this resource. For more information on the `prevent_destroy` lifecycle argument please see the [terraform documentation](https://developer.hashicorp.com/terraform/tutorials/state/resource-lifecycle#prevent-resource-deletion).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "my-kusto-rg"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "kustocluster"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }
}

resource "azurerm_kusto_database" "database" {
  name                = "my-kusto-database"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.cluster.name

  hot_cache_period   = "P7D"
  soft_delete_period = "P31D"

  # prevent the possibility of accidental data loss
  lifecycle {
    prevent_destroy = true
  }
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

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Kusto Cluster ID.

* `size` - The size of the database in bytes.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kusto Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Database.
* `update` - (Defaults to 1 hour) Used when updating the Kusto Database.
* `delete` - (Defaults to 1 hour) Used when deleting the Kusto Database.

## Import

Kusto Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1
```
