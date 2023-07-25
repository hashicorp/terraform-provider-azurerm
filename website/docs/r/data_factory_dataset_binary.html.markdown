---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_dataset_binary"
description: |-
  Manages a Data Factory Binary Dataset inside an Azure Data Factory.
---

# azurerm_data_factory_dataset_binary

Manages a Data Factory Binary Dataset inside an Azure Data Factory.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_linked_service_sftp" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id

  authentication_type = "Basic"
  host                = "http://www.bing.com"
  port                = 22
  username            = "foo"
  password            = "bar"
}

resource "azurerm_data_factory_dataset_binary" "example" {
  name                = "example"
  data_factory_id     = azurerm_data_factory.example.id
  linked_service_name = azurerm_data_factory_linked_service_sftp.example.name

  sftp_server_location {
    path     = "/test/"
    filename = "**"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Binary Dataset. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `linked_service_name` - (Required) The Data Factory Linked Service name in which to associate the Binary Dataset with.

---

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Binary Dataset.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Binary Dataset.

* `compression` - (Optional) A `compression` block as defined below.

* `description` - (Optional) The description for the Data Factory Dataset.

* `folder` - (Optional) The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.

* `parameters` - (Optional) Specifies a list of parameters to associate with the Data Factory Binary Dataset.

The following supported locations for a Binary Dataset. One of these should be specified:

* `http_server_location` - (Optional) A `http_server_location` block as defined below.

* `azure_blob_storage_location` - (Optional) A `azure_blob_storage_location` block as defined below.

* `sftp_server_location` - (Optional) A `sftp_server_location` block as defined below.

---

A `compression` block supports the following:

* `type` - (Required) The type of compression used during transport. Possible values are `BZip2`, `Deflate`, `GZip`, `Tar`, `TarGZip` and `ZipDeflate`.

* `level` - (Optional) The level of compression. Possible values are `Fastest` and `Optimal`.

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

* `path` - (Optional) The folder path to the file in the blob container.

* `filename` - (Optional) The filename of the file in the blob container.

* `dynamic_container_enabled` - (Optional) Is the `container` using dynamic expression, function or system variables? Defaults to `false`.

* `dynamic_path_enabled` - (Optional) Is the `path` using dynamic expression, function or system variables? Defaults to `false`.

* `dynamic_filename_enabled` - (Optional) Is the `filename` using dynamic expression, function or system variables? Defaults to `false`.

---

A `sftp_server_location` block supports the following:

* `path` - (Required) The folder path to the file on the SFTP server.

* `filename` - (Required) The filename of the file on the SFTP server.

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

Data Factory Binary Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_dataset_binary.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/datasets/example
```
