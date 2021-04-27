---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_dataset_http"
description: |-
  Manages an Azure Delimited Text Dataset inside an Azure Data Factory.
---

# azurerm_data_factory_dataset_http

Manages an Azure HTTP Dataset inside an Azure Data Factory.

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
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
  authentication_type = "Anonymous"
  url                 = "https://www.bing.com"
}

resource "azurerm_data_factory_dataset_http" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
  linked_service_name = azurerm_data_factory_linked_service_web.example.name

  relative_url   = "http://www.bing.com"
  request_body   = "foo=bar"
  request_method = "POST"

}
```

## Argument Reference

The following supported arguments are common across all Azure Data Factory Datasets:

* `name` - (Required) Specifies the name of the Data Factory Dataset. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Dataset. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Dataset with. Changing this forces a new resource.

* `linked_service_name` - (Required) The Data Factory Linked Service name in which to associate the Dataset with.

* `folder` - (Optional) The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.

* `schema_column` - (Optional) A `schema_column` block as defined below.

* `description` - (Optional) The description for the Data Factory Dataset.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Dataset.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Dataset.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Dataset.

The following supported arguments are specific to HTTP Dataset:

* `relative_url` - (Required) The relative URL based on the URL in the HTTP Linked Service.

* `request_body` - (Required) The body for the HTTP request.

* `request_method` - (Required) The HTTP method for the HTTP request. (e.g. GET, POST)

---

A `schema_column` block supports the following:

* `name` - (Required) The name of the column.

* `type` - (Optional) Type of the column. Valid values are `Byte`, `Byte[]`, `Boolean`, `Date`, `DateTime`,`DateTimeOffset`, `Decimal`, `Double`, `Guid`, `Int16`, `Int32`, `Int64`, `Single`, `String`, `TimeSpan`. Please note these values are case sensitive.

* `description` - (Optional) The description of the column.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Dataset.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Dataset.

## Import

Data Factory Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_dataset_http.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/datasets/example
```
