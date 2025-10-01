---
subcategory: "Managed Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_redis_database"
description: |-
  Manages a Managed Redis Database.
---

# azurerm_managed_redis_database

Manages a Managed Redis Database.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-managedredis"
  location = "West Europe"
}

resource "azurerm_managed_redis_cluster" "example" {
  name                = "example-managedredis"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "Balanced_B3"
}

resource "azurerm_managed_redis_database" "example" {
  cluster_id        = azurerm_managed_redis_cluster.example.id
  client_protocol   = "Encrypted"
  clustering_policy = "OSSCluster"
  eviction_policy   = "NoEviction"
  port              = 10000

  geo_replication_group_name = "tftestGeoGroup"
}
```

## Arguments Reference

The following arguments are supported:

* `cluster_id` - (Required) The resource id of the Managed Redis Cluster to deploy this Managed Redis Database. Changing this forces a new Managed Redis Database to be created.

* `access_keys_authentication_enabled` - (Optional) Whether access key authentication is enabled for the database. Defaults to `false`. Changing this forces a new Managed Redis Database to be created.

* `client_protocol` - (Optional) Specifies whether redis clients can connect using TLS-encrypted or plaintext redis protocols. Possible values are `Encrypted` and `Plaintext`. Defaults to `Encrypted`. Changing this forces a new Managed Redis Database to be created.

* `clustering_policy` - (Optional) Clustering policy Specified at create time. Possible values are `EnterpriseCluster` and `OSSCluster`. Defaults to `OSSCluster`. Changing this forces a new Managed Redis Database to be created.

* `eviction_policy` - (Optional) Specifies the Redis eviction policy. Possible values are `AllKeysLFU`, `AllKeysLRU`, `AllKeysRandom`, `VolatileLRU`, `VolatileLFU`, `VolatileTTL`, `VolatileRandom` and `NoEviction`. Changing this forces a new Managed Redis Database to be created. Defaults to `VolatileLRU`.

* `geo_replication_group_name` - (Optional) The name of the geo-replication group. If provided, a geo-replication group will be created for this database with itself as the only member. Use [azurerm_managed_redis_database_geo_replication](azurerm_managed_redis_database_geo_replication.html) resource to manage group membership, linking and unlinking. All databases to be linked have to have the same group name. Refer to the [Managed Redis geo-replication documentation](https://learn.microsoft.com/azure/redis/how-to-active-geo-replication) for more information. Changing this forces a new resource to be created.

* `module` - (Optional) A `module` block as defined below. Changing this forces a new resource to be created.

-> **Note:** Only `RediSearch` and `RedisJSON` modules are allowed with geo-replication

<!-- * `linked_database_id` - (Optional) A list of database resources to link with this database with a maximum of 5. Reference the database using the cluster address prefix to avoid cyclic dependency, for example: `${azurerm_managed_redis_cluster.example1.id}/databases/default`.

-> **Note:** Adding an existing database to a geo-replication group will discard all cache data and causes an outage for several minutes. Please refer to [the Microsoft documentation about geo-replication](https://learn.microsoft.com/en-us/azure/redis/how-to-active-geo-replication#remove-from-an-active-geo-replication-group). -->

* `port` - (Optional) TCP port of the database endpoint. Specified at create time. Defaults to an available port. Changing this forces a new Managed Redis Database to be created. Defaults to `10000`.

---

An `module` block exports the following:

* `name` - (Required) The name which should be used for this module. Possible values are `RedisBloom`, `RedisTimeSeries`, `RediSearch` and `RedisJSON`. Changing this forces a new Managed Redis Database to be created.

* `args` - (Optional) Configuration options for the module (e.g. `ERROR_RATE 0.00 INITIAL_SIZE 400`). Changing this forces a new resource to be created. Defaults to `""`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Redis Database.

* `primary_access_key` - The Primary Access Key for the Managed Redis Database Instance.

* `secondary_access_key` - The Secondary Access Key for the Managed Redis Database Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Redis Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Redis Database.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Redis Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Redis Database.

## Import

Managed Redis Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_redis_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/cluster1/databases/database1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cache` - 2025-04-01
