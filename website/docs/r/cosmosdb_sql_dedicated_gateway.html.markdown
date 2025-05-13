---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_sql_dedicated_gateway"
description: |-
  Manages a SQL Dedicated Gateway within a Cosmos DB Account.
---

# azurerm_cosmosdb_sql_dedicated_gateway

Manages a SQL Dedicated Gateway within a Cosmos DB Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-ca"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "BoundedStaleness"
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_dedicated_gateway" "example" {
  cosmosdb_account_id = azurerm_cosmosdb_account.example.id
  instance_count      = 1
  instance_size       = "Cosmos.D4s"
}
```

## Argument Reference

The following arguments are supported:

* `cosmosdb_account_id` - (Required) The resource ID of the CosmosDB Account. Changing this forces a new resource to be created.

* `instance_size` - (Required) The instance size for the CosmosDB SQL Dedicated Gateway. Changing this forces a new resource to be created. Possible values are `Cosmos.D4s`, `Cosmos.D8s` and `Cosmos.D16s`.

* `instance_count` - (Required) The instance count for the CosmosDB SQL Dedicated Gateway. Possible value is between `1` and `5`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CosmosDB SQL Dedicated Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CosmosDB SQL Dedicated Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB SQL Dedicated Gateway.
* `update` - (Defaults to 30 minutes) Used when updating the CosmosDB SQL Dedicated Gateway.
* `delete` - (Defaults to 30 minutes) Used when deleting the CosmosDB SQL Dedicated Gateway.

## Import

CosmosDB SQL Dedicated Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_sql_dedicated_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/account1/services/SqlDedicatedGateway
```
