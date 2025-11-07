---
subcategory: "Managed Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_redis"
description: |-
  Manages a Managed Redis instance.
---

# azurerm_managed_redis

Manages a [Managed Redis](https://learn.microsoft.com/azure/redis/overview). This resource supersedes [azurerm_redis_enterprise_cluster](azurerm_redis_enterprise_cluster.html) and [azurerm_redis_enterprise_database](azurerm_redis_enterprise_database.html) resources. Please refer to the migration guide for more information on migrating from Redis Enterprise to Managed Redis: [Migrating from Redis Enterprise to Managed Redis](https://learn.microsoft.com/azure/redis/migrate/migrate-overview).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_redis" "example" {
  name                = "example-managed-redis"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Balanced_B3"

  default_database {
    geo_replication_group_name = "myGeoGroup"
  }
}
```

## Example Usage with Customer Managed Key

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_key_vault" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "List",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
      "SetRotationPolicy"
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.example.principal_id

    key_permissions = [
      "Get",
      "WrapKey",
      "UnwrapKey"
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "managedrediscmk"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "unwrapKey", "wrapKey"
  ]
}

resource "azurerm_managed_redis" "example" {
  name                = "example-managed-redis"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Balanced_B3"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  }

  default_database {
    geo_replication_group_name = "myGeoGroup"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Managed Redis instance. Changing this forces a new Managed Redis instance to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed Redis instance should exist. Changing this forces a new Managed Redis instance to be created.

* `location` - (Required) The Azure Region where the Managed Redis instance should exist. Refer to "Redis Cache" on the [product availability documentation](https://azure.microsoft.com/explore/global-infrastructure/products-by-region/table) for supported locations. Changing this forces a new Managed Redis instance to be created.

* `sku_name` - (Required) The features and specification of the Managed Redis instance to deploy. Refer to [the documentation](https://learn.microsoft.com/rest/api/redis/redisenterprisecache/redis-enterprise/create?view=rest-redis-redisenterprisecache-2025-04-01&tabs=HTTP#skuname) for valid values. `Balanced_B3` SKU or higher is required for geo-replication. Changing this forces a new Managed Redis instance to be created.

~> **Note:** `Enterprise_` and `EnterpriseFlash_` prefixed SKUs were previously used by Redis Enterprise, and [not supported by Managed Redis](https://learn.microsoft.com/azure/redis/migrate/migrate-overview).

* `default_database` - (Optional) A `default_database` block as defined below. A Managed Redis instance will not be functional without a database. This block is intentionally optional to allow removal and re-creation of the database for troubleshooting purposes. A default database can be created or deleted in-place, however most properties will trigger an entire cluster replacement if changed.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as defined below.

* `high_availability_enabled` - (Optional) Whether to enable high availability for the Managed Redis instance. Defaults to `true`. Changing this forces a new Managed Redis instance to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `public_network_access` - (Optional) The public network access setting for the Managed Redis instance. Possible values are `Enabled` and `Disabled`. Defaults to `Enabled`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Managed Redis instance.

---

A `default_database` block supports the following:

~> **Note:** Updating the following properties will force a new database to be created, data will be lost and Managed Redis will be unavailable during the operation: `clustering_policy`, `geo_replication_group_name`, and `module`

* `access_keys_authentication_enabled` - (Optional) Whether access key authentication is enabled for the database. Defaults to `false`.

* `client_protocol` - (Optional) Specifies whether redis clients can connect using TLS-encrypted or plaintext redis protocols. Possible values are `Encrypted` and `Plaintext`. Defaults to `Encrypted`.

* `clustering_policy` - (Optional) Clustering policy specified at create time. Possible values are `EnterpriseCluster` and `OSSCluster`. Defaults to `OSSCluster`. Changing this forces a new database to be created, data will be lost and Managed Redis will be unavailable during the operation.

* `eviction_policy` - (Optional) Specifies the Redis eviction policy. Possible values are `AllKeysLFU`, `AllKeysLRU`, `AllKeysRandom`, `VolatileLRU`, `VolatileLFU`, `VolatileTTL`, `VolatileRandom` and `NoEviction`. Defaults to `VolatileLRU`.

* `geo_replication_group_name` - (Optional) The name of the geo-replication group. If provided, a geo-replication group will be created for this database with itself as the only member. Use [`azurerm_managed_redis_database_geo_replication`](azurerm_managed_redis_database_geo_replication.html) resource to manage group membership, linking and unlinking. All databases to be linked have to have the same group name. Refer to the [Managed Redis geo-replication documentation](https://learn.microsoft.com/azure/redis/how-to-active-geo-replication) for more information. Changing this forces a new database to be created, data will be lost and Managed Redis will be unavailable during the operation.

* `module` - (Optional) A `module` block as defined below. Refer to [the modules documentation](https://learn.microsoft.com/azure/redis/redis-modules) to learn more.

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Required) The ID of the key vault key used for encryption. For example: `https://example-vault-name.vault.azure.net/keys/example-key-name/a1b2c3d4`.

* `user_assigned_identity_id` - (Required) The ID of the User Assigned Identity that has access to the Key Vault Key.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Managed Redis instance. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Managed Redis instance.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `module` block supports the following:

* `name` - (Required) The name which should be used for this module. Possible values are `RedisBloom`, `RedisTimeSeries`, `RediSearch` and `RedisJSON`. Changing this forces a new database to be created, data will be lost and Managed Redis will be unavailable during the operation.

* `args` - (Optional) Configuration options for the module (e.g. `ERROR_RATE 0.00 INITIAL_SIZE 400`). Changing this forces a new database to be created, data will be lost and Managed Redis will be unavailable during the operation.

~> **Note:** Only `RediSearch` and `RedisJSON` modules are allowed with geo-replication.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Redis instance.

* `hostname` - DNS name of the cluster endpoint.

---

A `default_database` block exports the following:

* `port` - TCP port of the database endpoint.

* `primary_access_key` - The Primary Access Key for the Managed Redis Database Instance. Only exported if `access_keys_authentication_enabled` is set to `true`.

* `secondary_access_key` - The Secondary Access Key for the Managed Redis Database Instance. Only exported if `access_keys_authentication_enabled` is set to `true`.

---

A `module` block exports the following:

* `version` - Version of the module to be used.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 45 minutes) Used when creating the Managed Redis instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Redis instance.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Redis instance.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Redis instance.

## Import

Managed Redis instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_redis.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cache` - 2025-07-01
