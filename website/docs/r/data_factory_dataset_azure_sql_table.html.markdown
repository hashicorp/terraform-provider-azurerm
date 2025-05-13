---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_dataset_azure_sql_table"
description: |-
  Manages an Azure SQL Table Dataset inside an Azure Data Factory.
---

# azurerm_data_factory_dataset_azure_sql_table

Manages an Azure SQL Table Dataset inside an Azure Data Factory.

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

resource "azurerm_data_factory_linked_service_azure_sql_database" "example" {
  name              = "example"
  data_factory_id   = azurerm_data_factory.example.id
  connection_string = "Integrated Security=False;Data Source=test;Initial Catalog=test;Initial Catalog=test;User ID=test;Password=test"
}

resource "azurerm_data_factory_dataset_azure_sql_table" "example" {
  name              = "example"
  data_factory_id   = azurerm_data_factory.example.id
  linked_service_id = azurerm_data_factory_linked_service_azure_sql_database.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Dataset Azure SQL Table. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `linked_service_id` - (Required) The Data Factory Linked Service ID in which to associate the Dataset with.

* `schema` - (Optional) The schema name of the table in the Azure SQL Database.

* `table` - (Optional) The table name of the table in the Azure SQL Database.

* `folder` - (Optional) The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.

* `schema_column` - (Optional) A `schema_column` block as defined below.

* `description` - (Optional) The description for the Data Factory Dataset Azure SQL Table.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Dataset Azure SQL Table.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Dataset Azure SQL Table.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Dataset Azure SQL Table.

---

A `schema_column` block supports the following:

* `name` - (Required) The name of the column.

* `type` - (Optional) Type of the column. Valid values are `Byte`, `Byte[]`, `Boolean`, `Date`, `DateTime`,`DateTimeOffset`, `Decimal`, `Double`, `Guid`, `Int16`, `Int32`, `Int64`, `Single`, `String`, `TimeSpan`. Please note these values are case sensitive.

* `description` - (Optional) The description of the column.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Azure SQL Table Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Azure SQL Table Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Azure SQL Table Dataset.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Azure SQL Table Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Azure SQL Table Dataset.

## Import

Data Factory Azure SQL Table Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_dataset_azure_sql_table.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/datasets/example
```
