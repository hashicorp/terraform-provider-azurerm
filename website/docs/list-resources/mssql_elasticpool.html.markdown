---
subcategory: "mssql"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_elasticpool"
description: |-
    Lists mssql elasticpool resources.
---

# List resource: azurerm_mssql_elasticpool

Lists mssql elasticpool resources.

## Example Usage

### List mssql elasticpools in a mssql server

```hcl
list "azurerm_mssql_elasticpool" "example" {
  provider = azurerm
  config {
    server_id = "example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `server_id` - (Required) The id of the mssql server to query.

