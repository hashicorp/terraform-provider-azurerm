---
subcategory: "mssql"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database"
description: |-
    Lists mssql database resources.
---

# List resource: azurerm_mssql_database

Lists mssql database resources.

## Example Usage

### List mssql databases in a mssql server

```hcl
list "azurerm_mssql_database" "example" {
  provider = azurerm
  config {
    server_id = "example"
  }
}
```

### List mssql databases in a mssql elastic pool

```hcl
list "azurerm_mssql_database" "example" {
  provider = azurerm
  config {
    elastic_pool_id = "example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `server_id` - (Optional) The id of the mssql server to query.

* `elastic_pool_id` - (Optional) The id of the mssql elastic pool to query.

