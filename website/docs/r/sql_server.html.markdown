---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_server"
sidebar_current: "docs-azurerm-resource-database-sql-server"
description: |-
  Manages a SQL Azure Database Server.

---

# azurerm\_sql\_server

Manages a SQL Azure Database Server.

~> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "database-rg"
  location = "West US"
}

resource "azurerm_sql_server" "test" {
  name                         = "mysqlserver"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  tags {
    environment = "production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL Server. This needs to be globally unique within Azure.

* `resource_group_name` - (Required) The name of the resource group in which to create the SQL Server.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `version` - (Required) The version for the new server. Valid values are: 2.0 (for v11 server) and 12.0 (for v12 server).

* `administrator_login` - (Required) The administrator login name for the new server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx)

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The SQL Server ID.
* `fully_qualified_domain_name` - The fully qualified domain name of the Azure SQL Server (e.g. myServerName.database.windows.net)

## Import

SQL Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_server.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver
```
