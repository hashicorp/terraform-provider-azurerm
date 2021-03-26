---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_dataset_sql_server_table"
description: |-
  Manages a SQL Server Table Dataset inside a Azure Data Factory.
---

# azurerm_data_factory_dataset_sql_server

Manages a SQL Server Table Dataset inside a Azure Data Factory.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_linked_service_sql_server" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
  connection_string   = "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test;Password=test"
}

resource "azurerm_data_factory_dataset_sql_server_table" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
  linked_service_name = azurerm_data_factory_linked_service_sql_server.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Dataset SQL Server Table. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Dataset SQL Server Table. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Dataset with. Changing this forces a new resource.

* `linked_service_name` - (Required) The Data Factory Linked Service name in which to associate the Dataset with.

* `table_name` - (Optional) The table name of the Data Factory Dataset SQL Server Table.

* `folder` - (Optional) The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.

* `schema_column` - (Optional) A `schema_column` block as defined below.

* `description` - (Optional) The description for the Data Factory Dataset SQL Server Table.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Dataset SQL Server Table.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Dataset SQL Server Table.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Dataset SQL Server Table.

---

A `schema_column` block supports the following:

* `name` - (Required) The name of the column.

* `type` - (Optional) Type of the column. Valid values are `Byte`, `Byte[]`, `Boolean`, `Date`, `DateTime`,`DateTimeOffset`, `Decimal`, `Double`, `Guid`, `Int16`, `Int32`, `Int64`, `Single`, `String`, `TimeSpan`. Please note these values are case sensitive.

* `description` - (Optional) The description of the column.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory SQL Server Table Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory SQL Server Table Dataset.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory SQL Server Table Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory SQL Server Table Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory SQL Server Table Dataset.

## Import

Data Factory SQL Server Table Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_dataset_sql_server_table.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/datasets/example
```
