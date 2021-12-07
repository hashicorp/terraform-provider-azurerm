---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_restorable_database_account"
description: |-
  Gets information about Cosmos DB Restorable Database Account.
---

# Data Source: azurerm_cosmosdb_restorable_database_account

Use this data source to access information about Cosmos DB Restorable Database Account.

## Example Usage

```hcl
data "azurerm_cosmosdb_restorable_database_account" "example" {
  name     = "example-ca"
  location = "West Europe"
}

output "id" {
  value = data.azurerm_cosmosdb_restorable_database_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Cosmos DB Database Account.

* `location` - (Required) The location where the Cosmos DB Database Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cosmos DB Restorable Database Account.

* `restorable_db_account_ids` - A list of the Cosmos DB Restorable Database Account IDs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos DB Restorable Database Account.
