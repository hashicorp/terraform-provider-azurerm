---
subcategory: "Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_cache"
description: |-
  Manages a Redis Cache

---

# azurerm_redis_cache

Manages a Redis Cache.

-> **Note:** Redis version 4 is being retired and no longer supports creating new instances. Version 4 will be removed in a future release. [Redis Version 4 Retirement](https://learn.microsoft.com/azure/azure-cache-for-redis/cache-retired-features#important-upgrade-timelines)

## Example Usage

This example provisions a Standard Redis Cache. Other examples of the `azurerm_redis_cache` resource can be found in [the `./examples/redis-cache` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/redis-cache)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

# NOTE: the Name used for Redis needs to be globally unique
resource "azurerm_redis_cache" "example" {
  name                 = "example-cache"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  capacity             = 2
  family               = "C"
  sku_name             = "Standard"
  non_ssl_port_enabled = false
  minimum_tls_version  = "1.2"

  redis_configuration {
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Redis instance. Changing this forces a new resource to be created.

* `location` - (Required) The location of the resource group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Redis instance. Changing this forces a new resource to be created.

* `capacity` - (Required) The size of the Redis cache to deploy. Valid values for a SKU `family` of C (Basic/Standard) are `0, 1, 2, 3, 4, 5, 6`, and for P (Premium) `family` are `1, 2, 3, 4, 5`.

* `family` - (Required) The SKU family/pricing group to use. Valid values are `C` (for Basic/Standard SKU family) and `P` (for `Premium`)

* `sku_name` - (Required) The SKU of Redis to use. Possible values are `Basic`, `Standard` and `Premium`.

~> **Note:** Downgrading the SKU will force a new resource to be created.

---

* `access_keys_authentication_enabled` - (Optional) Whether access key authentication is enabled? Defaults to `true`. `active_directory_authentication_enabled` must be set to `true` to disable access key authentication.

* `non_ssl_port_enabled` - (Optional) Enable the non-SSL port (6379) - disabled by default.

* `identity` - (Optional) An `identity` block as defined below.

* `minimum_tls_version` - (Optional) The minimum TLS version. Possible values are `1.0`, `1.1` and `1.2`. Defaults to `1.0`.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

* `patch_schedule` - (Optional) A list of `patch_schedule` blocks as defined below.

* `private_static_ip_address` - (Optional) The Static IP Address to assign to the Redis Cache when hosted inside the Virtual Network. This argument implies the use of `subnet_id`. Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this Redis Cache. `true` means this resource could be accessed by both public and private endpoint. `false` means only private endpoint access is allowed. Defaults to `true`.

* `redis_configuration` - (Optional) A `redis_configuration` block as defined below - with some limitations by SKU - defaults/details are shown below.

* `replicas_per_master` - (Optional) Amount of replicas to create per master for this Redis Cache.

~> **Note:** Configuring the number of replicas per master is only available when using the Premium SKU and cannot be used in conjunction with shards.

* `replicas_per_primary` - (Optional) Amount of replicas to create per primary for this Redis Cache. If both `replicas_per_primary` and `replicas_per_master` are set, they need to be equal.

* `redis_version` - (Optional) Redis version. Only major version needed. Possible values are `4` and `6`. Defaults to `6`.

* `tenant_settings` - (Optional) A mapping of tenant settings to assign to the resource.

* `shard_count` - (Optional) *Only available when using the Premium SKU* The number of Shards to create on the Redis Cluster.

* `subnet_id` - (Optional) *Only available when using the Premium SKU* The ID of the Subnet within which the Redis Cache should be deployed. This Subnet must only contain Azure Cache for Redis instances without any other type of resources. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Redis Cache should be located. Changing this forces a new Redis Cache to be created.

-> **Note:** Availability Zones are [in Preview and only supported in several regions at this time](https://docs.microsoft.com/azure/availability-zones/az-overview) - as such you must be opted into the Preview to use this functionality. You can [opt into the Availability Zones Preview in the Azure Portal](https://aka.ms/azenroll).

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Redis Cluster. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Redis Cluster.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `patch_schedule` block supports the following:

* `day_of_week` - (Required) the Weekday name - possible values include `Monday`, `Tuesday`, `Wednesday` etc.

* `start_hour_utc` - (Optional) the Start Hour for maintenance in UTC - possible values range from `0 - 23`.

~> **Note:** The Patch Window lasts for `5` hours from the `start_hour_utc`.

* `maintenance_window` - (Optional) The ISO 8601 timespan which specifies the amount of time the Redis Cache can be updated. Defaults to `PT5H`.

---

A `redis_configuration` block supports the following:

* `aof_backup_enabled` - (Optional) Enable or disable AOF persistence for this Redis Cache. Defaults to `false`.

~> **Note:** `aof_backup_enabled` can only be set when SKU is `Premium`.

* `aof_storage_connection_string_0` - (Optional) First Storage Account connection string for AOF persistence.
* `aof_storage_connection_string_1` - (Optional) Second Storage Account connection string for AOF persistence.

Example usage:

```hcl
redis_configuration {
  aof_backup_enabled              = true
  aof_storage_connection_string_0 = "DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.nc-cruks-storage-account.primary_blob_endpoint};AccountName=${azurerm_storage_account.mystorageaccount.name};AccountKey=${azurerm_storage_account.mystorageaccount.primary_access_key}"
  aof_storage_connection_string_1 = "DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.mystorageaccount.primary_blob_endpoint};AccountName=${azurerm_storage_account.mystorageaccount.name};AccountKey=${azurerm_storage_account.mystorageaccount.secondary_access_key}"
}
```

* `authentication_enabled` - (Optional) If set to `false`, the Redis instance will be accessible without authentication. Defaults to `true`.

-> **Note:** `authentication_enabled` can only be set to `false` if a `subnet_id` is specified; and only works if there aren't existing instances within the subnet with `authentication_enabled` set to `true`.

* `active_directory_authentication_enabled` - (Optional) Enable Microsoft Entra (AAD) authentication. Defaults to `false`.

* `maxmemory_reserved` - (Optional) Value in megabytes reserved for non-cache usage e.g. failover. Defaults are shown below.
* `maxmemory_delta` - (Optional) The max-memory delta for this Redis instance. Defaults are shown below.
* `maxmemory_policy` - (Optional) How Redis will select what to remove when `maxmemory` is reached. Defaults to `volatile-lru`.

* `data_persistence_authentication_method` - (Optional) Preferred auth method to communicate to storage account used for data persistence. Possible values are `SAS` and `ManagedIdentity`.

* `maxfragmentationmemory_reserved` - (Optional) Value in megabytes reserved to accommodate for memory fragmentation. Defaults are shown below.

* `rdb_backup_enabled` - (Optional) Is Backup Enabled? Only supported on Premium SKUs. Defaults to `false`.

-> **Note:** If `rdb_backup_enabled` set to `true`, `rdb_storage_connection_string` must also be set.

* `rdb_backup_frequency` - (Optional) The Backup Frequency in Minutes. Only supported on Premium SKUs. Possible values are: `15`, `30`, `60`, `360`, `720` and `1440`.
* `rdb_backup_max_snapshot_count` - (Optional) The maximum number of snapshots to create as a backup. Only supported for Premium SKUs.
* `rdb_storage_connection_string` - (Optional) The Connection String to the Storage Account. Only supported for Premium SKUs. In the format: `DefaultEndpointsProtocol=https;BlobEndpoint=${azurerm_storage_account.example.primary_blob_endpoint};AccountName=${azurerm_storage_account.example.name};AccountKey=${azurerm_storage_account.example.primary_access_key}`.

~> **Note:** There's a bug in the Redis API where the original storage connection string isn't being returned, which [is being tracked in this issue](https://github.com/Azure/azure-rest-api-specs/issues/3037). In the interim you can use [the `ignore_changes` attribute to ignore changes to this field](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changess) e.g.:

* `storage_account_subscription_id` - (Optional) The ID of the Subscription containing the Storage Account.

```hcl
resource "azurerm_redis_cache" "example" {
  # ...
  ignore_changes = [redis_configuration[0].rdb_storage_connection_string]
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

### Default Redis Configuration Values

| Redis Value                     | Basic        | Standard     | Premium      |
|---------------------------------| ------------ | ------------ | ------------ |
| authentication_enabled          | true         | true         | true         |
| maxmemory_reserved              | 2            | 50           | 200          |
| maxfragmentationmemory_reserved | 2            | 50           | 200          |
| maxmemory_delta                 | 2            | 50           | 200          |
| maxmemory_policy                | volatile-lru | volatile-lru | volatile-lru |

~> **Note:** The `maxmemory_reserved`, `maxmemory_delta` and `maxfragmentationmemory_reserved` settings are only available for Standard and Premium caches. More details are available in the Relevant Links section below.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

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

* [Azure Cache for Redis planning](https://docs.microsoft.com/azure/azure-cache-for-redis/cache-planning-faq)
* [Redis: Available Configuration Settings](https://redis.io/topics/config)

## Timeouts

 The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Redis Cache.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Cache.
* `update` - (Defaults to 3 hours) Used when updating the Redis Cache.
* `delete` - (Defaults to 3 hours) Used when deleting the Redis Cache.

## Import

Redis Cache's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_cache.cache1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redis/cache1
```
