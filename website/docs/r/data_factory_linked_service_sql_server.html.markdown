---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_sql_server"
sidebar_current: "docs-azurerm-resource-data-factory-linked-service-sql-server"
description: |-
  Manages a Linked Service (connection) between a SQL Server and Azure Data Factory.
---

# azurerm_data_factory_linked_service_sql_server

Manages a Linked Service (connection) between a SQL Server and Azure Data Factory.

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

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

resource "azurerm_data_factory_linked_service_sql_server" "example" {
  name                = "example"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory.test.name}"
  connection_string   = "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test;Password=test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Linked Service SQL Server. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Linked Service SQL Server. Changing this forces a new resource

* `data_factory_name` - (Required) The Data Factory name in which to associate the Linked Service with. Changing this forces a new resource.

* `connection_string` - (Required) The connection string in which to authenticate with the SQL Server.

* `description` - (Optional) The description for the Data Factory Linked Service SQL Server.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service SQL Server.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service SQL Server.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service SQL Server.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service SQL Server.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Linked Service.

## Import

Data Factory Linked Service SQL Server can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_sql_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
