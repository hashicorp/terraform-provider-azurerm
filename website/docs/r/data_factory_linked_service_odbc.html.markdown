---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_odbc"
description: |-
  Manages a Linked Service (connection) between a Database and Azure Data Factory through ODBC protocol.
---

# azurerm_data_factory_linked_service_odbc

Manages a Linked Service (connection) between a Database and Azure Data Factory through ODBC protocol.

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

resource "azurerm_data_factory_linked_service_odbc" "anonymous" {
  name              = "anonymous"
  data_factory_id   = azurerm_data_factory.example.id
  connection_string = "Driver={SQL Server};Server=test;Database=test;Uid=test;Pwd=test;"
}

resource "azurerm_data_factory_linked_service_odbc" "basic_auth" {
  name              = "basic_auth"
  data_factory_id   = azurerm_data_factory.example.id
  connection_string = "Driver={SQL Server};Server=test;Database=test;Uid=test;Pwd=test;"
  basic_authentication {
    username = "onrylmz"
    password = "Ch4ngeM3!"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Linked Service ODBC. Changing this forces a new resource to be created. Must be unique within a data factory. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `connection_string` - (Required) The connection string in which to authenticate with ODBC.

* `basic_authentication` - (Optional) A `basic_authentication` block as defined below.

* `description` - (Optional) The description for the Data Factory Linked Service ODBC.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service ODBC.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service ODBC.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service ODBC.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service ODBC.

---

A `basic_authentication` block supports the following:

* `username` - (Required) The username which can be used to authenticate to the ODBC endpoint.

* `password` - (Required) The password associated with the username, which can be used to authenticate to the ODBC endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory ODBC Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory ODBC Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory ODBC Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory ODBC Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory ODBC Linked Service.

## Import

Data Factory ODBC Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_odbc.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
