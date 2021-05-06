---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_notebook_workspace"
description: |-
  Manages an SQL Notebook Workspace.
---

# azurerm_cosmosdb_notebook_workspace

Manages an SQL Notebook Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmosdb-account"
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

resource "azurerm_cosmosdb_notebook_workspace" "example" {
  name                = "default"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
}
```

## Arguments Reference

The following arguments are supported:
* `name` - (Required) The name which should be used for this SQL Notebook Workspace. Possible value is `default`. Changing this forces a new SQL Notebook Workspace to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SQL Notebook Workspace should exist. Changing this forces a new SQL Notebook Workspace to be created.

* `account_name` - (Required) The name of the Cosmos DB Account to create the SQL Notebook Workspace within. Changing this forces a new SQL Notebook Workspace to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the SQL Notebook Workspace.

* `server_endpoint` - Specifies the endpoint of Notebook server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the =SQL Notebook Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the =SQL Notebook Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the =SQL Notebook Workspace.

## Import

=SQL Notebook Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_notebook_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/databaseAccounts/account1/notebookWorkspaces/notebookWorkspace1
```
