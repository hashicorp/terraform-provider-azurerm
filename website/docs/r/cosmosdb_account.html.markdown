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

  enable_automatic_failover = true

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
    location          = "eastus"
    failover_priority = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the CosmosDB Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the CosmosDB Account is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `offer_type` - (Required) Specifies the Offer Type to use for this CosmosDB Account - currently this can only be set to `Standard`.

* `analytical_storage` - (Optional) An `analytical_storage` block as defined below.

* `capacity` - (Optional) A `capacity` block as defined below.

* `create_mode` - (Optional) The creation mode for the CosmosDB Account. Possible values are `Default` and `Restore`. Changing this forces a new resource to be created.

~> **NOTE:** `create_mode` only works when `backup.type` is `Continuous`.

* `default_identity_type` - (Optional) The default identity for accessing Key Vault. Possible values are `FirstPartyIdentity`, `SystemAssignedIdentity` or start with `UserAssignedIdentity`. Defaults to `FirstPartyIdentity`.

* `kind` - (Optional) Specifies the Kind of CosmosDB to create - possible values are `GlobalDocumentDB` and `MongoDB`. Defaults to `GlobalDocumentDB`. Changing this forces a new resource to be created.

* `consistency_policy` - (Required) Specifies a `consistency_policy` resource, used to define the consistency policy for this CosmosDB account.

* `geo_location` - (Required) Specifies a `geo_location` resource, used to define where data should be replicated with the `failover_priority` 0 specifying the primary location. Value is a `geo_location` block as defined below.

* `ip_range_filter` - (Optional) CosmosDB Firewall Support: This value specifies the set of IP addresses or IP address ranges in CIDR form to be included as the allowed list of client IP's for a given database account. IP addresses/ranges must be comma separated and must not contain any spaces.

