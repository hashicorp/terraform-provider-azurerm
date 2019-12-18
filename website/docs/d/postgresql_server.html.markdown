---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_server"
sidebar_current: "docs-azurerm-datasource-postgresql-server"
description: |-
  Gets information about an existing PostgreSQL Azure Database Server.
---

# Data Source: azurerm_postgresql_server

Use this data source to access information about an existing PostgreSQL Azure Database Server.

## Example Usage

```hcl
data "azurerm_postgresql_server" "example" {
  name                = "postgresql-server-1"
  resource_group_name = "api-rg-pro"
}

output "postgresql_server_id" {
  value = "${data.azurerm_postgresql_server.example.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the PostgreSQL Server.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the PostgreSQL Server exists.

## Attributes Reference

* `location` - The location of the Resource Group in which the PostgreSQL Server exists.

* `fqdn` - The fully qualified domain name of the PostgreSQL Server.

* `version` - The version of the PostgreSQL Server.

* `administrator_login` - The administrator username of the PostgreSQL Server.

* `tags` - A mapping of tags assigned to the resource.
