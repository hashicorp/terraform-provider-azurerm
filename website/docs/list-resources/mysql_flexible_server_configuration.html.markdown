---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_flexible_server"
description: |-
  Lists MySQL Flexible Servers.
---

# List resource: azurerm_mysql_flexible_server_configuration

Lists MySQL Flexible Server configuration properties.

-> **Note:** This will return all available configuration properties available in the version of the deployed server, including those properties not being managed with Terraform. The count of configuration items will surpass the default `limit` of Terraform List (100), so you will need to increase the `limit` to include all configuration items.

## Example Usage

### List all configuration properties of a specific MySQL Flexible Server

```hcl
list "azurerm_mysql_flexible_server_configuration" "example" {
  provider = azurerm
  config {
    flexible_server_id = "some-mysql-flexible-serverid"
  }
  # increase limit to include all 300+ configurations
  limit = 400
}
```

## Argument Reference

This list resource supports the following arguments:

* `flexible_server_id` - (Required) The full ID of an existing Azure MySQL Flexible Server.

````