~> **NOTE:** To enable the "Allow access from the Azure portal" behavior, you should add the IP addresses provided by the [documentation](https://docs.microsoft.com/azure/cosmos-db/how-to-configure-firewall#allow-requests-from-the-azure-portal) to this list.

~> **NOTE:** To enable the "Accept connections from within public Azure datacenters" behavior, you should add `0.0.0.0` to the list, see the [documentation](https://docs.microsoft.com/azure/cosmos-db/how-to-configure-firewall#allow-requests-from-global-azure-datacenters-or-other-sources-within-azure) for more details.

* `enable_free_tier` - (Optional) Enable Free Tier pricing option for this Cosmos DB account. Defaults to `false`. Changing this forces a new resource to be created.

* `analytical_storage_enabled` - (Optional) Enable Analytical Storage option for this Cosmos DB account. Defaults to `false`. Changing this forces a new resource to be created.

* `enable_automatic_failover` - (Optional) Enable automatic fail over for this Cosmos DB account.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this CosmosDB account.

* `capabilities` - (Optional) The capabilities which should be enabled for this Cosmos DB account. Value is a `capabilities` block as defined below. Changing this forces a new resource to be created.

* `is_virtual_network_filter_enabled` - (Optional) Enables virtual network filtering for this Cosmos DB account.

* `key_vault_key_id` - (Optional) A versionless Key Vault Key ID for CMK encryption. Changing this forces a new resource to be created.

~> **NOTE:** When referencing an `azurerm_key_vault_key` resource, use `versionless_id` instead of `id`

~> **NOTE:** In order to use a `Custom Key` from Key Vault for encryption you must grant Azure Cosmos DB Service access to your key vault. For instructions on how to configure your Key Vault correctly please refer to the [product documentation](https://docs.microsoft.com/azure/cosmos-db/how-to-setup-cmk#add-an-access-policy-to-your-azure-key-vault-instance)

* `virtual_network_rule` - (Optional) Specifies a `virtual_network_rules` resource, used to define which subnets are allowed to access this CosmosDB account.

* `enable_multiple_write_locations` - (Optional) Enable multiple write locations for this Cosmos DB account.

* `access_key_metadata_writes_enabled` - (Optional) Is write operations on metadata resources (databases, containers, throughput) via account keys enabled? Defaults to `true`.

* `mongo_server_version` - (Optional) The Server Version of a MongoDB account. Possible values are `4.2`, `4.0`, `3.6`, and `3.2`.

* `network_acl_bypass_for_azure_services` - (Optional) If Azure services can bypass ACLs. Defaults to `false`.

* `network_acl_bypass_ids` - (Optional) The list of resource Ids for Network Acl Bypass for this Cosmos DB account.

* `local_authentication_disabled` - (Optional) Disable local authentication and ensure only MSI and AAD can be used exclusively for authentication. Defaults to `false`. Can be set only when using the SQL API.

* `backup` - (Optional) A `backup` block as defined below.

* `cors_rule` - (Optional) A `cors_rule` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `restore` - (Optional) A `restore` block as defined below.

~> **NOTE:** `restore` should be set when `create_mode` is `Restore`.

---

`consistency_policy` Configures the database consistency and supports the following:

* `consistency_level` - (Required) The Consistency Level to use for this CosmosDB Account - can be either `BoundedStaleness`, `Eventual`, `Session`, `Strong` or `ConsistentPrefix`.
* `max_interval_in_seconds` - (Optional) When used with the Bounded Staleness consistency level, this value represents the time amount of staleness (in seconds) tolerated. Accepted range for this value is `5` - `86400` (1 day). Defaults to `5`. Required when `consistency_level` is set to `BoundedStaleness`.
* `max_staleness_prefix` - (Optional) When used with the Bounded Staleness consistency level, this value represents the number of stale requests tolerated. Accepted range for this value is `10` – `2147483647`. Defaults to `100`. Required when `consistency_level` is set to `BoundedStaleness`.

~> **Note:** `max_interval_in_seconds` and `max_staleness_prefix` can only be set to custom values when `consistency_level` is set to `BoundedStaleness` - otherwise they will return the default values shown above.

---

`geo_location` Configures the geographic locations the data is replicated to and supports the following:

* `location` - (Required) The name of the Azure region to host replicated data.
* `failover_priority` - (Required) The failover priority of the region. A failover priority of `0` indicates a write region. The maximum value for a failover priority = (total number of regions - 1). Failover priority values must be unique for each of the regions in which the database account exists. Changing this causes the location to be re-provisioned and cannot be changed for the location with failover priority `0`.
* `zone_redundant` - (Optional) Should zone redundancy be enabled for this region? Defaults to `false`.

---

`capabilities` Configures the capabilities to enable for this Cosmos DB account:

* `name` - (Required) The capability to enable - Possible values are `AllowSelfServeUpgradeToMongo36`, `DisableRateLimitingResponses`, `EnableAggregationPipeline`, `EnableCassandra`, `EnableGremlin`, `EnableMongo`, `EnableTable`, `EnableServerless`, `MongoDBv3.4` and `mongoEnableDocLevelTTL`. 

**NOTE:**  Setting `MongoDBv3.4` also requires setting `EnableMongo`.

---

`virtual_network_rule` Configures the virtual network subnets allowed to access this Cosmos DB account and supports the following:

* `id` - (Required) The ID of the virtual network subnet.
* `ignore_missing_vnet_service_endpoint` - (Optional) If set to true, the specified subnet will be added as a virtual network rule even if its CosmosDB service endpoint is not active. Defaults to `false`.

---

A `analytical_storage` block supports the following:

* `schema_type` - (Required) The schema type of the Analytical Storage for this Cosmos DB account. Possible values are `FullFidelity` and `WellDefined`.

---

A `capacity` block supports the following:

* `total_throughput_limit` - (Required) The total throughput limit imposed on this Cosmos DB account (RU/s). Possible values are at least `-1`. `-1` means no limit.

---

A `backup` block supports the following:

* `type` - (Required) The type of the `backup`. Possible values are `Continuous` and `Periodic`. Defaults to `Periodic`. Migration of `Periodic` to `Continuous` is one-way, changing `Continuous` to `Periodic` forces a new resource to be created.

* `interval_in_minutes` - (Optional) The interval in minutes between two backups. This is configurable only when `type` is `Periodic`. Possible values are between 60 and 1440.

* `retention_in_hours` - (Optional) The time in hours that each backup is retained. This is configurable only when `type` is `Periodic`. Possible values are between 8 and 720.

* `storage_redundancy` - (Optional) The storage redundancy which is used to indicate type of backup residency. This is configurable only when `type` is `Periodic`. Possible values are `Geo`, `Local` and `Zone`.

---

A `cors_rule` block supports the following:

* `allowed_headers` - (Required) A list of headers that are allowed to be a part of the cross-origin request.

* `allowed_methods` - (Required) A list of HTTP headers that are allowed to be executed by the origin. Valid options are  `DELETE`, `GET`, `HEAD`, `MERGE`, `POST`, `OPTIONS`, `PUT` or `PATCH`.

* `allowed_origins` - (Required) A list of origin domains that will be allowed by CORS.

* `exposed_headers` - (Required) A list of response headers that are exposed to CORS clients.

* `max_age_in_seconds` - (Required) The number of seconds the client should cache a preflight response.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Cosmos Account. The only possible value is `SystemAssigned`.

---

A `restore` block supports the following:

* `source_cosmosdb_account_id` - (Required) The resource ID of the restorable database account from which the restore has to be initiated. The example is `/subscriptions/{subscriptionId}/providers/Microsoft.DocumentDB/locations/{location}/restorableDatabaseAccounts/{restorableDatabaseAccountName}`. Changing this forces a new resource to be created.

**NOTE:** Any database account with `Continuous` type (live account or accounts deleted in last 30 days) are the restorable database accounts and there cannot be Create/Update/Delete operations on the restorable database accounts. They can only be read and be retrieved by `azurerm_cosmosdb_restorable_database_accounts`.

* `restore_timestamp_in_utc` - (Required) The creation time of the database or the collection (Datetime Format `RFC 3339`). Changing this forces a new resource to be created.

* `database` - (Optional) A `database` block as defined below. Changing this forces a new resource to be created.

---

A `database` block supports the following:

* `name` - (Required) The database name for the restore request. Changing this forces a new resource to be created.

* `collection_names` - (Optional) A list of the collection names for the restore request. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The CosmosDB Account ID.

* `endpoint` - The endpoint used to connect to the CosmosDB account.

* `read_endpoints` - A list of read endpoints available for this CosmosDB account.

* `write_endpoints` - A list of write endpoints available for this CosmosDB account.

* `primary_key` - The Primary key for the CosmosDB Account.

* `secondary_key` - The Secondary key for the CosmosDB Account.

* `primary_readonly_key` - The Primary read-only Key for the CosmosDB Account.

* `secondary_readonly_key` - The Secondary read-only key for the CosmosDB Account.

* `connection_strings` - A list of connection strings available for this CosmosDB account.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 180 minutes) Used when creating the CosmosDB Account.
* `update` - (Defaults to 180 minutes) Used when updating the CosmosDB Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB Account.
* `delete` - (Defaults to 180 minutes) Used when deleting the CosmosDB Account.

## Import

CosmosDB Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_account.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/databaseAccounts/account1
```
