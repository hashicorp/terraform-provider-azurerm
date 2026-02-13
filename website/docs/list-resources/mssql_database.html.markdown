---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database"
description: |-
    Lists Mssql Database resources.
---

# List resource: azurerm_mssql_database

Lists Mssql Database resources.

## Example Usage

### List Mssql Databases in a Mssql Server

```hcl
list "azurerm_mssql_database" "example" {
  provider = azurerm
  config {
    server_id = /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver
  }
}
```

### List Mssql Databases in a Mssql Elastic Pool

```hcl
list "azurerm_mssql_database" "example" {
  provider = azurerm
  config {
    elastic_pool_id = /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/elasticPools/myelasticpoolname
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `server_id` - (Optional) The ID of the Mssql Server to query.

* `elastic_pool_id` - (Optional) The ID of the Mssql Elastic Pool to query.
