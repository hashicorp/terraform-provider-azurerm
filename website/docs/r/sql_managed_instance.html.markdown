---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_managed_instance"
sidebar_current: "docs-azurerm-resource-database-sql-managed-instance"
description: |-
  Manages a SQL Azure Managed Instance.

---

# azurerm_sql_managed_instance

Manages a SQL Azure Managed Instance.

~> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "database-rg"
  location = "West Europe"
}

resource "azurerm_sql_managed_instance" "test" {
  name                         = "misqlserver"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  license_type                 = "BasePrice"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsJpm81"
  subnet_id                    = "${azurerm_subnet.test.id}"

  tags {
    environment = "production"
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL Managed Instance. This needs to be globally unique within Azure.

* `resource_group_name` - (Required) The name of the resource group in which to create the SQL Server.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) One `sku` blocks as defined below. 

* `vcores` - (Optional) Number of cores that should be assigned to your instance. Values can be 8, 16, or 24 if you select GP_Gen4 sku name, or 8, 16, 24, 32, or 40 if you select GP_Gen5.

* `storage_size_in_gb` - (Optional) Maximum storage space for your instance. It should be multiple of 32GB.

* `license_type` - License of the Managed Instance. Values can be PriceIncluded or BasePrice.

* `administrator_login` - (Required) The administrator login name for the new server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx)

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) Sku of the managed instance. Values can be GP_Gen4 or GP_Gen5.

## Attributes Reference

The following attributes are exported:

* `id` - The SQL Managed Instance ID.
* `fully_qualified_domain_name` - The fully qualified domain name of the Azure SQL Server (e.g. myServerName.database.windows.net)

## Import

SQL Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_managed_instance.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/managedInstances/myserver
```
