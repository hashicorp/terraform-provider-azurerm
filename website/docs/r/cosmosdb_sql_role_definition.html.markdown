---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_sql_role_definition"
description: |-
  Manages a Cosmos DB SQL Role Definition.
---

# azurerm_cosmosdb_sql_role_definition

Manages a Cosmos DB SQL Role Definition.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmosdb"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_role_definition" "example" {
  role_definition_id  = "84cf3a8b-4122-4448-bce2-fa423cfe0a15"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_cosmosdb_account.example.name
  name                = "acctestsqlrole"
  assignable_scopes   = ["${azurerm_cosmosdb_account.example.id}/dbs/sales"]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group in which the Cosmos DB SQL Role Definition is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Account. Changing this forces a new resource to be created.

* `assignable_scopes` - (Required) A list of fully qualified scopes at or below which Role Assignments may be created using this Cosmos DB SQL Role Definition. It will allow application of this Cosmos DB SQL Role Definition on the entire Database Account or any underlying Database/Collection. Scopes higher than Database Account are not enforceable as assignable scopes.

~> **Note:** The resources referenced in assignable scopes need not exist.

* `name` - (Required) An user-friendly name for the Cosmos DB SQL Role Definition which must be unique for the Database Account.

* `permissions` - (Required) A `permissions` block as defined below.

* `role_definition_id` - (Optional) The GUID as the name of the Cosmos DB SQL Role Definition - one will be generated if not specified. Changing this forces a new resource to be created.

* `type` - (Optional) The type of the Cosmos DB SQL Role Definition. Possible values are `BuiltInRole` and `CustomRole`. Defaults to `CustomRole`. Changing this forces a new resource to be created.

---

A `permissions` block supports the following:

* `data_actions` - (Required) A list of data actions that are allowed for the Cosmos DB SQL Role Definition.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cosmos DB SQL Role Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cosmos DB SQL Role Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos DB SQL Role Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Cosmos DB SQL Role Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cosmos DB SQL Role Definition.

## Import

Cosmos DB SQL Role Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_sql_role_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlRoleDefinitions/28b3c337-f436-482b-a167-c2618dc52033
```
