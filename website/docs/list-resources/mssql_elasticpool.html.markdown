---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_elasticpool"
description: |-
    Lists Mssql Elasticpool resources.
---

# List resource: azurerm_mssql_elasticpool

Lists Mssql Elasticpool resources.

## Example Usage

### List Mssql Elasticpools in a Mssql Server

```hcl
list "azurerm_mssql_elasticpool" "example" {
  provider = azurerm
  config {
    mssql_server_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `mssql_server_id` - (Required) The ID of the Mssql Server to query.
