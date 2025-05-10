---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_dataset_json"
description: |-
  Manages an Azure JSON Dataset inside an Azure Data Factory.
---

# azurerm_data_factory_dataset_json

Manages an Azure JSON Dataset inside an Azure Data Factory.

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

resource "azurerm_data_factory_linked_service_web" "example" {
  name                = "example"
  data_factory_id     = azurerm_data_factory.example.id
  authentication_type = "Anonymous"
  url                 = "https://www.bing.com"
}

resource "azurerm_data_factory_dataset_json" "example" {
  name                = "example"
  data_factory_id     = azurerm_data_factory.example.id
  linked_service_name = azurerm_data_factory_linked_service_web.example.name

  http_server_location {
    relative_url = "/fizz/buzz/"
    path         = "foo/bar/"
    filename     = "foo.txt"
  }

  encoding = "UTF-8"

}
```

## Argument Reference

The following supported arguments are common across all Azure Data Factory Datasets:

* `name` - (Required) Specifies the name of the Data Factory Dataset. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `linked_service_name` - (Required) The Data Factory Linked Service name in which to associate the Dataset with.

* `folder` - (Optional) The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.

* `schema_column` - (Optional) A `schema_column` block as defined below.

* `description` - (Optional) The description for the Data Factory Dataset.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Dataset.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Dataset.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Dataset.

The following supported arguments are specific to JSON Dataset:

* `http_server_location` - (Optional) A `http_server_location` block as defined below.

* `azure_blob_storage_location` - (Optional) A `azure_blob_storage_location` block as defined below.

The following supported arguments are specific to Delimited Text Dataset:

* `encoding` - (Optional) The encoding format for the file.

---

A `schema_column` block supports the following:

* `name` - (Required) The name of the column.

* `type` - (Optional) Type of the column. Valid values are `Byte`, `Byte[]`, `Boolean`, `Date`, `DateTime`,`DateTimeOffset`, `Decimal`, `Double`, `Guid`, `Int16`, `Int32`, `Int64`, `Single`, `String`, `TimeSpan`. Please note these values are case sensitive.

* `description` - (Optional) The description of the column.

---

A `http_server_location` block supports the following:

* `relative_url` - (Required) The base URL to the web server hosting the file.

* `path` - (Required) The folder path to the file on the web server.

* `filename` - (Required) The filename of the file on the web server.

* `dynamic_path_enabled` - (Optional) Is the `path` using dynamic expression, function or system variables? Defaults to `false`.

* `dynamic_filename_enabled` - (Optional) Is the `filename` using dynamic expression, function or system variables? Defaults to `false`.

---

A `azure_blob_storage_location` block supports the following:

* `container` - (Required) The container on the Azure Blob Storage Account hosting the file.

* `path` - (Required) The folder path to the file on the web server.

* `filename` - (Required) The filename of the file on the web server.

* `dynamic_container_enabled` - (Optional) Is the `container` using dynamic expression, function or system variables? Defaults to `false`.

* `dynamic_path_enabled` - (Optional) Is the `path` using dynamic expression, function or system variables? Defaults to `false`.

* `dynamic_filename_enabled` - (Optional) Is the `filename` using dynamic expression, function or system variables? Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Dataset.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Dataset.

## Import

Data Factory Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_dataset_json.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/datasets/example
```
