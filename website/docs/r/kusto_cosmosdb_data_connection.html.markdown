---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_cosmosdb_data_connection"
description: |-
  Manages a Kusto / Cosmos Database Data Connection.
---

# azurerm_kusto_cosmosdb_data_connection

Manages a Kusto / Cosmos Database Data Connection.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "accexampleRG"
  location = "West Europe"
}

data "azurerm_role_definition" "builtin" {
  role_definition_id = "fbdf93bf-df7d-467e-a4d2-9458aa1360c8"
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = data.azurerm_role_definition.builtin.name
  principal_id         = azurerm_kusto_cluster.example.identity[0].principal_id
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "accexample-ca"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "Session"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "accexamplecosmosdbsqldb"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
}

resource "azurerm_cosmosdb_sql_container" "example" {
  name                = "accexamplecosmosdbsqlcon"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_sql_database.example.name
  partition_key_path  = "/part"
  throughput          = 400
}

data "azurerm_cosmosdb_sql_role_definition" "example" {
  role_definition_id  = "00000000-0000-0000-0000-000000000001"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_cosmosdb_account.example.name
}


resource "azurerm_cosmosdb_sql_role_assignment" "example" {
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_cosmosdb_account.example.name
  role_definition_id  = data.azurerm_cosmosdb_sql_role_definition.example.id
  principal_id        = azurerm_kusto_cluster.example.identity[0].principal_id
  scope               = azurerm_cosmosdb_account.example.id
}

resource "azurerm_kusto_cluster" "example" {
  name                = "accexamplekc%s"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kusto_database" "example" {
  name                = "accexamplekd"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.example.name
}

resource "azurerm_kusto_script" "example" {
  name           = "create-table-script"
  database_id    = azurerm_kusto_database.example.id
  script_content = <<SCRIPT
.create table TestTable(Id:string, Name:string, _ts:long, _timestamp:datetime)
.create table TestTable ingestion json mapping "TestMapping"
'['
'    {"column":"Id","path":"$.id"},'
'    {"column":"Name","path":"$.name"},'
'    {"column":"_ts","path":"$._ts"},'
'    {"column":"_timestamp","path":"$._ts", "transform":"DateTimeFromUnixSeconds"}'
']'
.alter table TestTable policy ingestionbatching "{'MaximumBatchingTimeSpan': '0:0:10', 'MaximumNumberOfItems': 10000}"
SCRIPT
}

resource "azurerm_kusto_cosmosdb_data_connection" "test" {
  name                 = "acctestkcd%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  cosmosdb_account_id  = azurerm_cosmosdb_account.test.id
  cosmosdb_database    = azurerm_cosmosdb_sql_database.test.name
  cosmosdb_container   = azurerm_cosmosdb_sql_container.test.name
  cluster_name         = azurerm_kusto_cluster.test.name
  database_name        = azurerm_kusto_database.test.name
  managed_identity_id  = azurerm_kusto_cluster.test.id
  table_name           = "TestTable"
  mapping_rule_name    = "TestMapping"
  retrieval_start_date = "2023-06-26T12:00:00.6554616Z"
}
```

## Arguments Reference

The following arguments are supported:

* `cluster_name` - (Required) The name of the Kusto cluster.Changing this forces a new Data Explorer to be created.

* `cosmosdb_account_id` - (Required) The resource ID of the Cosmos DB account used to create the data connection. Changing this forces a new Data Explorer to be created.

* `cosmosdb_container` - (Required) The name of an existing container in the Cosmos DB database. Changing this forces a new Data Explorer to be created.

* `cosmosdb_database` - (Required) The name of an existing database in the Cosmos DB account. Changing this forces a new Data Explorer to be created.

* `database_name` - (Required) The name of the database in the Kusto cluster. Changing this forces a new Data Explorer to be created.

* `location` - (Required) The Azure Region where the Data Explorer should exist. Changing this forces a new Data Explorer to be created.

* `managed_identity_id` - (Required) The resource ID of a managed system or user-assigned identity. The identity is used to authenticate with Cosmos DB. Changing this forces a new Data Explorer to be created.

* `name` - (Required) The name of the data connection. Changing this forces a new Data Explorer to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Data Explorer should exist. Changing this forces a new Data Explorer to be created.

* `table_name` - (Required) The case-sensitive name of the existing target table in your cluster. Retrieved data is ingested into this table. Changing this forces a new Data Explorer to be created.

---

* `mapping_rule_name` - (Optional) The name of an existing mapping rule to use when ingesting the retrieved data.

* `retrieval_start_date` - (Optional) If defined, the data connection retrieves Cosmos DB documents created or updated after the specified retrieval start date.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Explorer.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Explorer.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Explorer.
* `update` - (Defaults to 30 minutes) Used when updating the Data Explorer.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Explorer.

## Import

Data Explorers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_cosmosdb_data_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/dataConnections/dataConnection1
```
