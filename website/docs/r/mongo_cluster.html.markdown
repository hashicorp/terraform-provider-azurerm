---
subcategory: "MongoCluster"
layout: "azurerm"
page_title: "Azure Resource Manager: `azurerm_mongo_cluster"
description: |-
  Manages an Azure Cosmos DB for MongoDB vCore.
---

# azurerm_mongo_cluster

Manages an Azure Cosmos DB for MongoDB vCore.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

resource "azurerm_mongo_cluster" "example" {
  name                   = "example-mc"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123"
  shard_count            = "1"
  compute_tier           = "Free"
  high_availability_mode = "Disabled"
  storage_size_in_gb     = "32"
}

```

## Example Usage (Preview feature GeoReplicas)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

resource "azurerm_mongo_cluster" "example" {
  name                   = "example-mc"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_username = "adminTerraform"
  administrator_password = "QAZwsx123"
  shard_count            = "1"
  compute_tier           = "M30"
  high_availability_mode = "ZoneRedundantPreferred"
  storage_size_in_gb     = "64"
  preview_features       = ["GeoReplicas"]
}

resource "azurerm_mongo_cluster" "example_geo_replica" {
  name                = "example-mc-geo"
  resource_group_name = azurerm_resource_group.example.name
  location            = "Central US"
  source_server_id    = azurerm_mongo_cluster.example.id
  source_location     = azurerm_mongo_cluster.example.location
  create_mode         = "GeoReplica"

  lifecycle {
    ignore_changes = ["administrator_username", "high_availability_mode", "preview_features", "shard_count", "storage_size_in_gb", "compute_tier", "version"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for the Azure Cosmos DB for MongoDB vCore. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Azure Cosmos DB for MongoDB vCore. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the Azure Cosmos DB for MongoDB vCore exists. Changing this forces a new resource to be created.

* `administrator_username` - (Optional) The administrator username of the Azure Cosmos DB for MongoDB vCore. Changing this forces a new resource to be created.

* `create_mode` - (Optional) The creation mode for the Azure Cosmos DB for MongoDB vCore. Possibles values are `Default` and `GeoReplica`. Defaults to `Default`. Changing this forces a new resource to be created.

-> **Note** The creation mode `GeoReplica` is currently in preview. It is only available when `preview_features` is set.

* `preview_features` - (Optional) The preview features that can be enabled on the Azure Cosmos DB for MongoDB vCore. Changing this forces a new resource to be created.

* `shard_count` -  (Optional) The Number of shards to provision on the Azure Cosmos DB for MongoDB vCore. Changing this forces a new resource to be created.

* `source_location` - (Optional) The location of the source Azure Cosmos DB for MongoDB vCore. Changing this forces a new resource to be created.

* `source_server_id` - (Optional) The ID of the replication source Azure Cosmos DB for MongoDB vCore. Changing this forces a new resource to be created.

* `administrator_password` - (Optional) The Password associated with the `administrator_username` for the Azure Cosmos DB for MongoDB vCore.

* `compute_tier` - (Optional) The compute tier to assign to the Azure Cosmos DB for MongoDB vCore. Possible values are `Free`, `M25`, `M30`, `M40`, `M50`, `M60` and `M80`.

* `high_availability_mode` - (Optional) The high availability mode for the Azure Cosmos DB for MongoDB vCore. Possibles values are `Disabled` and `ZoneRedundantPreferred`.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for the Azure Cosmos DB for MongoDB vCore. Defaults to `true`.

* `storage_size_in_gb` - (Optional) The size of the data disk space for the Azure Cosmos DB for MongoDB vCore.

* `tags` - (Optional) A mapping of tags to assign to the Azure Cosmos DB for MongoDB vCore.

* `version` - (Optional) The version for the Azure Cosmos DB for MongoDB vCore. Possibles values are `5.0`, `6.0` and `7.0`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Cosmos DB for MongoDB vCore.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Cosmos DB for MongoDB vCore.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Cosmos DB for MongoDB vCore.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Cosmos DB for MongoDB vCore.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Cosmos DB for MongoDB vCore.

## Import

Monitor Azure Active Directory Diagnostic Settings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mongo_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/mongoClusters/myMongoCluster
```
