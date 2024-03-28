---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_azure_postgresql"
description: |-
  Manages a Linked Service (connection) between Azure PostgreSQL Database and Azure Data Factory.
---

# azurerm_data_factory_linked_service_azure_postgresql

Manages a Linked Service (connection) between Azure PostgreSQL Database and Azure Data Factory.

~> **Note:** All arguments including the connection_string will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_linked_service_azure_postgresql" "test" {
  name              = "acctestlsazpsql%d"
  data_factory_id   = azurerm_data_factory.test.id
  connection_string = "Host=serverHostname;Port=5432;Database=postgres;UID=psqladmin@serverHostname;EncryptionMethod=0;validateservercertificate=1;Password=password123"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Linked Service Azure PostgreSQL Database. Changing this forces a new resource to be created. Must be unique within a data factory. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

---

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service Azure PostgreSQL Database.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service Azure PostgreSQL Database.

* `connection_string` - (Optional) The connection string in which to authenticate with Azure PostgreSQL Database. Exactly one of either `connection_string` or `key_vault_connection_string` is required.

* `description` - (Optional) The description for the Data Factory Linked Service Azure PostgreSQL Database.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service Azure PostgreSQL Database.

* `key_vault_connection_string` - (Optional) A `key_vault_connection_string` block as defined below. Use this argument to store Azure PostgreSQL Database connection string in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service. Exactly one of either `connection_string` or `key_vault_connection_string` is required.

* `key_vault_password` - (Optional) A `key_vault_password` block as defined below. Use this argument to store PostgreSQL Server password in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service Azure PostgreSQL Database.

---

A `key_vault_connection_string` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores PostgreSQL Server connection string.

---

A `key_vault_password` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores PostgreSQL Server password.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Factory Azure PostgreSQL Database Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Data Factory Azure PostgreSQL Database Linked Service
* `update` - (Defaults to 5 minutes) Used when updating the Data Factory Azure PostgreSQL Database Linked Service
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Azure PostgreSQL Database Linked Service
* `delete` - (Defaults to 5 minutes) Used when deleting the Data Factory Azure PostgreSQL Database Linked Service

## Import

Data Factory Azure PostgreSQL Database Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_azure_postgresql.example terraform import azurerm_data_factory_linked_service_azure_postgresql.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```