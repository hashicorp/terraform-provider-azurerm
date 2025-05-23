---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_custom_dataset"
description: |-
  Manages a Dataset inside an Azure Data Factory. This is a generic resource that supports all different Dataset Types.
---

# azurerm_data_factory_custom_dataset

Manages a Dataset inside an Azure Data Factory. This is a generic resource that supports all different Dataset Types.

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
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_data_factory_linked_custom_service" "example" {
  name                 = "example"
  data_factory_id      = azurerm_data_factory.example.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString":"${azurerm_storage_account.example.primary_connection_string}"
}
JSON
}

resource "azurerm_storage_container" "example" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_data_factory_custom_dataset" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id
  type            = "Json"

  linked_service {
    name = azurerm_data_factory_linked_custom_service.example.name
    parameters = {
      key1 = "value1"
    }
  }

  type_properties_json = <<JSON
{
  "location": {
    "container":"${azurerm_storage_container.example.name}",
    "fileName":"foo.txt",
    "folderPath": "foo/bar/",
    "type":"AzureBlobStorageLocation"
  },
  "encodingName":"UTF-8"
}
JSON

  description = "test description"
  annotations = ["test1", "test2", "test3"]
  folder      = "testFolder"

  parameters = {
    foo = "test1"
    Bar = "Test2"
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }

  schema_json = <<JSON
{
  "type": "object",
  "properties": {
    "name": {
      "type": "object",
      "properties": {
        "firstName": {
          "type": "string"
        },
        "lastName": {
          "type": "string"
        }
      }
    },
    "age": {
      "type": "integer"
    }
  }
}
JSON
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Data Factory Dataset. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Dataset with. Changing this forces a new resource.

* `linked_service` - (Required) A `linked_service` block as defined below.

* `type` - (Required) The type of dataset that will be associated with Data Factory. Changing this forces a new resource to be created.

* `type_properties_json` - (Required) A JSON object that contains the properties of the Data Factory Dataset.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Dataset.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Dataset.

* `description` - (Optional) The description for the Data Factory Dataset.

* `folder` - (Optional) The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Dataset.

* `schema_json` - (Optional) A JSON object that contains the schema of the Data Factory Dataset.

---

A `linked_service` block supports the following:

* `name` - (Required) The name of the Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

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
terraform import azurerm_data_factory_custom_dataset.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/datasets/example
```
