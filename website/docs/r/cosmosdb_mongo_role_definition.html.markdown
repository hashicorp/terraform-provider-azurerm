---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_mongo_role_definition"
description: |-
  Manages a Cosmos DB Mongo Role Definition.
---

# azurerm_cosmosdb_mongo_role_definition

Manages a Cosmos DB Mongo Role Definition.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-ca"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  capabilities {
    name = "EnableMongoRoleBasedAccessControl"
  }

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_mongo_database" "example" {
  name                = "example-mongodb"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
}

resource "azurerm_cosmosdb_mongo_role_definition" "example" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.example.id
  role_name                = "example-roledefinition"
}
```

## Arguments Reference

The following arguments are supported:

* `cosmos_mongo_database_id` - (Required) The resource ID of the Mongo DB. Changing this forces a new resource to be created.

* `role_name` - (Required) The user-friendly name for the Mongo Role Definition. It must be unique for the database account. Changing this forces a new resource to be created.

* `inherited_role_names` - (Optional) A list of Mongo Roles which are inherited to the Mongo Role Definition.

~> **Note:** The role that needs to be inherited should exist in the Mongo DB of `cosmos_mongo_database_id`.

* `privilege` - (Optional) A `privilege` block as defined below.

---

A `privilege` block supports the following:

* `actions` - (Required) A list of actions that are allowed.

* `resource` - (Required) A `resource` block as defined below.

---

A `resource` block supports the following:

* `collection_name` - (Optional) The name of the Mongo DB Collection that the Role Definition is applied.

* `db_name` - (Optional) The name of the Mongo DB that the Role Definition is applied.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cosmos DB Mongo Role Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cosmos DB Mongo Role Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos DB Mongo Role Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Cosmos DB Mongo Role Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cosmos DB Mongo Role Definition.

## Import

Cosmos DB Mongo Role Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_mongo_role_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/databaseAccounts/account1/mongodbRoleDefinitions/dbname1.rolename1
```
