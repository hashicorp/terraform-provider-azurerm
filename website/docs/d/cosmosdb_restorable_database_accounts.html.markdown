---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_restorable_database_accounts"
description: |-
  Gets information about Cosmos DB Restorable Database Accounts.
---

# Data Source: azurerm_cosmosdb_restorable_database_accounts

Use this data source to access information about Cosmos DB Restorable Database Accounts.

## Example Usage

```hcl
data "azurerm_cosmosdb_restorable_database_accounts" "example" {
  name     = "example-ca"
  location = "West Europe"
}

output "id" {
  value = data.azurerm_cosmosdb_restorable_database_accounts.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Cosmos DB Database Account.

* `location` - (Required) The location where the Cosmos DB Database Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cosmos DB Restorable Database Accounts.

* `accounts` - One or more `accounts` blocks as defined below.

---

An `accounts` block exports the following:

* `id` - The ID of the Cosmos DB Restorable Database Account.

* `api_type` - The API type of the Cosmos DB Restorable Database Account.

* `creation_time` - The creation time of the Cosmos DB Restorable Database Account.

* `deletion_time` - The deletion time of the Cosmos DB Restorable Database Account.

* `restorable_locations` - One or more `restorable_locations` blocks as defined below.

---

An `restorable_locations` block exports the following:

* `creation_time` - The creation time of the regional Cosmos DB Restorable Database Account.

* `deletion_time` - The deletion time of the regional Cosmos DB Restorable Database Account.

* `location` - The location of the regional Cosmos DB Restorable Database Account.

* `regional_database_account_instance_id` - The instance ID of the regional Cosmos DB Restorable Database Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos DB Restorable Database Accounts.
