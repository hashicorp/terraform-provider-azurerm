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

  sku_name = "EnterpriseFlash_F300-2"
}

resource "azurerm_redis_enterprise_database" "example" {
  name                = "example-database"
  resource_group_name = azurerm_resource_group.example.name

  cluster_id = azurerm_redis_enterprise_cluster.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Redis Enterprise Database. Changing this forces a new Redis Enterprise Database to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Redis Enterprise Database should exist. Changing this forces a new Redis Enterprise Database to be created.

* `cluster_name` - (Required) The name of the RedisEnterprise cluster. Changing this forces a new Redis Enterprise Database to be created.

---

* `client_protocol` - (Optional) Specifies whether redis clients can connect using TLS-encrypted or plaintext redis protocols. Default is TLS-encrypted. Possible values are "Encrypted" and "Plaintext" is allowed.

* `clustering_policy` - (Optional) Clustering policy - default is OSSCluster. Specified at create time. Possible values are "EnterpriseCluster" and "OSSCluster" is allowed.

* `eviction_policy` - (Optional) Redis eviction policy - default is VolatileLRU. Possible values are "AllKeysLFU", "AllKeysLRU", "AllKeysRandom", "VolatileLRU", "VolatileLFU", "VolatileTTL", "VolatileRandom" and "NoEviction" is allowed.

* `module` - (Optional)  A `module` block as defined below.

* `persistence` - (Optional)  A `persistence` block as defined below.

* `port` - (Optional) TCP port of the database endpoint. Specified at create time. Defaults to an available port.

---

An `module` block exports the following:

* `name` - (Required) The name which should be used for this module.

---

* `args` - (Optional) Configuration options for the module, e.g. 'ERROR_RATE 0.00 INITIAL_SIZE 400'.

---

An `persistence` block exports the following:

* `aof_enabled` - (Optional) Should the aof be enabled?

* `aof_frequency` - (Optional) Sets the frequency at which data is written to disk. Possible values are "1s" and "always" is allowed.

* `rdb_enabled` - (Optional) Should the rdb be enabled?

* `rdb_frequency` - (Optional) Sets the frequency at which a snapshot of the database is created. Possible values are "1h", "6h" and "12h" is allowed.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Redis Enterprise Database.

* `resource_state` - Current resource status of the database.

* `type` - The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts".

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Redis Enterprise Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Enterprise Database.
* `update` - (Defaults to 30 minutes) Used when updating the Redis Enterprise Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Redis Enterprise Database.

## Import

Redis Enterprise Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redisenterprise_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/cluster1/databases/database1
```