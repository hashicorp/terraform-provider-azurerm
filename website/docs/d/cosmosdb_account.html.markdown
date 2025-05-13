---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_account"
description: |-
  Gets information about an existing CosmosDB (formally DocumentDB) Account.
---

# Data Source: azurerm_cosmosdb_account

Use this data source to access information about an existing CosmosDB (formally DocumentDB) Account.

## Example Usage

```hcl
data "azurerm_cosmosdb_account" "example" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = "tfex-cosmosdb-account-rg"
}

output "cosmosdb_account_endpoint" {
  value = data.azurerm_cosmosdb_account.example.endpoint
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the CosmosDB Account.

* `resource_group_name` - Specifies the name of the resource group in which the CosmosDB Account resides.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CosmosDB Account.

* `location` - The Azure location where the resource exists.

* `tags` - A mapping of tags assigned to the resource.

* `offer_type` - The Offer Type to used by this CosmosDB Account.

* `kind` - The Kind of the CosmosDB account.

* `key_vault_key_id` - The Key Vault key URI for CMK encryption.

~> **Note:** The CosmosDB service always uses the latest version of the specified key.

* `ip_range_filter` - The current IP Filter for this CosmosDB account

* `free_tier_enabled` - If Free Tier pricing option is enabled for this CosmosDB Account. You can have up to one free tier Azure Cosmos DB account per Azure subscription.

* `automatic_failover_enabled` - If automatic failover is enabled for this CosmosDB Account.

* `capabilities` - Capabilities enabled on this Cosmos DB account.

* `is_virtual_network_filter_enabled` - If virtual network filtering is enabled for this Cosmos DB account.

* `virtual_network_rule` - Subnets that are allowed to access this CosmosDB account.

* `multiple_write_locations_enabled` - If multiple write locations are enabled for this Cosmos DB account.

`consistency_policy` The current consistency Settings for this CosmosDB account with the following properties:

* `consistency_level` - The Consistency Level used by this CosmosDB Account.
* `max_interval_in_seconds` - The amount of staleness (in seconds) tolerated when the consistency level is Bounded Staleness.
* `max_staleness_prefix` - The number of stale requests tolerated when the consistency level is Bounded Staleness.

`geo_location` The geographic locations data is replicated to with the following properties:

* `id` - The ID of the location.
* `location` - The name of the Azure region hosting replicated data.
* `priority` - The locations fail over priority.

`virtual_network_rule` The virtual network subnets allowed to access this Cosmos DB account with the following properties:

* `id` - The ID of the virtual network subnet.

* `endpoint` - The endpoint used to connect to the CosmosDB account.

* `read_endpoints` - A list of read endpoints available for this CosmosDB account.

* `write_endpoints` - A list of write endpoints available for this CosmosDB account.

* `primary_key` - The primary key for the CosmosDB account.

* `secondary_key` - The secondary key for the CosmosDB account.

* `primary_readonly_key` - The primary read-only Key for the CosmosDB account.

* `secondary_readonly_key` - The secondary read-only key for the CosmosDB account.

* `primary_sql_connection_string` - The primary SQL connection string for the CosmosDB Account.

* `secondary_sql_connection_string` - The secondary SQL connection string for the CosmosDB Account.

* `primary_readonly_sql_connection_string` - The primary read-only SQL connection string for the CosmosDB account.

* `secondary_readonly_sql_connection_string` - The secondary read-only SQL connection string for the CosmosDB account.

* `primary_mongodb_connection_string` - The primary Mongodb connection string for the CosmosDB account.

* `secondary_mongodb_connection_string` - The secondary Mongodb connection string for the CosmosDB account.

* `primary_readonly_mongodb_connection_string` - The primary readonly Mongodb connection string for the CosmosDB account.

* `secondary_readonly_mongodb_connection_string` - The secondary readonly Mongodb connection string for the CosmosDB account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB Account.
