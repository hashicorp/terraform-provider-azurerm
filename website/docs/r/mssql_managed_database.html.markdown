---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_database"
description: |-
  Manages an Azure SQL Azure Managed Database.
---

# azurerm_mssql_managed_database

Manages an Azure SQL Azure Managed Database for a SQL Managed Instance.

## Example Usage

```hcl
resource "azurerm_mssql_managed_database" "example" {
  name                = "example"
  managed_instance_id = azurerm_mssql_managed_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Database to create. Changing this forces a new resource to be created.

* `managed_instance_id` - (Required) The ID of the Azure SQL Managed Instance on which to create this Managed Database. Changing this forces a new resource to be created.

---

The following attributes are exported:

* `id` - The Azure SQL Managed Database ID.

## Import

SQL Managed Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/managedInstances/myserver/databases/mydatabase
```
