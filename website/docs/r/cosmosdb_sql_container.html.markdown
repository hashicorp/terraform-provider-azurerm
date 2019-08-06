---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_sql_container"
sidebar_current: "docs-azurerm-resource-cosmosdb-sql-container"
description: |-
  Manages a SQL Container within a Cosmos DB Account.
---

# azurerm_cosmosdb_sql_container

Manages a SQL Container within a Cosmos DB Account.

## Example Usage

```hcl
data "azurerm_resource_group" "example" {
  name     = "example-rg"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmos-account"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  capabilities                      = []
  enable_automatic_failover         = false
  enable_multiple_write_locations   = false
  is_virtual_network_filter_enabled = false

  consistency_policy {
    consistency_level       = "Session"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }
  geo_location {  
    location          = "${azurerm_resource_group.example.location}"
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "example-db"
  resource_group_name = "${azurerm_cosmosdb_account.example.resource_group_name}"
  account_name        = "${azurerm_cosmosdb_account.example.name}"
}

resource "azurerm_cosmosdb_sql_container" "example" {
  name                = "example-container"
  resource_group_name = "${azurerm_cosmosdb_account.example.resource_group_name}"
  account_name        = "${azurerm_cosmosdb_account.example.name}"
  database_name       = "${azurerm_cosmosdb_sql_database.example.name}"
  partition_key_paths = "/definition/id"
  unique_key_policy {
    path = "/definition/idyard"
  }
  unique_key_policy {
    path = "/definition/idlong"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB SQL Database. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB SQL Database is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Account to create the container within. Changing this forces a new resource to be created.

* `database_name` - (Required) The name of the Cosmos DB SQL Database to create the container within. Changing this forces a new resource to be created.

* `partition_key_paths` - (Optional) Define a partition key. Changing this forces a new resource to be created.

* `unique_key_policy` - (Optional) Define Unique Keys for data integrity. A unique_key_policy block as defined below.

---

A `unique_key_policy` block supports the following:

* `path` - (Required) The path of the unique key in the path format /example/id.


## Attributes Reference

The following attributes are exported:

* `id` - the Cosmos DB SQL Database ID.

## Import

Cosmos SQL Database can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_sql_container.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/apis/sql/databases/database1/containers/example
```

