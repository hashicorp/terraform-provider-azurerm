---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_account"
sidebar_current: "docs-azurerm-resource-cosmosdb-account"
description: |-
  Manages a CosmosDB (formally DocumentDB) Account.
---

# azurerm_cosmosdb_account

Manages a CosmosDB (formally DocumentDB) Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
    name     = "${var.resource_group_name}"
    location = "${var.resource_group_location}"
}

resource "random_integer" "ri" {
    min = 10000
    max = 99999
}

resource "azurerm_cosmosdb_account" "db" {
    name                = "tfex-cosmos-db-${random_integer.ri.result}"
    location            = "${azurerm_resource_group.rg.location}"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    offer_type          = "Standard"
    kind                = "GlobalDocumentDB"

    enable_automatic_failover = true

    consistency_policy {
        consistency_level       = "BoundedStaleness"
        max_interval_in_seconds = 10
        max_staleness_prefix    = 200
    }

    geo_location {
        location          = "${var.failover_location}"
        failover_priority = 1
    }

    geo_location {
        id                = "tfex-cosmos-db-${random_integer.ri.result}-customid"
        location          = "${azurerm_resource_group.rg.location}"
        failover_priority = 0
    }
}
```

## Example Usage with virtual_network_rules

```hcl
resource "azurerm_resource_group" "rg" {
    name     = "cosmosDBVNetRules"
    location = "westeurope"
}

resource "azurerm_virtual_network" "test" {
  name                = "virtualNetwork1"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  location            = "${azurerm_resource_group.rg.location}"
  address_space       = ["10.0.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "test1" {
  name                 = "testsubnet1"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
  service_endpoints    = ["Microsoft.AzureCosmosDB"]
}

resource "azurerm_subnet" "test2" {
  name                 = "testsubnet2"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.AzureCosmosDB"]
}

resource "azurerm_cosmosdb_account" "db" {
    depends_on                        = ["azurerm_virtual_network.test"]
    name                              = "cosmos-db-test-vnet-rules"
    location                          = "${azurerm_resource_group.rg.location}"
    resource_group_name               = "${azurerm_resource_group.rg.name}"
    offer_type                        = "Standard"
    kind                              = "GlobalDocumentDB"
    is_virtual_network_filter_enabled = true

    virtual_network_rules {
        id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/cosmosDBVNetRules/providers/Microsoft.Network/virtualNetworks/virtualNetwork1/subnets/testsubnet1"
    }

    virtual_network_rules {
        id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/cosmosDBVNetRules/providers/Microsoft.Network/virtualNetworks/virtualNetwork1/subnets/testsubnet2"
    }

    enable_automatic_failover = true

    consistency_policy {
        consistency_level       = "BoundedStaleness"
        max_interval_in_seconds = 10
        max_staleness_prefix    = 200
    }

    geo_location {
        location          = "${azurerm_resource_group.rg.location}"
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

* `enable_automatic_failover` - (Optional) Enable automatic fail over for this Cosmos DB account.

* `capabilities` - (Optional) Enable capabilities for this Cosmos DB account. Possible values are `EnableTable` and `EnableGremlin`.

`consistency_policy` - Configures the database consistency and supports the following:

* `consistency_level` - (Required) The Consistency Level to use for this CosmosDB Account - can be either `BoundedStaleness`, `Eventual`, `Session`, `Strong` or `ConsistentPrefix`.
* `max_interval_in_seconds` - (Optional) When used with the Bounded Staleness consistency level, this value represents the time amount of staleness (in seconds) tolerated. Accepted range for this value is `5` - `86400` (1 day). Defaults to `5`. Required when `consistency_level` is set to `BoundedStaleness`.
* `max_staleness_prefix` - (Optional) When used with the Bounded Staleness consistency level, this value represents the number of stale requests tolerated. Accepted range for this value is `10` – `2147483647`. Defaults to `100`. Required when `consistency_level` is set to `BoundedStaleness`.

~> **Note**: `max_interval_in_seconds` and `max_staleness_prefix` can only be set to custom values when `consistency_level` is set to `BoundedStaleness` - otherwise they will return the default values shown above.

* `is_virtual_network_filter_enabled` - (Optional) Flag to indicate whether to enable/disable `virtual_network_rules`.

* `virtual_network_rules` - (Optional) Resource `ID` of a subnet for the `azurerm_cosmosdb_account`. 

~> **Note**: Resource `ID` of a subnet must be in this format: /subscriptions/`{subscriptionId}`/resourceGroups/`{groupName}`/providers/Microsoft.Network/virtualNetworks/`{virtualNetworkName}`/subnets/`{subnetName}`. The subnet that is listed as a `virtual_network_rules` must also have `service_endpoints` for `Microsoft.AzureCosmosDB` enabled.

`geo_location` Configures the geographic locations the data is replicated to and supports the following:

* `prefix` - (Optional) The string used to generate the document endpoints for this region. If not specified it defaults to `${cosmosdb_account.name}-${location}`. Changing this causes the location to be deleted and re-provisioned and cannot be changed for the location with failover priority `0`.
* `location` - (Required) The name of the Azure region to host replicated data.
* `failover_priority` - (Required) The failover priority of the region. A failover priority of `0` indicates a write region. The maximum value for a failover priority = (total number of regions - 1). Failover priority values must be unique for each of the regions in which the database account exists. Changing this causes the location to be re-provisioned and cannot be changed for the location with failover priority `0`.

**NOTE:** The `prefix` and `failover_priority` fields of a location cannot be changed for the location with a failover priority of `0`.

## Attributes Reference

The following attributes are exported:

* `id` - The CosmosDB Account ID.

* `endpoint` - The endpoint used to connect to the CosmosDB account.

* `read_endpoints` - A list of read endpoints available for this CosmosDB account.

* `write_endpoints` - A list of write endpoints available for this CosmosDB account.

* `primary_master_key` - The Primary master key for the CosmosDB Account.

* `secondary_master_key` - The Secondary master key for the CosmosDB Account.

* `primary_readonly_master_key` - The Primary read-only master Key for the CosmosDB Account.

* `secondary_readonly_master_key` - The Secondary read-only master key for the CosmosDB Account.

* `connection_strings` - A list of connection strings available for this CosmosDB account. If the kind is `GlobalDocumentDB`, this will be empty.


## Import

CosmosDB Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_account.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/databaseAccounts/account1
```
