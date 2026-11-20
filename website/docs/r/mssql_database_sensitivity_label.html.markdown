---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database_sensitivity_label"
description: |-
  Manages a sensitivity label on an MSSQL Database column.
---

# azurerm_mssql_database_sensitivity_label

Manages a sensitivity label on an MSSQL Database column.

-> **Note:** The schema, table, and column must already exist in the database. The following example uses the `AdventureWorksLT` sample to create them.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_mssql_server" "example" {
  name                         = "example-server"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "ChangeMe123!"
}

resource "azurerm_mssql_database" "example" {
  name        = "example-db"
  server_id   = azurerm_mssql_server.example.id
  sample_name = "AdventureWorksLT"
  sku_name    = "Basic"
}

resource "azurerm_mssql_database_sensitivity_label" "example" {
  column_name      = "EmailAddress"
  database_id      = azurerm_mssql_database.example.id
  schema_name      = "SalesLT"
  table_name       = "Customer"
  label_name       = "Confidential"
  information_type = "Email"
  rank             = "High"
}
```

## Arguments Reference

The following arguments are supported:

* `column_name` - (Required) The name of the column to classify. Must be between 1 and 128 characters in length. Changing this forces a new resource to be created.

* `database_id` - (Required) The ID of the MSSQL Database containing the column. Changing this forces a new resource to be created.

* `schema_name` - (Required) The name of the schema containing the table. Must be between 1 and 128 characters in length. Changing this forces a new resource to be created.

* `table_name` - (Required) The name of the table containing the column. Must be between 1 and 128 characters in length. Changing this forces a new resource to be created.

---

* `information_type` - (Optional) The type of information stored in the column. Must be between 1 and 128 characters in length.

* `information_type_id` - (Optional) The UUID of the information type.

* `label_id` - (Optional) The UUID of the sensitivity label.

* `label_name` - (Optional) The name of the sensitivity label assigned to the column. Must be between 1 and 128 characters in length.

* `rank` - (Optional) The sensitivity rank of the column. Possible values are `Critical`, `High`, `Low`, `Medium`, and `None`. Defaults to `None`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MSSQL Database Sensitivity Label.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MSSQL Database Sensitivity Label.
* `read` - (Defaults to 5 minutes) Used when retrieving the MSSQL Database Sensitivity Label.
* `update` - (Defaults to 30 minutes) Used when updating the MSSQL Database Sensitivity Label.
* `delete` - (Defaults to 30 minutes) Used when deleting the MSSQL Database Sensitivity Label.

## Import

MSSQL Database Sensitivity Labels can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_database_sensitivity_label.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resource-group/providers/Microsoft.Sql/servers/example-server/databases/example-db/schemas/SalesLT/tables/Customer/columns/EmailAddress/sensitivityLabels/current
```

-> **Note:** The resource ID must end in `/sensitivityLabels/current`. Recommended sensitivity labels cannot be imported.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Sql` - 2025-01-01
