---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database_sensitivity_label"
description: |-
  Lists MSSQL Database Sensitivity Labels.
---

# List resource: azurerm_mssql_database_sensitivity_label

Lists the current sensitivity labels assigned to columns in an MSSQL Database.

## Example Usage

### List all sensitivity labels in a database

```hcl
list "azurerm_mssql_database_sensitivity_label" "example" {
  provider         = azurerm
  include_resource = true
  config {
    database_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resource-group/providers/Microsoft.Sql/servers/example-server/databases/example-db"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `database_id` - (Required) The ID of the MSSQL Database whose sensitivity labels should be listed.
