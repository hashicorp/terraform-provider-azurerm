---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_data_lake_storage_gen2"
description: |-
  Manages a Linked Service (connection) between Data Lake Storage Gen2 and Azure Data Factory.
---

# azurerm_data_factory_linked_service_data_lake_storage_gen2

Manages a Linked Service (connection) between Data Lake Storage Gen2 and Azure Data Factory.

~> **Note:** All arguments including the `service_principal_key` will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

data "azurerm_client_config" "current" {
}

resource "azurerm_data_factory_linked_service_data_lake_storage_gen2" "example" {
  name                  = "example"
  resource_group_name   = azurerm_resource_group.example.name
  data_factory_name     = azurerm_data_factory.example.name
  service_principal_id  = data.azurerm_client_config.current.client_id
  service_principal_key = "exampleKey"
  tenant                = "11111111-1111-1111-1111-111111111111"
  url                   = "https://datalakestoragegen2"
}
```

## Argument Reference

The following supported arguments are common across all Azure Data Factory Linked Services:

* `name` - (Required) Specifies the name of the Data Factory Linked Service. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Linked Service. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Linked Service with. Changing this forces a new resource.

* `description` - (Optional) The description for the Data Factory Linked Service.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service.

The following supported arguments are specific to Data Lake Storage Gen2 Linked Service:

* `url` - (Required) The endpoint for the Azure Data Lake Storage Gen2 service.

* `use_managed_identity` - (Optional) Whether to use the Data Factory's managed identity to authenticate against the Azure Data Lake Storage Gen2 account. Incompatible with `service_principal_id` and `service_principal_key`  

* `service_principal_id` - (Optional) The service principal id in which to authenticate against the Azure Data Lake Storage Gen2 account. Required if `use_managed_identity` is true.

* `service_principal_key` - (Optional) The service principal key in which to authenticate against the Azure Data Lake Storage Gen2 account.  Required if `use_managed_identity` is true.

* `tenant` - (Required) The tenant id or name in which to authenticate against the Azure Data Lake Storage Gen2 account.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Data Lake Storage Gen2 Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Data Lake Storage Gen2 Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Data Lake Storage Gen2 Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Data Lake Storage Gen2 Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Data Lake Storage Gen2 Linked Service.

## Import

Data Factory Data Lake Storage Gen2 Linked Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_data_lake_storage_gen2.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
