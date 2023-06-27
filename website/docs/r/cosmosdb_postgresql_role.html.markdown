---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_postgresql_role"
description: |-
  Manages an Azure Cosmos DB for PostgreSQL Role.
---

# azurerm_cosmosdb_postgresql_role

Manages an Azure Cosmos DB for PostgreSQL Role.

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
  node_count                      = 0
}

resource "azurerm_cosmosdb_postgresql_role" "example" {
  name       = "examplerole"
  cluster_id = azurerm_cosmosdb_postgresql_cluster.example.id
  password   = "H@Sh1CoR3!"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Cosmos DB for PostgreSQL Role. Changing this forces a new resource to be created.

* `cluster_id` - (Required) The resource ID of the Azure Cosmos DB for PostgreSQL Cluster. Changing this forces a new resource to be created.

* `password` - (Required) The password of the Azure Cosmos DB for PostgreSQL Role. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Cosmos DB for PostgreSQL Role.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Cosmos DB for PostgreSQL Role.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Cosmos DB for PostgreSQL Role.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Cosmos DB for PostgreSQL Role.

## Import

Azure Cosmos DB for PostgreSQL Roles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_postgresql_role.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/cluster1/roles/role1
```
