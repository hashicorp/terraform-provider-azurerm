---
subcategory: "Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_cache"
sidebar_current: "docs-azurerm-datasource-redis-cache"
description: |-
  Gets information about an existing Azure Redis Cache.

---

# Data Source: azurerm_redis_cache

Use this data source to access information about an existing Redis Cache

# Example Usage
```hcl
data "azurerm_redis_cache" "example" {
  name                = "myrediscache"
  resource_group_name = "redis-cache"
}

output "primary_access_key" {
  value = "${data.azurerm_redis_cache.example.primary_access_key}"
}

output "hostname" {
  value = "${data.azurerm_redis_cache.example.hostname}"
}
```

## Argument Reference

* `name` - The name of the Redis cache

* `resource_group_name` - The name of the resource group the Redis cache instance is located in.

## Attribute Reference

* `id` - The Cache ID.

* `location` - The location of the Redis Cache.

* `capacity` - The size of the Redis Cache deployed.

* `family` - The SKU family/pricing group used. Possible values are `C` (for Basic/Standard SKU family) and `P` (for `Premium`)

* `sku_name` - The SKU of Redis used. Possible values are `Basic`, `Standard` and `Premium`.

* `enable_non_ssl_port` - Whether the SSL port is enabled.

* `minimum_tls_version` - The minimum TLS version.

* `patch_schedule` - A list of `patch_schedule` blocks as defined below - only available for Premium SKU's.

* `private_static_ip_address` The Static IP Address assigned to the Redis Cache when hosted inside the Virtual Network.

* `hostname` - The Hostname of the Redis Instance

* `ssl_port` - The SSL Port of the Redis Instance

* `port` - The non-SSL Port of the Redis Instance

* `primary_access_key` - The Primary Access Key for the Redis Instance

* `secondary_access_key` - The Secondary Access Key for the Redis Instance

* `redis_configuration` - A `redis_configuration` block as defined below.

---

A `patch_schedule` block supports the following (Requires Premium SKU's, attempting to access this value on Basic or Standard SKU's will result in an error):

* `day_of_week` - the Weekday name for the patch item

* `start_hour_utc` - The Start Hour for maintenance in UTC

~> **Note:** The Patch Window lasts for `5` hours from the `start_hour_utc`.

---

A `redis_configuration` block exports the following:

* `enable_authentication` - Specifies if authentication is enabled

* `maxmemory_reserved` - The value in megabytes reserved for non-cache usage e.g. failover

* `maxmemory_delta` - The max-memory delta for this Redis instance.

* `maxmemory_policy` - How Redis will select what to remove when `maxmemory` is reached.

* `maxfragmentationmemory_reserved` - Value in megabytes reserved to accommodate for memory fragmentation.

* `rdb_backup_enabled` - Is Backup Enabled? Only supported on Premium SKU's.

* `rdb_backup_frequency` - The Backup Frequency in Minutes. Only supported on Premium SKU's.

* `rdb_backup_max_snapshot_count` - The maximum number of snapshots that can be created as a backup.

* `rdb_storage_connection_string` - The Connection String to the Storage Account. Only supported for Premium SKU's.

~> **NOTE:** There's a bug in the Redis API where the original storage connection string isn't being returned, which [is being tracked in this issue](https://github.com/Azure/azure-rest-api-specs/issues/3037). In the interim you can use [the `ignore_changes` attribute to ignore changes to this field](https://www.terraform.io/docs/configuration/resources.html#ignore_changes) e.g.:
