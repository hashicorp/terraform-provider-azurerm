---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_sql_managed_instance"
description: |-
  Manages a Linked Service (connection) between Azure SQL Managed Instance and Azure Data Factory.
---

# azurerm_data_factory_linked_service_sql_managed_instance

Manages a Linked Service (connection) between Azure SQL Managed Instance and Azure Data Factory.

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

resource "azurerm_data_factory_linked_service_sql_managed_instance" "example" {
  name              = "example"
  data_factory_id   = azurerm_data_factory.example.id
  connection_string = "Server=myserver.database.windows.net;Database=mydatabase;User ID=myuser;Password=mypassword;"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Linked Service SQL Managed Instance. Changing this forces a new resource to be created. Must be unique within a data factory. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `connection_string` - (Optional) The connection string in which to authenticate with the SQL Managed Instance. Exactly one of either `connection_string` or `key_vault_connection_string` is required.

* `description` - (Optional) The description for the Data Factory Linked Service SQL Managed Instance.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service SQL Managed Instance.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service SQL Managed Instance.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service SQL Managed Instance.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service SQL Managed Instance.

* `service_principal_id` - (Optional) The service principal id in which to authenticate against the Azure SQL Managed Instance. Required if `service_principal_key` is set.

* `service_principal_key` - (Optional) The service principal key in which to authenticate against the Azure SQL Managed Instance. Required if `service_principal_id` is set.

* `tenant` - (Optional) The tenant id or name in which to authenticate against the Azure SQL Managed Instance.

* `key_vault_connection_string` - (Optional) A `key_vault_connection_string` block as defined below. Use this argument to store SQL Managed Instance connection string in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service. Exactly one of either `connection_string` or `key_vault_connection_string` is required.

* `key_vault_password` - (Optional) A `key_vault_password` block as defined below. Use this argument to store SQL Managed Instance password in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service.

---

A `key_vault_connection_string` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores SQL Managed Instance connection string.

---

A `key_vault_password` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores SQL Managed Instance password.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Factory Linked Service SQL Managed Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Linked Service SQL Managed Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Linked Service SQL Managed Instance.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Linked Service SQL Managed Instance.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Linked Service SQL Managed Instance.

## Import

Data Factory Linked Service SQL Managed Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_sql_managed_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```