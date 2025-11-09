---
subcategory: "Managed Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_redis"
description: |-
  Gets information about an existing Managed Redis instance.
---

# Data Source: azurerm_managed_redis

Use this data source to access information about an existing Managed Redis instance.

## Example Usage

```hcl
data "azurerm_managed_redis" "example" {
  name                = "example-managed-redis"
  resource_group_name = "example-resources"
}

output "managed_redis_hostname" {
  value = data.azurerm_managed_redis.example.hostname
}

output "managed_redis_primary_access_key" {
  value       = data.azurerm_managed_redis.example.default_database[0].primary_access_key
  description = "The Managed Redis primary access key."
  sensitive   = true
}

output "managed_redis_secondary_access_key" {
  value       = data.azurerm_managed_redis.example.default_database[0].secondary_access_key
  description = "The Managed Redis secondary access key."
  sensitive   = true
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Redis instance.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed Redis instance exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Redis instance.

* `location` - The Azure Region where the Managed Redis instance exists.

* `sku_name` - The SKU name of the Managed Redis instance.

* `hostname` - The DNS hostname of the Managed Redis instance.

* `high_availability_enabled` - Whether high availability is enabled for the Managed Redis instance.

* `customer_managed_key` - A `customer_managed_key` block as defined below.

* `default_database` - A `default_database` block as defined below.

* `identity` - An `identity` block as defined below.

* `public_network_access` - The public network access setting for the Managed Redis instance.

* `tags` - A mapping of tags assigned to the Managed Redis instance.

* `zones` - A list of Availability Zones in which the Managed Redis instance is located.

---

A `customer_managed_key` block exports the following:

* `key_vault_key_id` - The ID of the key vault key used for encryption.

* `user_assigned_identity_id` - The ID of the User Assigned Identity that has access to the Key Vault Key.

---

A `default_database` block exports the following:

* `access_keys_authentication_enabled` - Whether access key authentication is enabled for the database.

* `client_protocol` - The client protocol used by the database (either `Encrypted` or `Plaintext`).

* `clustering_policy` - The clustering policy used by the database.

* `eviction_policy` - The Redis eviction policy used by the database.

* `geo_replication_group_name` - The name of the geo-replication group.

* `geo_replication_linked_database_ids` - A list of linked database IDs for geo-replication.

* `module` - A list of `module` blocks as defined below.

* `port` - The TCP port of the database endpoint.

* `primary_access_key` - The Primary Access Key for the Managed Redis Database instance.

* `secondary_access_key` - The Secondary Access Key for the Managed Redis Database instance.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity configured on the Managed Redis instance.

* `identity_ids` - A list of User Assigned Managed Identity IDs assigned to the Managed Redis instance.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on the Managed Redis instance.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on the Managed Redis instance.

---

A `module` block exports the following:

* `name` - The name of the Redis module.

* `args` - The configuration options for the module.

* `version` - The version of the module.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Redis instance.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Cache` - 2025-07-01
