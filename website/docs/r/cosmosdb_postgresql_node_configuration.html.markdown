---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_postgresql_node_configuration"
description: |-
  Sets a Node Configuration value on Azure Cosmos DB for PostgreSQL Cluster.
---

# azurerm_cosmosdb_postgresql_node_configuration

Sets a Node Configuration value on Azure Cosmos DB for PostgreSQL Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cosmosdb_postgresql_cluster" "example" {
  name                            = "examplecluster"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  administrator_login_password    = "H@Sh1CoR3!"
  coordinator_storage_quota_in_mb = 131072
  coordinator_vcore_count         = 2
  node_count                      = 2
  node_storage_quota_in_mb        = 131072
  node_vcores                     = 2
}

resource "azurerm_cosmosdb_postgresql_node_configuration" "example" {
  name       = "array_nulls"
  cluster_id = azurerm_cosmosdb_postgresql_cluster.example.id
  value      = "on"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Node Configuration on Azure Cosmos DB for PostgreSQL Cluster. Changing this forces a new resource to be created.

* `cluster_id` - (Required) The resource ID of the Azure Cosmos DB for PostgreSQL Cluster where we want to change configuration. Changing this forces a new resource to be created.

* `value` - (Required) The value of the Node Configuration on Azure Cosmos DB for PostgreSQL Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Node Configuration on Azure Cosmos DB for PostgreSQL Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Node Configuration on Azure Cosmos DB for PostgreSQL Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Node Configuration on Azure Cosmos DB for PostgreSQL Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Node Configuration on Azure Cosmos DB for PostgreSQL Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Node Configuration on Azure Cosmos DB for PostgreSQL Cluster.

## Import

Node Configurations on Azure Cosmos DB for PostgreSQL Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_postgresql_node_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/cluster1/nodeConfigurations/array_nulls
```
