---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_account"
sidebar_current: "docs-azurerm-datasource-azurerm-cosmosdb-account"
description: |-
  Get information about the specified CosmosDB (formally DocumentDB) Account.
---

# Data Source: azurerm_cosmosdb_account

Use this data source to access the properties of an Azure CosmosDB (formally DocumentDB) Account.

## Example Usage

```hcl
data "azurerm_cosmosdb_account" "test" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = "tfex-cosmosdb-account-rg"
}

output "cosmosdb_account_endpoint" {
  value = "${data.azurerm_cosmosdb_account.jobs.endpoint}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the CosmosDB Account. 

* `resource_group_name` - (Required) Specifies the name of the resource group in which the CosmosDB Account resides. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CosmosDB Account.

* `location` - The Azure location where the resource exists. 

* `tags` - A mapping of tags assigned to the resource.

* `offer_type` - The Offer Type to used by this CosmosDB Account.

* `kind` - The Kind of the CosmosDB account.

* `ip_range_filter` - The current IP Filter for this CosmosDB account

* `enable_automatic_failover` - If automatic failover is enabled for this CosmosDB Account.

`consistency_policy` The current consistency Settings for this CosmosDB account with the following properties:

* `consistency_level` - The Consistency Level used by this CosmosDB Account. 
* `max_interval_in_seconds` - The amount of staleness (in seconds) tolerated when the consistency level is Bounded Staleness. 
* `max_staleness_prefix` - The number of stale requests tolerated when the consistency level is Bounded Staleness.  


`geo_location` The geographic locations data is replicated to with the following properties:

* `id` - The ID of the location.
* `location` - The name of the Azure region hosting replicated data.
* `priority` - The locations fail over priority.

* `endpoint` - The endpoint used to connect to the CosmosDB account.

* `read_endpoints` - A list of read endpoints available for this CosmosDB account.

* `write_endpoints` - A list of write endpoints available for this CosmosDB account.

* `primary_master_key` - The Primary master key for the CosmosDB Account.

* `secondary_master_key` - The Secondary master key for the CosmosDB Account.

* `primary_readonly_master_key` - The Primary read-only master Key for the CosmosDB Account.

* `secondary_readonly_master_key` - The Secondary read-only master key for the CosmosDB Account.
