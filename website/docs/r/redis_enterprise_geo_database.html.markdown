---
subcategory: "Redis Enterprise"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_enterprise_geo_database"
description: |-
  Manages a Redis Enterprise Active Geo-Replication Database.
---

# azurerm_redis_enterprise_geo_database

Manages a Redis Enterprise Active Geo-Replication Database.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-redisenterprise"
  location = "West Europe"
}

resource "azurerm_redis_enterprise_cluster" "example" {
  name                = "example-geo0"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "Enterprise_E20-4"
}

resource "azurerm_redis_enterprise_cluster" "example1" {
  name                = "example-geo1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "Enterprise_E20-4"
}

resource "azurerm_redis_enterprise_geo_database" "example" {
  name = "default"

  cluster_id                 = azurerm_redis_enterprise_cluster.example.id
  client_protocol            = "Encrypted"
  clustering_policy          = "EnterpriseCluster"
  eviction_policy            = "NoEviction"
  redi_search_module_enabled = true
  redi_search_module_args    = ""
  port                       = 1000

  linked_database_id = [
    "${azurerm_redis_enterprise_cluster.example.id}/databases/default",
    "${azurerm_redis_enterprise_cluster.example1.id}/databases/default",
  ]

  linked_database_group_nickname = "tftestGeoGroup"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Redis Enterprise Geo Database. Currently the acceptable value for this argument is `default`. Defaults to `default`. Changing this forces a new Redis Enterprise Geo Database to be created.

* `cluster_id` - (Required) The resource id of the Redis Enterprise Cluster to deploy this Redis Enterprise Geo Database. Changing this forces a new Redis Enterprise Geo Database to be created.

* `client_protocol` - (Optional) Specifies whether redis clients can connect using TLS-encrypted or plaintext redis protocols. Default is TLS-encrypted. Possible values are `Encrypted` and `Plaintext`. Defaults to `Encrypted`. Changing this forces a new Redis Enterprise Geo Database to be created.

* `clustering_policy` - (Optional) Clustering policy - default is OSSCluster. Specified at create time. Possible values are `EnterpriseCluster` and `OSSCluster`. Defaults to `OSSCluster`. Changing this forces a new Redis Enterprise Geo Database to be created.

* `eviction_policy` - (Optional) Redis eviction policy - default is VolatileLRU. Possible values are `AllKeysLFU`, `AllKeysLRU`, `AllKeysRandom`, `VolatileLRU`, `VolatileLFU`, `VolatileTTL`, `VolatileRandom` and `NoEviction`. Defaults to `VolatileLRU`. Eviction policy must be set to NoEviction when using RediSearch module. Changing this forces a new Redis Enterprise Geo Database to be created.

* `redi_search_module_enabled` - (Optional) Whether to enable Redi-Search module, only `RediSearch` is allowed with geo-replication. Defaults to false. Changing this forces a new Redis Enterprise Geo Database to be created.

* `redi_search_module_args` - (Optional) Configuration options for the module. 

* `port` - (Optional) TCP port of the database endpoint. Specified at create time. Defaults to an available port. Changing this forces a new Redis Enterprise Geo Database to be created.

* `linked_database_id` - (Required) A list of database resources to link with this database with a maximum of 5.

-> **NOTE:** Only the newly created databases can be added to an existing geo-replication group. Existing regular databases or recreated databases cannot be added to the existing geo-replication group.

* `linked_database_group_nickname` - (Optional) Nickname of the group of linked databases. Changing this force a new Redis Enterprise Geo Database to be created.

* `force_unlink_database_id`  - (Optional) A list of linked database resources to be removed from an existing geo-replication group. 

-> **NOTE:** Only the linked databases which exists in the `linked_database_id` list can be added to this list. The corresponding database must be removed from the `linked_database_id` list when adding to the unlinked list. The only recommended operation is to delete after force-unlink and the recommended scenario of force-unlink is region outrage. The database cannot be linked again after force-unlink.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Redis Enterprise Geo Database.

* `primary_access_key` - The Primary Access Key for the Redis Enterprise Geo Database Instance.

* `secondary_access_key` - The Secondary Access Key for the Redis Enterprise Geo Database Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Redis Enterprise Geo Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Enterprise Geo Database.
* `update` - (Defaults to 30 minutes) Used when updating the Redis Enterprise Geo Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Redis Enterprise Geo Database.

## Import

Redis Enterprise Geo Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_enterprise_geo_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/cluster1/databases/database1
```
