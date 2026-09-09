---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_flexible_server_configuration"
description: |-
    Lists Postgresql Flexible Server Configuration resources.
---

# List resource: azurerm_postgresql_flexible_server_configuration

Lists Postgresql Flexible Server Configuration resources.

## Example Usage

### List PostgreSQL Flexible Server Configurations in a Flexible Server

```hcl
list "azurerm_postgresql_flexible_server_configuration" "example" {
  provider = azurerm
  config {
    server_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.DBforPostgreSQL/flexibleServers/example-server"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `server_id` - (Required) The ID of the PostgreSQL Flexible Server to query.
