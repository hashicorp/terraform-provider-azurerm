---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_managed_database"
description: |-
  Manages a SQL Azure Managed Database.
---

# azurerm_sql_managed_database

Manages a SQL Azure Managed Database.

-> **Note:** The `azurerm_sql_managed_database` resource is deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use the [`azurerm_mssql_managed_database`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_managed_database) resource instead.

## Example Usage

```hcl
resource "azurerm_sql_managed_database" "example" {
  name                    = "exampledatabase"
  sql_managed_instance_id = azurerm_sql_managed_instance.example.id
  location                = azurerm_resource_group.example.location
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL Managed Instance. Changing this forces a new resource to be created.

* `sql_managed_instance_id` - (Required) The SQL Managed Instance ID that this Managed Database will be associated with. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

---

The following attributes are exported:

* `id` - The SQL Managed Database ID.

## Import

SQL Managed Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_managed_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/managedInstances/myserver/databases/mydatabase
```
