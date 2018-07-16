---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_cache"
sidebar_current: "docs-azurerm-resource-redis-cache"
description: |-
  Manages a Redis Cache

---

# azurerm_redis_cache

Manages a Redis Cache.

## Example Usage (Basic)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "redis-resources"
  location = "West US"
}

# NOTE: the Name used for Redis needs to be globally unique
resource "azurerm_redis_cache" "test" {
  name                = "tf-redis-basic"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  capacity            = 0
  family              = "C"
  sku_name            = "Basic"
  enable_non_ssl_port = false
}
```

## Example Usage (Standard)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "redis-resources"
  location = "West US"
}

# NOTE: the Name used for Redis needs to be globally unique
resource "azurerm_redis_cache" "test" {
  name                = "tf-redis-standard"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  capacity            = 2
  family              = "C"
  sku_name            = "Standard"
  enable_non_ssl_port = false
}
```

## Example Usage (Premium with Clustering)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "redis-resources"
  location = "West US"
}

# NOTE: the Name used for Redis needs to be globally unique
resource "azurerm_redis_cache" "test" {
  name                = "tf-redis-premium"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  shard_count         = 3

  redis_configuration {
    maxmemory_reserved = 2
    maxmemory_delta    = 2
    maxmemory_policy   = "allkeys-lru"
  }
}
```

## Example Usage (Premium with Backup)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "redis-resources"
  location = "West US"
}

resource "azurerm_storage_account" "test" {
  name                     = "redissa"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

# NOTE: the Name used for Redis needs to be globally unique
resource "azurerm_redis_cache" "test" {
  name                = "tf-redis-pbkup"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  capacity            = 3
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false
  redis_configuration {
    rdb_backup_enabled            = true
    rdb_backup_frequency          = 60
    rdb_backup_max_snapshot_count = 1
    rdb_storage_connection_string = "DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.test.primary_blob_endpoint};AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.primary_access_key}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Redis instance. Changing this forces a
    new resource to be created.

* `location` - (Required) The location of the resource group.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the Redis instance.

* `capacity` - (Required) The size of the Redis cache to deploy. Valid values for a SKU `family` of C (Basic/Standard) are `0, 1, 2, 3, 4, 5, 6`, and for P (Premium) `family` are `1, 2, 3, 4`.

* `family` - (Required) The SKU family to use. Valid values are `C` and `P`, where C = Basic/Standard, P = Premium.

The pricing group for the Redis Family - either "C" or "P" at present.

* `sku_name` - (Required) The SKU of Redis to use - can be either Basic, Standard or Premium.

* `enable_non_ssl_port` - (Optional) Enable the non-SSL port (6789) - disabled by default.

* `patch_schedule` - (Optional) A list of `patch_schedule` blocks as defined below - only available for Premium SKU's.

* `private_static_ip_address` - (Optional) The Static IP Address to assign to the Redis Cache when hosted inside the Virtual Network. Changing this forces a new resource to be created.

* `redis_configuration` - (Required) A `redis_configuration` as defined below - with some limitations by SKU - defaults/details are shown below.

* `shard_count` - (Optional) *Only available when using the Premium SKU* The number of Shards to create on the Redis Cluster.

* `subnet_id` - (Optional) The ID of the Subnet within which the Redis Cache should be deployed. Changing this forces a new resource to be created.

---

* `redis_configuration` supports the following:

* `maxmemory_reserved` - (Optional) Value in megabytes reserved for non-cache usage e.g. failover. Defaults are shown below.
* `maxmemory_delta` - (Optional) The max-memory delta for this Redis instance. Defaults are shown below.
* `maxmemory_policy` - (Optional) How Redis will select what to remove when `maxmemory` is reached. Defaults are shown below.

* `rdb_backup_enabled` - (Optional) Is Backup Enabled? Only supported on Premium SKU's.
* `rdb_backup_frequency` - (Optional) The Backup Frequency in Minutes. Only supported on Premium SKU's. Possible values are: `15`, `30`, `60`, `360`, `720` and `1440`.
* `rdb_backup_max_snapshot_count` - (Optional) The maximum number of snapshots to create as a backup. Only supported for Premium SKU's.
* `rdb_storage_connection_string` - (Optional) The Connection String to the Storage Account. Only supported for Premium SKU's. In the format: `DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.test.primary_blob_endpoint};AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.primary_access_key}`.

~> **NOTE:** There's a bug in the Redis API where the original storage connection string isn't being returned, which [is being tracked in this issue](https://github.com/Azure/azure-rest-api-specs/issues/3037). In the interim you can use [the `ignore_changes` attribute to ignore changes to this field](https://www.terraform.io/docs/configuration/resources.html#ignore_changes) e.g.:

```
resource "azurerm_redis_cache" "test" {
  # ...
  ignore_changes = ["redis_configuration.0.rdb_storage_connection_string"]
}
```

* `notify_keyspace_events` - (Optional) Keyspace notifications allows clients to subscribe to Pub/Sub channels in order to receive events affecting the Redis data set in some way. [Reference](https://redis.io/topics/notifications#configuration)

```hcl
redis_configuration {
  maxmemory_reserve  = 10
  maxmemory_delta    = 2
  maxmemory_policy   = "allkeys-lru"
}
```

## Default Redis Configuration Values
| Redis Value        | Basic        | Standard     | Premium      |
| ------------------ | ------------ | ------------ | ------------ |
| maxmemory_reserved | 2            | 50           | 200          |
| maxmemory_delta    | 2            | 50           | 200          |
| maxmemory_policy   | volatile-lru | volatile-lru | volatile-lru |

_*Important*: The `maxmemory_reserved` and `maxmemory_delta` settings are only available for Standard and Premium caches. More details are available in the Relevant Links section below._

* `patch_schedule` supports the following:

* `day_of_week` (Required) the Weekday name - possible values include `Monday`, `Tuesday`, `Wednesday` etc.
* `start_hour_utc` - (Optional) the Start Hour for maintenance in UTC - possible values range from `0 - 23`.

~> **Note:** The Patch Window lasts for `5` hours from the `start_hour_utc`.

## Attributes Reference

The following attributes are exported:

* `id` - The Route ID.

* `hostname` - The Hostname of the Redis Instance

* `ssl_port` - The SSL Port of the Redis Instance

* `port` - The non-SSL Port of the Redis Instance

* `primary_access_key` - The Primary Access Key for the Redis Instance

* `secondary_access_key` - The Secondary Access Key for the Redis Instance

* `redis_configuration` - A `redis_configuration` block as defined below:

---

A `redis_configuration` block exports the following:

* `maxclients` - Returns the max number of connected clients at the same time.

## Relevant Links
 - [Azure Redis Cache: SKU specific configuration limitations](https://azure.microsoft.com/en-us/documentation/articles/cache-configure/#advanced-settings)
 - [Redis: Available Configuration Settings](http://redis.io/topics/config)

## Import

Redis Cache's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_cache.cache1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/Redis/cache1
```
