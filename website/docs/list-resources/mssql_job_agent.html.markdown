---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_job_agent"
description: |-
    Lists mssql job agent resources.
---

# List resource: azurerm_mssql_job_agent

Lists mssql job agent resources.

## Example Usage

### List mssql job agents in a mssql server

```hcl
list "azurerm_mssql_job_agent" "example" {
  provider = azurerm
  config {
    server_id = "example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `server_id` - (Required) The id of the mssql server to query.

