---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_odata"
description: |-
  Manages a Linked Service (connection) between a Database and Azure Data Factory through OData protocol.
---

# azurerm_data_factory_linked_service_odata

Manages a Linked Service (connection) between a Database and Azure Data Factory through OData protocol.

~> **Note:** All arguments including the connection_string will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

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

resource "azurerm_data_factory_linked_service_odata" "anonymous" {
  name                = "anonymous"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
  url                 = "https://services.odata.org/v4/TripPinServiceRW/People"
}

resource "azurerm_data_factory_linked_service_odata" "basic_auth" {
  name                = "basic_auth"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
  url                 = "https://services.odata.org/v4/TripPinServiceRW/People"
  basic_authentication {
    username = "emma"
    password = "Ch4ngeM3!"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Linked Service OData. Changing this forces a new resource to be created. Must be unique within a data factory. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Linked Service OData. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Linked Service with. Changing this forces a new resource.

* `url` - (Required) The URL of the OData service endpoint.

* `basic_authentication` - (Optional) A `basic_authentication` block as defined below.

* `description` - (Optional) The description for the Data Factory Linked Service OData.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service OData.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service OData.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service OData.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service OData.

---

A `basic_authentication` block supports the following:

* `username` - (Required) The username which can be used to authenticate to the OData endpoint.

* `password` - (Required) The password associated with the username, which can be used to authenticate to the OData endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory OData Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory OData Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory OData Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory OData Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory OData Linked Service.

## Import

Data Factory OData Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_odata.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
