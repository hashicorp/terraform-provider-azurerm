---
subcategory: "MsSql"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database"
description: |-
    Lists MSSQL Server resources.
---

# List resource: azurerm_mssql_database

Lists MSSQL Database resources.

## Example Usage

### List all MSSQL Databases in a server

```hcl
list "azurerm_mssql_database" "example" {
  provider = azurerm
  config {
    server_id = "example-server_id"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `server_id` - (Optional) The id of the mssql server to query.
