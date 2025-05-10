---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_attached_database_configuration"
description: |-
  Manages Kusto / Data Explorer Attached Database Configuration
---

# azurerm_kusto_attached_database_configuration

Manages a Kusto (also known as Azure Data Explorer) Attached Database Configuration

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "my-kusto-rg"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "follower_cluster" {
  name                = "cluster1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_cluster" "followed_cluster" {
  name                = "cluster2"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "followed_database" {
  name                = "my-followed-database"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.follower_cluster.name
}

resource "azurerm_kusto_database" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.follower_cluster.name
}

resource "azurerm_kusto_attached_database_configuration" "example" {
  name                = "configuration1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.follower_cluster.name
  cluster_resource_id = azurerm_kusto_cluster.followed_cluster.id
  database_name       = azurerm_kusto_database.example.name

  sharing {
    external_tables_to_exclude    = ["ExternalTable2"]
    external_tables_to_include    = ["ExternalTable1"]
    materialized_views_to_exclude = ["MaterializedViewTable2"]
    materialized_views_to_include = ["MaterializedViewTable1"]
    tables_to_exclude             = ["Table2"]
    tables_to_include             = ["Table1"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Kusto Attached Database Configuration to create. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the location of the Kusto Cluster for which the configuration will be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group of the Kusto Cluster for which the configuration will be created. Changing this forces a new resource to be created.

* `cluster_name` - (Required) Specifies the name of the Kusto Cluster for which the configuration will be created. Changing this forces a new resource to be created.

* `cluster_id` - (Required) The resource id of the cluster where the databases you would like to attach reside. Changing this forces a new resource to be created.

* `database_name` - (Required) The name of the database which you would like to attach, use * if you want to follow all current and future databases. Changing this forces a new resource to be created.

* `default_principal_modification_kind` - (Optional) The default principals modification kind. Valid values are: `None` (default), `Replace` and `Union`. Defaults to `None`.

* `sharing` - (Optional) A `sharing` block as defined below.

---

An `sharing` block exports the following:

* `external_tables_to_exclude` - (Optional) List of external tables exclude from the follower database.

* `external_tables_to_include` - (Optional) List of external tables to include in the follower database.

* `materialized_views_to_exclude` - (Optional) List of materialized views exclude from the follower database.

* `materialized_views_to_include` - (Optional) List of materialized views to include in the follower database.

* `tables_to_exclude` - (Optional) List of tables to exclude from the follower database.

* `tables_to_include` - (Optional) List of tables to include in the follower database.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Kusto Attached Database Configuration ID.

* `attached_database_names` - The list of databases from the `cluster_resource_id` which are currently attached to the cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kusto Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Database.
* `update` - (Defaults to 1 hour) Used when updating the Kusto Database.
* `delete` - (Defaults to 1 hour) Used when deleting the Kusto Database.

## Import

Kusto Attached Database Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_attached_database_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1/attachedDatabaseConfigurations/configuration1
```
