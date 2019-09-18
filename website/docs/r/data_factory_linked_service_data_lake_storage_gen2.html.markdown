---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_data_lake_storage_gen2"
sidebar_current: "docs-azurerm-resource-data-factory-linked-service-data-lake-storage-gen2"
description: |-
  Manages a Linked Service (connection) between Data Lake Storage Gen2 and Azure Data Factory.
---

# azurerm_data_factory_linked_service_data_lake_storage_gen2

Manages a Linked Service (connection) between Data Lake Storage Gen2 and Azure Data Factory.

~> **Note:** All arguments including the `service_principal_key` will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "northeurope"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

data "azurerm_client_config" "current" {}

resource "azurerm_data_factory_linked_service_data_lake_storage_gen2" "example" {
  name                  = "example"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  data_factory_name     = "${azurerm_data_factory.test.name}"
  service_principal_id  = "${data.azurerm_client_config.current.client_id}"
  service_principal_key = "exampleKey"
  tenant                = "11111111-1111-1111-1111-111111111111"
  url                   = "https://datalakestoragegen2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Linked Service MySQL. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Linked Service MySQL. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Linked Service with. Changing this forces a new resource.

* `url` - (Required) The endpoint for the Azure Data Lake Storage Gen2 service.

* `service_principal_id` - (Required) The service principal id in which to authenticate against the Azure Data Lake Storage Gen2 account.

* `service_principal_key` - (Required) The service principal key in which to authenticate against the Azure Data Lake Storage Gen2 account.

* `tenant` - (Required) The tenant id or name in which to authenticate against the Azure Data Lake Storage Gen2 account.

* `description` - (Optional) The description for the Data Factory Linked Service MySQL.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service MySQL.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service MySQL.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service MySQL.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service MySQL.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Linked Service.

## Import

Data Factory Linked Service MySQL can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_mysql.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
