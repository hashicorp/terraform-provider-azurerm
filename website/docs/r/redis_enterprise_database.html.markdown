---
subcategory: "Redis Enterprise"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_enterprise_database"
description: |-
  Manages a Redis Enterprise Database.
---

# azurerm_redis_enterprise_database

Manages a Redis Enterprise Database.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-redisenterprise"
  location = "West Europe"
}

resource "azurerm_redis_enterprise_cluster" "example" {
  name                = "example-redisenterprise"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "Enterprise_E20-4"
}

resource "azurerm_redis_enterprise_cluster" "example1" {
  name                = "example-redisenterprise1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "Enterprise_E20-4"
}

resource "azurerm_redis_enterprise_database" "example" {
  name = "default"

  cluster_id        = azurerm_redis_enterprise_cluster.example.id
  client_protocol   = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy   = "NoEviction"
  port              = 10000

  linked_database_id = [
    "${azurerm_redis_enterprise_cluster.example.id}/databases/default",
    "${azurerm_redis_enterprise_cluster.example1.id}/databases/default"
  ]

  linked_database_group_nickname = "tftestGeoGroup"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Optional) The name which should be used for this Redis Enterprise Database. Currently the acceptable value for this argument is `default`. Defaults to `default`. Changing this forces a new Redis Enterprise Database to be created.

* `cluster_id` - (Required) The resource id of the Redis Enterprise Cluster to deploy this Redis Enterprise Database. Changing this forces a new Redis Enterprise Database to be created.

* `client_protocol` - (Optional) Specifies whether redis clients can connect using TLS-encrypted or plaintext redis protocols. Possible values are `Encrypted` and `Plaintext`. Defaults to `Encrypted`. Changing this forces a new Redis Enterprise Database to be created.

* `clustering_policy` - (Optional) Clustering policy Specified at create time. Possible values are `EnterpriseCluster` and `OSSCluster`. Defaults to `OSSCluster`. Changing this forces a new Redis Enterprise Database to be created.

* `eviction_policy` - (Optional) Redis eviction policy possible values are `AllKeysLFU`, `AllKeysLRU`, `AllKeysRandom`, `VolatileLRU`, `VolatileLFU`, `VolatileTTL`, `VolatileRandom` and `NoEviction`. Changing this forces a new Redis Enterprise Database to be created. Defaults to `VolatileLRU`.

* `module` - (Optional) A `module` block as defined below. Changing this forces a new resource to be created.

-> **Note:** Only `RediSearch` and `RedisJSON` modules are allowed with geo-replication

* `linked_database_id` - (Optional) A list of database resources to link with this database with a maximum of 5.

-> **Note:** Only the newly created databases can be added to an existing geo-replication group. Existing regular databases or recreated databases cannot be added to the existing geo-replication group. Any linked database be removed from the list will be forcefully unlinked.The only recommended operation is to delete after force-unlink and the recommended scenario of force-unlink is region outrage. The database cannot be linked again after force-unlink.

* `linked_database_group_nickname` - (Optional) Nickname of the group of linked databases. Changing this force a new Redis Enterprise Geo Database to be created.

* `port` - (Optional) TCP port of the database endpoint. Specified at create time. Defaults to an available port. Changing this forces a new Redis Enterprise Database to be created. Defaults to `10000`.

---

An `module` block exports the following:

* `name` - (Required) The name which should be used for this module. Possible values are `RedisBloom`, `RedisTimeSeries`, `RediSearch` and `RedisJSON`. Changing this forces a new Redis Enterprise Database to be created.

* `args` - (Optional) Configuration options for the module (e.g. `ERROR_RATE 0.00 INITIAL_SIZE 400`). Changing this forces a new resource to be created. Defaults to `""`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Redis Enterprise Database.

* `primary_access_key` - The Primary Access Key for the Redis Enterprise Database Instance.

* `secondary_access_key` - The Secondary Access Key for the Redis Enterprise Database Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Redis Enterprise Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Enterprise Database.
* `update` - (Defaults to 30 minutes) Used when updating the Redis Enterprise Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Redis Enterprise Database.

## Import

Redis Enterprise Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_enterprise_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/cluster1/databases/database1
```
