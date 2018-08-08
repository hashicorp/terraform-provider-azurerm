---
layout: "azurerm"
page_title: "Azure Resource manager: azurerm_sql_active_directory_administrator"
sidebar_current: "docs-azurerm-resource-database-sql-administrator"
description: |-
  Manages an Active Directory administrator on a SQL server

---

# azurerm_sql_active_directory_administrator

Allows you to set a user or group as the AD administrator for an Azure SQL server

## Example Usage

```hcl
data "azurerm_client_config" "example" {}

resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_sql_server" "example" {
  # ...
}

resource "azurerm_sql_active_directory_administrator" "example" {
  server_name         = "${azurerm_sql_server.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  login               = "sqladmin"
  tenant_id           = "${data.azurerm_client_config.example.tenant_id}"
  object_id           = "${data.azurerm_client_config.example.service_principal_object_id}"
}
```

## Argument Reference

The following arguments are supported:

* `server_name` - (Required) The name of the SQL Server on which to set the administrator. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group for the SQL server. Changing this forces a new resource to be created.

* `login` - (Required) The login name of the principal to set as the server administrator

* `object_id` - (Required) The ID of the principal to set as the server administrator

* `tenant_id` - (Required) The Azure Tenant ID

## Attributes Reference

The following attributes are exported:

* `id` - The SQL Active Directory Administrator ID.

## Import

A SQL Active Directory Administrator can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_active_directory_administrator.administrator /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/administrators/activeDirectory
```
