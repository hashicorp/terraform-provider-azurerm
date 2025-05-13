---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_account"
description: |-
  Manages a CosmosDB (formally DocumentDB) Account.
---

# azurerm_cosmosdb_account

Manages a CosmosDB (formally DocumentDB) Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "random_integer" "ri" {
  min = 10000
  max = 99999
}

resource "azurerm_cosmosdb_account" "db" {
  name                = "tfex-cosmos-db-${random_integer.ri.result}"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  automatic_failover_enabled = true

  capabilities {
    name = "EnableAggregationPipeline"
  }

  capabilities {
    name = "mongoEnableDocLevelTTL"
  }

  capabilities {
    name = "MongoDBv3.4"
  }

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 100000
  }

  geo_location {
    location          = "eastus"
    failover_priority = 1
  }

  geo_location {
    location          = "westus"
    failover_priority = 0
  }
}
```

## User Assigned Identity Example Usage

```hcl
resource "azurerm_user_assigned_identity" "example" {
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  name                = "example-resource"
}

resource "azurerm_cosmosdb_account" "example" {
  name                  = "example-resource"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  default_identity_type = join("=", ["UserAssignedIdentity", azurerm_user_assigned_identity.example.id])
  offer_type            = "Standard"
  kind                  = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = "westus"
    failover_priority = 0
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the CosmosDB Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the CosmosDB Account is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `minimal_tls_version` - (Optional) Specifies the minimal TLS version for the CosmosDB account. Possible values are: `Tls`, `Tls11`, and `Tls12`. Defaults to `Tls12`.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more details.

* `offer_type` - (Required) Specifies the Offer Type to use for this CosmosDB Account; currently, this can only be set to `Standard`.

* `analytical_storage` - (Optional) An `analytical_storage` block as defined below.

* `capacity` - (Optional) A `capacity` block as defined below.

* `create_mode` - (Optional) The creation mode for the CosmosDB Account. Possible values are `Default` and `Restore`. Changing this forces a new resource to be created.

~> **Note:** `create_mode` can only be defined when the `backup.type` is set to `Continuous`.

* `default_identity_type` - (Optional) The default identity for accessing Key Vault. Possible values are `FirstPartyIdentity`, `SystemAssignedIdentity` or `UserAssignedIdentity`. Defaults to `FirstPartyIdentity`.

~> **Note:** When `default_identity_type` is a `UserAssignedIdentity` it must include the User Assigned Identity ID in the following format: `UserAssignedIdentity=/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{userAssignedIdentityName}`.

* `kind` - (Optional) Specifies the Kind of CosmosDB to create - possible values are `GlobalDocumentDB`, `MongoDB` and `Parse`. Defaults to `GlobalDocumentDB`. Changing this forces a new resource to be created.

* `consistency_policy` - (Required) Specifies one `consistency_policy` block as defined below, used to define the consistency policy for this CosmosDB account.

* `geo_location` - (Required) Specifies a `geo_location` resource, used to define where data should be replicated with the `failover_priority` 0 specifying the primary location. Value is a `geo_location` block as defined below.

* `ip_range_filter` - (Optional) A set of IP addresses or IP address ranges in CIDR form to be included as the allowed list of client IPs for a given database account. For example `["55.0.1.0/24", "55.0.2.0/24"]`.

~> **Note:** To enable the "Allow access from the Azure portal" behavior, you should add the IP addresses provided by the [documentation](https://docs.microsoft.com/azure/cosmos-db/how-to-configure-firewall#allow-requests-from-the-azure-portal) to this list.

~> **Note:** To enable the "Accept connections from within public Azure datacenters" behavior, you should add `0.0.0.0` to the list, see the [documentation](https://docs.microsoft.com/azure/cosmos-db/how-to-configure-firewall#allow-requests-from-global-azure-datacenters-or-other-sources-within-azure) for more details.

* `free_tier_enabled` - (Optional) Enable the Free Tier pricing option for this Cosmos DB account. Defaults to `false`. Changing this forces a new resource to be created.

* `analytical_storage_enabled` - (Optional) Enable Analytical Storage option for this Cosmos DB account. Defaults to `false`. Enabling and then disabling analytical storage forces a new resource to be created.

* `automatic_failover_enabled` - (Optional) Enable automatic failover for this Cosmos DB account.

* `partition_merge_enabled` - (Optional) Is partition merge on the Cosmos DB account enabled? Defaults to `false`.

* `burst_capacity_enabled` - (Optional) Enable burst capacity for this Cosmos DB account. Defaults to `false`.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this CosmosDB account. Defaults to `true`.

* `capabilities` - (Optional) The capabilities which should be enabled for this Cosmos DB account. Value is a `capabilities` block as defined below.

* `is_virtual_network_filter_enabled` - (Optional) Enables virtual network filtering for this Cosmos DB account.

* `key_vault_key_id` - (Optional) A versionless Key Vault Key ID for CMK encryption. Changing this forces a new resource to be created.

~> **Note:** When referencing an `azurerm_key_vault_key` resource, use `versionless_id` instead of `id`

~> **Note:** In order to use a `Custom Key` from Key Vault for encryption you must grant Azure Cosmos DB Service access to your key vault. For instructions on how to configure your Key Vault correctly please refer to the [product documentation](https://docs.microsoft.com/azure/cosmos-db/how-to-setup-cmk#add-an-access-policy-to-your-azure-key-vault-instance)

* `managed_hsm_key_id` - (Optional) A versionless Managed HSM Key ID for CMK encryption. Changing this forces a new resource to be created.

~> **Note:** When referencing an `azurerm_key_vault_managed_hardware_security_module_key` resource, use `id` instead of `versioned_id`

~> **Note:** In order to use a `Custom Key` from Managed HSM for encryption you must grant Azure Cosmos DB Service access to your Managed HSM. For instructions on how to configure your Key Vault correctly please refer to the [product documentation](https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-setup-customer-managed-keys-mhsm)

* `virtual_network_rule` - (Optional) Specifies a `virtual_network_rule` block as defined below, used to define which subnets are allowed to access this CosmosDB account.

* `multiple_write_locations_enabled` - (Optional) Enable multiple write locations for this Cosmos DB account.

* `access_key_metadata_writes_enabled` - (Optional) Is write operations on metadata resources (databases, containers, throughput) via account keys enabled? Defaults to `true`.

* `mongo_server_version` - (Optional) The Server Version of a MongoDB account. Possible values are `7.0`, `6.0`, `5.0`, `4.2`, `4.0`, `3.6`, and `3.2`.

* `network_acl_bypass_for_azure_services` - (Optional) If Azure services can bypass ACLs. Defaults to `false`.

* `network_acl_bypass_ids` - (Optional) The list of resource Ids for Network Acl Bypass for this Cosmos DB account.

* `local_authentication_disabled` - (Optional) Disable local authentication and ensure only MSI and AAD can be used exclusively for authentication. Defaults to `false`. Can be set only when using the SQL API.

* `backup` - (Optional) A `backup` block as defined below.

* `cors_rule` - (Optional) A `cors_rule` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `restore` - (Optional) A `restore` block as defined below.

~> **Note:** `restore` should be set when `create_mode` is `Restore`.

---

The `consistency_policy` block Configures the database consistency and supports the following:

* `consistency_level` - (Required) The Consistency Level to use for this CosmosDB Account - can be either `BoundedStaleness`, `Eventual`, `Session`, `Strong` or `ConsistentPrefix`.

* `max_interval_in_seconds` - (Optional) When used with the Bounded Staleness consistency level, this value represents the time amount of staleness (in seconds) tolerated. The accepted range for this value is `5` - `86400` (1 day). Defaults to `5`. Required when `consistency_level` is set to `BoundedStaleness`.

* `max_staleness_prefix` - (Optional) When used with the Bounded Staleness consistency level, this value represents the number of stale requests tolerated. The accepted range for this value is `10` – `2147483647`. Defaults to `100`. Required when `consistency_level` is set to `BoundedStaleness`.

~> **Note:** `max_interval_in_seconds` and `max_staleness_prefix` can only be set to values other than default when the `consistency_level` is set to `BoundedStaleness`.

---

The `geo_location` block Configures the geographic locations the data is replicated to and supports the following:

* `location` - (Required) The name of the Azure region to host replicated data.

* `failover_priority` - (Required) The failover priority of the region. A failover priority of `0` indicates a write region. The maximum value for a failover priority = (total number of regions - 1). Failover priority values must be unique for each of the regions in which the database account exists. Changing this causes the location to be re-provisioned and cannot be changed for the location with failover priority `0`.

* `zone_redundant` - (Optional) Should zone redundancy be enabled for this region? Defaults to `false`.

---

A `capabilities` block Configures the capabilities to be enabled for this Cosmos DB account:

* `name` - (Required) The capability to enable - Possible values are `AllowSelfServeUpgradeToMongo36`, `DeleteAllItemsByPartitionKey`, `DisableRateLimitingResponses`, `EnableAggregationPipeline`, `EnableCassandra`, `EnableGremlin`, `EnableMongo`, `EnableMongo16MBDocumentSupport`, `EnableMongoRetryableWrites`, `EnableMongoRoleBasedAccessControl`, `EnableNoSQLVectorSearch`, `EnableNoSQLFullTextSearch`, `EnablePartialUniqueIndex`,  `EnableServerless`, `EnableTable`, `EnableTtlOnCustomPath`, `EnableUniqueCompoundNestedDocs`, `MongoDBv3.4` and `mongoEnableDocLevelTTL`.

~> **Note:** Setting `MongoDBv3.4` also requires setting `EnableMongo`.

~> **Note:** Only `AllowSelfServeUpgradeToMongo36`, `DeleteAllItemsByPartitionKey`, `DisableRateLimitingResponses`, `EnableAggregationPipeline`, `MongoDBv3.4`, `EnableMongoRetryableWrites`, `EnableMongoRoleBasedAccessControl`, `EnableUniqueCompoundNestedDocs`, `EnableMongo16MBDocumentSupport`, `mongoEnableDocLevelTTL`, `EnableTtlOnCustomPath` and `EnablePartialUniqueIndex` can be added to an existing Cosmos DB account.

~> **Note:** Only `DisableRateLimitingResponses` and `EnableMongoRetryableWrites` can be removed from an existing Cosmos DB account.

---

The `virtual_network_rule` block Configures the virtual network subnets allowed to access this Cosmos DB account and supports the following:

* `id` - (Required) The ID of the virtual network subnet.
* `ignore_missing_vnet_service_endpoint` - (Optional) If set to true, the specified subnet will be added as a virtual network rule even if its CosmosDB service endpoint is not active. Defaults to `false`.

---

An `analytical_storage` block supports the following:

* `schema_type` - (Required) The schema type of the Analytical Storage for this Cosmos DB account. Possible values are `FullFidelity` and `WellDefined`.

---

A `capacity` block supports the following:

* `total_throughput_limit` - (Required) The total throughput limit imposed on this Cosmos DB account (RU/s). Possible values are at least `-1`. `-1` means no limit.

---

A `backup` block supports the following:

* `type` - (Required) The type of the `backup`. Possible values are `Continuous` and `Periodic`.

~> **Note:** Migration of `Periodic` to `Continuous` is one-way, changing `Continuous` to `Periodic` forces a new resource to be created.

* `tier` - (Optional) The continuous backup tier. Possible values are `Continuous7Days` and `Continuous30Days`.

* `interval_in_minutes` - (Optional) The interval in minutes between two backups. Possible values are between 60 and 1440. Defaults to `240`.

* `retention_in_hours` - (Optional) The time in hours that each backup is retained. Possible values are between 8 and 720. Defaults to `8`.

* `storage_redundancy` - (Optional) The storage redundancy is used to indicate the type of backup residency. Possible values are `Geo`, `Local` and `Zone`. Defaults to `Geo`.

~> **Note:** You can only configure `interval_in_minutes`, `retention_in_hours` and `storage_redundancy` when the `type` field is set to `Periodic`.

---

A `cors_rule` block supports the following:

* `allowed_headers` - (Required) A list of headers that are allowed to be a part of the cross-origin request.

* `allowed_methods` - (Required) A list of HTTP headers that are allowed to be executed by the origin. Valid options are `DELETE`, `GET`, `HEAD`, `MERGE`, `POST`, `OPTIONS`, `PUT` or `PATCH`.

* `allowed_origins` - (Required) A list of origin domains that will be allowed by CORS.

* `exposed_headers` - (Required) A list of response headers that are exposed to CORS clients.

* `max_age_in_seconds` - (Optional) The number of seconds the client should cache a preflight response. Possible values are between `1` and `2147483647`.

---

An `identity` block supports the following:

* `type` - (Required) The Type of Managed Identity assigned to this Cosmos account. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Cosmos Account.

---

A `restore` block supports the following:

* `source_cosmosdb_account_id` - (Required) The resource ID of the restorable database account from which the restore has to be initiated. The example is `/subscriptions/{subscriptionId}/providers/Microsoft.DocumentDB/locations/{location}/restorableDatabaseAccounts/{restorableDatabaseAccountName}`. Changing this forces a new resource to be created.

~> **Note:** Any database account with `Continuous` type (live account or accounts deleted in last 30 days) is a restorable database account and there cannot be Create/Update/Delete operations on the restorable database accounts. They can only be read and retrieved by `azurerm_cosmosdb_restorable_database_accounts`.

* `restore_timestamp_in_utc` - (Required) The creation time of the database or the collection (Datetime Format `RFC 3339`). Changing this forces a new resource to be created.

* `database` - (Optional) A `database` block as defined below. Changing this forces a new resource to be created.

* `gremlin_database` - (Optional) One or more `gremlin_database` blocks as defined below. Changing this forces a new resource to be created.

* `tables_to_restore` - (Optional) A list of specific tables available for restore. Changing this forces a new resource to be created.

---

A `database` block supports the following:

* `name` - (Required) The database name for the restore request. Changing this forces a new resource to be created.

* `collection_names` - (Optional) A list of the collection names for the restore request. Changing this forces a new resource to be created.

---

A `gremlin_database` block supports the following:

* `name` - (Required) The Gremlin Database name for the restore request. Changing this forces a new resource to be created.

* `graph_names` - (Optional) A list of the Graph names for the restore request. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The CosmosDB Account ID.

* `endpoint` - The endpoint used to connect to the CosmosDB account.

* `read_endpoints` - A list of read endpoints available for this CosmosDB account.

* `write_endpoints` - A list of write endpoints available for this CosmosDB account.

* `primary_key` - The Primary key for the CosmosDB Account.

* `secondary_key` - The Secondary key for the CosmosDB Account.

* `primary_readonly_key` - The Primary read-only Key for the CosmosDB Account.

* `secondary_readonly_key` - The Secondary read-only key for the CosmosDB Account.

* `primary_sql_connection_string` - Primary SQL connection string for the CosmosDB Account.

* `secondary_sql_connection_string` - Secondary SQL connection string for the CosmosDB Account.

* `primary_readonly_sql_connection_string` - Primary readonly SQL connection string for the CosmosDB Account.

* `secondary_readonly_sql_connection_string` - Secondary readonly SQL connection string for the CosmosDB Account.

* `primary_mongodb_connection_string` - Primary Mongodb connection string for the CosmosDB Account.

* `secondary_mongodb_connection_string` - Secondary Mongodb connection string for the CosmosDB Account.

* `primary_readonly_mongodb_connection_string` - Primary readonly Mongodb connection string for the CosmosDB Account.

* `secondary_readonly_mongodb_connection_string` - Secondary readonly Mongodb connection string for the CosmosDB Account.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the CosmosDB Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB Account.
* `update` - (Defaults to 3 hours) Used when updating the CosmosDB Account.
* `delete` - (Defaults to 3 hours) Used when deleting the CosmosDB Account.

## Import

CosmosDB Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_account.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/databaseAccounts/account1
```
