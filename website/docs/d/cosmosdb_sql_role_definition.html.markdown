---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_sql_role_definition"
description: |-
  Gets information about an existing Cosmos DB SQL Role Definition.
---

# azurerm_cosmosdb_sql_role_definition

Use this data source to access information about an existing Cosmos DB SQL Role Definition.

## Example Usage

```hcl
data "azurerm_cosmosdb_sql_role_definition" "example" {
  name                = "tfex-cosmosdb-sql-role-definition"
  resource_group_name = "tfex-cosmosdb-sql-role-definition-rg"
  account_name        = "tfex-cosmosdb-sql-role-definition-account-name"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) An user-friendly name for the Cosmos DB SQL Role Definition.

* `resource_group_name` - (Required) The name of the Resource Group in which the Cosmos DB SQL Role Definition is created.

* `account_name` - (Required) The name of the Cosmos DB Account.

* `role_definition_id` - (Optional) The GUID as the name of the Cosmos DB SQL Role Definition.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cosmos DB SQL Role Definition.

* `assignable_scopes` - A list of fully qualified scopes at or below which Role Assignments may be created using this Cosmos DB SQL Role Definition.

* `type` - The type of the Cosmos DB SQL Role Definition.

* `permissions` - A `permissions` block as defined below.

---

A `permissions` block supports the following:

* `data_actions` - A list of data actions that are allowed for the Cosmos DB SQL Role Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos DB SQL Role Definition.
