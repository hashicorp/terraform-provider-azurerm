---
subcategory: "MsSql"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_elasticpool"
description: |-
    Lists MSSQL Elastic Pool resources.
---

# List resource: azurerm_mssql_elasticpool

Lists MSSQL Elastic Pool resources.

## Example Usage

### List all MSSQL Elastic Pools in a server

```hcl
list "azurerm_mssql_elasticpool" "example" {
  provider = azurerm
  config {
    server_id = "example-server_id"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `server_id` - (Required) The id of the mssql server to query.

