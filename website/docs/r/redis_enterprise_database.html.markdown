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

resource "azurerm_redis_enterprise_database" "example" {
  name                = "example-database"
  resource_group_name = azurerm_resource_group.example.name

  cluster_id = azurerm_redis_enterprise_cluster.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Redis Enterprise Database. Currently the acceptable value for this argument is `default`. Defaults to `default`. Changing this forces a new Redis Enterprise Database to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Redis Enterprise Database should exist. Changing this forces a new Redis Enterprise Database to be created.

* `cluster_id` - (Required) The resource id of the Redis Enterprise Cluster to deploy this Redis Enterprise Database. Changing this forces a new Redis Enterprise Database to be created.

* `client_protocol` - (Optional) Specifies whether redis clients can connect using TLS-encrypted or plaintext redis protocols. Default is TLS-encrypted. Possible values are `Encrypted` and `Plaintext`. Defaults to `Encrypted`. Changing this forces a new Redis Enterprise Database to be created.

* `clustering_policy` - (Optional) Clustering policy - default is OSSCluster. Specified at create time. Possible values are `EnterpriseCluster` and `OSSCluster`. Defaults to `OSSCluster`. Changing this forces a new Redis Enterprise Database to be created.

* `eviction_policy` - (Optional) Redis eviction policy - default is VolatileLRU. Possible values are `AllKeysLFU`, `AllKeysLRU`, `AllKeysRandom`, `VolatileLRU`, `VolatileLFU`, `VolatileTTL`, `VolatileRandom` and `NoEviction`. Defaults to `VolatileLRU`. Changing this forces a new Redis Enterprise Database to be created.

* `module` - (Optional)  A `module` block as defined below.

* `port` - (Optional) TCP port of the database endpoint. Specified at create time. Defaults to an available port. Changing this forces a new Redis Enterprise Database to be created.

---

An `module` block exports the following:

* `name` - (Required) The name which should be used for this module. Possible values are `RediSearch`, `RedisBloom` and `RedisTimeSeries`. Changing this forces a new Redis Enterprise Database to be created.

* `args` - (Optional) Configuration options for the module (e.g. `ERROR_RATE 0.00 INITIAL_SIZE 400`).

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Redis Enterprise Database.

* `primary_access_key` - The Primary Access Key for the Redis Enterprise Database Instance.

* `secondary_access_key` - The Secondary Access Key for the Redis Enterprise Database Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Redis Enterprise Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Enterprise Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Redis Enterprise Database.

## Import

Redis Enterprise Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redisenterprise_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/cluster1/databases/database1
```
