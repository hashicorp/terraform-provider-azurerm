---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_dataset_snowflake"
description: |-
  Manages a Snowflake Dataset inside a Azure Data Factory.
---

# azurerm_data_factory_dataset_snowflake

Manages a Snowflake Dataset inside an Azure Data Factory.

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

resource "azurerm_data_factory_linked_service_snowflake" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
  connection_string   = "jdbc:snowflake://account.region.snowflakecomputing.com/?user=user&db=db&warehouse=wh"
}

resource "azurerm_data_factory_dataset_snowflake" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  linked_service_name = azurerm_data_factory_linked_service_snowflake.test.name

  schema_name = "foo_schema"
  table_name  = "foo_table"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Dataset Snowflake. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Dataset Snowflake. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Dataset with. Changing this forces a new resource.

* `linked_service_name` - (Required) The Data Factory Linked Service name in which to associate the Dataset with.

* `schema_name` - (Optional) The schema name of the Data Factory Dataset Snowflake.

* `table_name` - (Optional) The table name of the Data Factory Dataset Snowflake.

* `folder` - (Optional) The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.

* `schema_column` - (Optional) A `schema_column` block as defined below.

* `description` - (Optional) The description for the Data Factory Dataset Snowflake.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Dataset Snowflake.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Dataset Snowflake.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Dataset Snowflake.

---

A `schema_column` block supports the following:

* `name` - (Required) The name of the column.

* `type` - (Optional) Type of the column. Valid values are `Byte`, `Byte[]`, `Boolean`, `Date`, `DateTime`,`DateTimeOffset`, `Decimal`, `Double`, `Guid`, `Int16`, `Int32`, `Int64`, `Single`, `String`, `TimeSpan`. Please note these values are case sensitive.

* `description` - (Optional) The description of the column.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Snowflake Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Snowflake Dataset.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Snowflake Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Snowflake Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Snowflake Dataset.

## Import

Data Factory Snowflake Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_dataset_snowflake.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/datasets/example
```
