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
  name              = "example"
  data_factory_id   = azurerm_data_factory.example.id
  connection_string = "jdbc:snowflake://account.region.snowflakecomputing.com/?user=user&db=db&warehouse=wh"
}

resource "azurerm_data_factory_dataset_snowflake" "example" {
  name                = "example"
  data_factory_id     = azurerm_data_factory.example.id
  linked_service_name = azurerm_data_factory_linked_service_snowflake.example.name

  schema_name = "foo_schema"
  table_name  = "foo_table"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Dataset Snowflake. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

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

* `type` - (Optional) Type of the column. Valid values are `NUMBER`, `DECIMAL`, `NUMERIC`, `INT`, `INTEGER`, `BIGINT`, `SMALLINT`, `FLOAT``FLOAT4`, `FLOAT8`, `DOUBLE`, `DOUBLE PRECISION`, `REAL`, `VARCHAR`, `CHAR`, `CHARACTER`, `STRING`, `TEXT`, `BINARY`, `VARBINARY`, `BOOLEAN`, `DATE`, `DATETIME`, `TIME`, `TIMESTAMP`, `TIMESTAMP_LTZ`, `TIMESTAMP_NTZ`, `TIMESTAMP_TZ`, `VARIANT`, `OBJECT`, `ARRAY`, `GEOGRAPHY`. Please note these values are case sensitive.

* `precision` - (Optional) The total number of digits allowed.

* `scale` - (Optional) The number of digits allowed to the right of the decimal point.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Snowflake Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Snowflake Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Snowflake Dataset.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Snowflake Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Snowflake Dataset.

## Import

Data Factory Snowflake Datasets can be imported using the `resource id`,  e.g.

```shell
terraform import azurerm_data_factory_dataset_snowflake.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/datasets/example
```
