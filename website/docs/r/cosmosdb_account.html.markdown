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
resource "azurerm_resource_group" "rg" {
  name     = var.resource_group_name
  location = var.resource_group_location
}

resource "random_integer" "ri" {
  min = 10000
  max = 99999
}

resource "azurerm_cosmosdb_account" "db" {
  name                = "tfex-cosmos-db-${random_integer.ri.result}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

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

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = var.failover_location
    failover_priority = 1
  }

  geo_location {
    location          = azurerm_resource_group.rg.location
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

* `kind` - (Optional) Specifies the Kind of CosmosDB to create - possible values are `GlobalDocumentDB` and `MongoDB`. Defaults to `GlobalDocumentDB`. Changing this forces a new resource to be created.

* `consistency_policy` - (Required) Specifies a `consistency_policy` resource, used to define the consistency policy for this CosmosDB account.

* `geo_location` - (Required) Specifies a `geo_location` resource, used to define where data should be replicated with the `failover_priority` 0 specifying the primary location.

* `ip_range_filter` - (Optional) CosmosDB Firewall Support: This value specifies the set of IP addresses or IP address ranges in CIDR form to be included as the allowed list of client IP's for a given database account. IP addresses/ranges must be comma separated and must not contain any spaces.

* `enable_free_tier` - (Optional) Enable Free Tier pricing option for this Cosmos DB account. Defaults to `false`. Changing this forces a new resource to be created.

* `enable_automatic_failover` - (Optional) Enable automatic fail over for this Cosmos DB account.

* `capabilities` - (Optional) The capabilities which should be enabled for this Cosmos DB account. Possible values are `EnableAggregationPipeline`, `EnableCassandra`, `EnableGremlin`, `EnableTable`, `MongoDBv3.4`, `EnableServerless`, and `mongoEnableDocLevelTTL`.

* `is_virtual_network_filter_enabled` - (Optional) Enables virtual network filtering for this Cosmos DB account.

* `virtual_network_rule` - (Optional) Specifies a `virtual_network_rules` resource, used to define which subnets are allowed to access this CosmosDB account.

* `enable_multiple_write_locations` - (Optional) Enable multi-master support for this Cosmos DB account.

`consistency_policy` Configures the database consistency and supports the following:

* `consistency_level` - (Required) The Consistency Level to use for this CosmosDB Account - can be either `BoundedStaleness`, `Eventual`, `Session`, `Strong` or `ConsistentPrefix`.
* `max_interval_in_seconds` - (Optional) When used with the Bounded Staleness consistency level, this value represents the time amount of staleness (in seconds) tolerated. Accepted range for this value is `5` - `86400` (1 day). Defaults to `5`. Required when `consistency_level` is set to `BoundedStaleness`.
* `max_staleness_prefix` - (Optional) When used with the Bounded Staleness consistency level, this value represents the number of stale requests tolerated. Accepted range for this value is `10` – `2147483647`. Defaults to `100`. Required when `consistency_level` is set to `BoundedStaleness`.

~> **Note**: `max_interval_in_seconds` and `max_staleness_prefix` can only be set to custom values when `consistency_level` is set to `BoundedStaleness` - otherwise they will return the default values shown above.

`geo_location` Configures the geographic locations the data is replicated to and supports the following:

* `prefix` - (Optional) The string used to generate the document endpoints for this region. If not specified it defaults to `${cosmosdb_account.name}-${location}`. Changing this causes the location to be deleted and re-provisioned and cannot be changed for the location with failover priority `0`.
* `location` - (Required) The name of the Azure region to host replicated data.
* `failover_priority` - (Required) The failover priority of the region. A failover priority of `0` indicates a write region. The maximum value for a failover priority = (total number of regions - 1). Failover priority values must be unique for each of the regions in which the database account exists. Changing this causes the location to be re-provisioned and cannot be changed for the location with failover priority `0`.
* `zone_redundant` - (Optional) Should zone redundancy be enabled for this region? Defaults to `false`.

`capabilities` Configures the capabilities to enable for this Cosmos DB account:

* `name` - (Required) The capability to enable - Possible values are `AllowSelfServeUpgradeToMongo36`, `DisableRateLimitingResponses`, `EnableAggregationPipeline`, `EnableCassandra`, `EnableGremlin`,`EnableMongo`, `EnableTable`, `EnableServerless`, `MongoDBv3.4` and `mongoEnableDocLevelTTL`.

**NOTE:** The `prefix` and `failover_priority` fields of a location cannot be changed for the location with a failover priority of `0`.

`virtual_network_rule` Configures the virtual network subnets allowed to access this Cosmos DB account and supports the following:

* `id` - (Required) The ID of the virtual network subnet.
* `ignore_missing_vnet_service_endpoint` - (Optional) If set to true, the specified subnet will be added as a virtual network rule even if its CosmosDB service endpoint is not active. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The CosmosDB Account ID.

* `endpoint` - The endpoint used to connect to the CosmosDB account.

* `read_endpoints` - A list of read endpoints available for this CosmosDB account.

* `write_endpoints` - A list of write endpoints available for this CosmosDB account.

* `primary_key` - The Primary master key for the CosmosDB Account.

* `secondary_key` - The Secondary master key for the CosmosDB Account.

* `primary_readonly_key` - The Primary read-only master Key for the CosmosDB Account.

* `secondary_readonly_key` - The Secondary read-only master key for the CosmosDB Account.

* `connection_strings` - A list of connection strings available for this CosmosDB account.

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
