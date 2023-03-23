---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_sql_role_assignment"
description: |-
  Manages a Cosmos DB SQL Role Assignment.
---

# azurerm_cosmosdb_sql_role_assignment

Manages a Cosmos DB SQL Role Assignment.

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
  name                = "examplesqlroledef"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_cosmosdb_account.example.name
  type                = "CustomRole"
  assignable_scopes   = [azurerm_cosmosdb_account.example.id]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read"]
  }
}

resource "azurerm_cosmosdb_sql_role_assignment" "example" {
  name                = "736180af-7fbc-4c7f-9004-22735173c1c3"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_cosmosdb_account.example.name
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.example.id
  principal_id        = data.azurerm_client_config.current.object_id
  scope               = azurerm_cosmosdb_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group in which the Cosmos DB SQL Role Assignment is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Account. Changing this forces a new resource to be created.

* `principal_id` - (Required) The ID of the Principal (Client) in Azure Active Directory. Changing this forces a new resource to be created.

* `role_definition_id` - (Required) The resource ID of the Cosmos DB SQL Role Definition.

* `scope` - (Required) The data plane resource path for which access is being granted through this Cosmos DB SQL Role Assignment. Changing this forces a new resource to be created.

* `name` - (Optional) The GUID as the name of the Cosmos DB SQL Role Assignment - one will be generated if not specified. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cosmos DB SQL Role Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cosmos DB SQL Role Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos DB SQL Role Assignment.
* `update` - (Defaults to 30 minutes) Used when updating the Cosmos DB SQL Role Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cosmos DB SQL Role Assignment.

## Import

Cosmos DB SQL Role Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_sql_role_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlRoleAssignments/9e007587-dbcd-4190-84cb-fcab5a09ca39
```
