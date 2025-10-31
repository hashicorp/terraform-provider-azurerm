---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_dataset_cosmosdb_mongoapi"
description: |-
  Manages an Azure Cosmos DB Mongo API Dataset inside an Azure Data Factory.
---

# azurerm_data_factory_dataset_cosmosdb_mongoapi

Manages an Azure Cosmos DB Mongo API Dataset inside an Azure Data Factory.

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

resource "azurerm_data_factory_linked_service_cosmosdb_mongoapi" "example" {
  name              = "example"
  data_factory_id   = azurerm_data_factory.example.id
  connection_string = "mongodb://accname:secretpass@foobar.documents.azure.com:10255"
  database          = "foo"
}

resource "azurerm_data_factory_dataset_cosmosdb_mongoapi" "example" {
  name                = "example"
  collection_name     = "bar"
  data_factory_id     = azurerm_data_factory.example.id
  linked_service_name = azurerm_data_factory_linked_service_cosmosdb_mongoapi.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `collection_name` - (Required) The collection name of the Data Factory Dataset Azure Cosmos DB Mongo API.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `linked_service_name` - (Required) The Data Factory Linked Service name in which to associate the Dataset with.

* `name` - (Required) Specifies the name of the Data Factory Dataset. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

---

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Dataset.

* `description` - (Optional) The description for the Data Factory Dataset.

* `folder` - (Optional) The folder that this Dataset is in. If not specified, the Dataset will appear at the root level.

* `parameters` - (Optional) A map of string key-value pairs of parameters to associate with the Data Factory Dataset. Only string parameter type is supported at the moment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Dataset CosmosDB Mongo API.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Dataset CosmosDB Mongo API.
* `read` - (Defaults to 30 minutes) Used when retrieving the Data Factory Dataset CosmosDB Mongo API.
* `update` - (Defaults to 5 minutes) Used when updating the Data Factory Dataset CosmosDB Mongo API.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Dataset CosmosDB Mongo API.

## Import

Data Factory Dataset CosmosDB Mongo APIs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_dataset_cosmosdb_mongoapi.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/datasets/example
```
