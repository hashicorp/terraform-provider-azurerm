---
subcategory: "Data Lake"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_store_file"
description: |-
  Manages a Azure Data Lake Store File.
---

# azurerm_data_lake_store_file

Manages a Azure Data Lake Store File.

~> **Note:** If you want to change the data in the remote file without changing the `local_file_path`, then
taint the resource so the `azurerm_data_lake_store_file` gets recreated with the new data.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_lake_store" "example" {
  name                = "consumptiondatalake"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_data_lake_store_file" "example" {
  resource_group_name = azurerm_resource_group.example.name
  local_file_path     = "/path/to/local/file"
  remote_file_path    = "/path/created/for/remote/file"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) Specifies the name of the Data Lake Store for which the File should created.

* `local_file_path` - (Required) The path to the local file to be added to the Data Lake Store.

* `remote_file_path` - (Required) The path created for the file on the Data Lake Store.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Lake Store File.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Lake Store File.
* `update` - (Defaults to 30 minutes) Used when updating the Data Lake Store File.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Lake Store File.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Lake Store File.

## Import

Data Lake Store File's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_lake_store_file.exampleexample.azuredatalakestore.net/test/example.txt
```
