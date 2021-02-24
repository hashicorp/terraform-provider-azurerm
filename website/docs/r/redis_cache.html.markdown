---
subcategory: "Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_cache"
description: |-
  Manages a Redis Cache

---

# azurerm_redis_cache

Manages a Redis Cache.

## Example Usage

This example provisions a Standard Redis Cache. Other examples of the `azurerm_redis_cache` resource can be found in [the `./examples/redis-cache` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/redis-cache)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

# NOTE: the Name used for Redis needs to be globally unique
resource "azurerm_redis_cache" "example" {
  name                = "example-cache"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  capacity            = 2
  family              = "C"
  sku_name            = "Standard"
  enable_non_ssl_port = false
  minimum_tls_version = "1.2"

  redis_configuration {
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

* `family` - (Required) The SKU family/pricing group to use. Valid values are `C` (for Basic/Standard SKU family) and `P` (for `Premium`)

* `sku_name` - (Required) The SKU of Redis to use. Possible values are `Basic`, `Standard` and `Premium`.

---

* `enable_non_ssl_port` - (Optional) Enable the non-SSL port (6379) - disabled by default.

* `minimum_tls_version` - (Optional) The minimum TLS version.  Defaults to `1.0`.

* `patch_schedule` - (Optional) A list of `patch_schedule` blocks as defined below.

* `private_static_ip_address` - (Optional) The Static IP Address to assign to the Redis Cache when hosted inside the Virtual Network. Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this Redis Cache. `true` means this resource could be accessed by both public and private endpoint. `false` means only private endpoint access is allowed. Defaults to `true`.

* `redis_configuration` - (Optional) A `redis_configuration` as defined below - with some limitations by SKU - defaults/details are shown below.

* `shard_count` - (Optional) *Only available when using the Premium SKU* The number of Shards to create on the Redis Cluster.

* `subnet_id` - (Optional) *Only available when using the Premium SKU* The ID of the Subnet within which the Redis Cache should be deployed. This Subnet must only contain Azure Cache for Redis instances without any other type of resources. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `zones` - (Optional) A list of a one or more Availability Zones, where the Redis Cache should be allocated.

-> **Please Note**: Availability Zones are [in Preview and only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview) - as such you must be opted into the Preview to use this functionality. You can [opt into the Availability Zones Preview in the Azure Portal](http://aka.ms/azenroll).

---

A `redis_configuration` block supports the following:

* `enable_authentication` - (Optional) If set to `false`, the Redis instance will be accessible without authentication. Defaults to `true`.

-> **NOTE:** `enable_authentication` can only be set to `false` if a `subnet_id` is specified; and only works if there aren't existing instances within the subnet with `enable_authentication` set to `true`.

* `maxmemory_reserved` - (Optional) Value in megabytes reserved for non-cache usage e.g. failover. Defaults are shown below.
* `maxmemory_delta` - (Optional) The max-memory delta for this Redis instance. Defaults are shown below.
* `maxmemory_policy` - (Optional) How Redis will select what to remove when `maxmemory` is reached. Defaults are shown below.

* `maxfragmentationmemory_reserved` - (Optional) Value in megabytes reserved to accommodate for memory fragmentation. Defaults are shown below.

* `rdb_backup_enabled` - (Optional) Is Backup Enabled? Only supported on Premium SKU's.

-> **NOTE:** If `rdb_backup_enabled` set to `true`, `rdb_storage_connection_string` must also be set.

* `rdb_backup_frequency` - (Optional) The Backup Frequency in Minutes. Only supported on Premium SKU's. Possible values are: `15`, `30`, `60`, `360`, `720` and `1440`.
* `rdb_backup_max_snapshot_count` - (Optional) The maximum number of snapshots to create as a backup. Only supported for Premium SKU's.
* `rdb_storage_connection_string` - (Optional) The Connection String to the Storage Account. Only supported for Premium SKU's. In the format: `DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.example.primary_blob_endpoint};AccountName=${azurerm_storage_account.example.name};AccountKey=${azurerm_storage_account.example.primary_access_key}`.

~> **NOTE:** There's a bug in the Redis API where the original storage connection string isn't being returned, which [is being tracked in this issue](https://github.com/Azure/azure-rest-api-specs/issues/3037). In the interim you can use [the `ignore_changes` attribute to ignore changes to this field](https://www.terraform.io/docs/configuration/resources.html#ignore_changes) e.g.:

```
resource "azurerm_redis_cache" "example" {
  # ...
  ignore_changes = [redis_configuration.0.rdb_storage_connection_string]
}
```

* `notify_keyspace_events` - (Optional) Keyspace notifications allows clients to subscribe to Pub/Sub channels in order to receive events affecting the Redis data set in some way. [Reference](https://redis.io/topics/notifications#configuration)

```hcl
redis_configuration {
  maxmemory_reserved = 10
  maxmemory_delta    = 2
  maxmemory_policy   = "allkeys-lru"
}
```

## Default Redis Configuration Values

| Redis Value                     | Basic        | Standard     | Premium      |
| ------------------------------- | ------------ | ------------ | ------------ |
| enable_authentication           | true         | true         | true         |
| maxmemory_reserved              | 2            | 50           | 200          |
| maxfragmentationmemory_reserved | 2            | 50           | 200          |
| maxmemory_delta                 | 2            | 50           | 200          |
| maxmemory_policy                | volatile-lru | volatile-lru | volatile-lru |

~> **NOTE:** The `maxmemory_reserved`, `maxmemory_delta` and `maxfragmentationmemory-reserved` settings are only available for Standard and Premium caches. More details are available in the Relevant Links section below._

---

A `patch_schedule` block supports the following:

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

* `primary_connection_string` - The primary connection string of the Redis Instance.

* `secondary_connection_string` - The secondary connection string of the Redis Instance.

* `redis_configuration` - A `redis_configuration` block as defined below:

---

A `redis_configuration` block exports the following:

* `maxclients` - Returns the max number of connected clients at the same time.

## Relevant Links
 - [Azure Redis Cache: SKU specific configuration limitations](https://azure.microsoft.com/en-us/documentation/articles/cache-configure/#advanced-settings)
 - [Redis: Available Configuration Settings](http://redis.io/topics/config)

## Timeouts

 The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

 * `create` - (Defaults to 90 minutes) Used when creating the Redis Cache.
 * `update` - (Defaults to 90 minutes) Used when updating the Redis Cache.
 * `read` - (Defaults to 5 minutes) Used when retrieving the Redis Cache.
 * `delete` - (Defaults to 90 minutes) Used when deleting the Redis Cache.

## Import

Redis Cache's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_cache.cache1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/Redis/cache1
```
