---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_azure_postgresql"
description: |-
  Manages a Linked Service (connection) between Azure Database for PostgreSQL and Azure Data Factory.
---

# azurerm_data_factory_linked_service_azure_postgresql

Manages a Linked Service (connection) between Azure Database for PostgreSQL and Azure Data Factory.

~> **Note:** All arguments including the password will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

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

resource "azurerm_data_factory_linked_service_azure_postgresql" "example" {
  name                = "example"
  data_factory_id     = azurerm_data_factory.example.id
  authentication_type = "SystemAssignedManagedIdentity"
  server              = "example.postgres.database.azure.com"
  port                = 5432
  database_name       = "exampledb"
  ssl_mode            = "5"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Azure Database for PostgreSQL Linked Service. Changing this forces a new resource to be created. Must be unique within a data factory. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource to be created.

* `authentication_type` - (Required) The authentication type used to connect to the Azure Database for PostgreSQL server. Possible values are `Basic`, `SystemAssignedManagedIdentity`, and `UserAssignedManagedIdentity`.

~> **Note:** When `authentication_type` is set to `Basic`, `username` and `key_vault_password` are required. When `authentication_type` is set to `UserAssignedManagedIdentity`, `credential_name` is required.

* `database_name` - (Required) The name of the database on the Azure Database for PostgreSQL server instance.

* `port` - (Required) The port number used to connect to the Azure Database for PostgreSQL server instance.

* `server` - (Required) The fully qualified host name of the Azure Database for PostgreSQL server instance.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Azure Database for PostgreSQL Linked Service.

* `credential_name` - (Optional) The name of a Data Factory credential that uses a User-Assigned Managed Identity to authenticate against the Azure Database for PostgreSQL server.

~> **Note:** `credential_name` conflicts with `key_vault_password`.

* `description` - (Optional) The description for the Data Factory Azure Database for PostgreSQL Linked Service.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Azure Database for PostgreSQL Linked Service.

* `key_vault_password` - (Optional) A `key_vault_password` block as defined below. Use this argument to store the password in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Azure Database for PostgreSQL Linked Service.

* `ssl_mode` - (Optional) The SSL connection mode used to connect to the Azure Database for PostgreSQL server. Possible values are `0` (Disabled), `1` (Allow), `2` (Prefer), `3` (Require), `4` (VerifyCA), and `5` (VerifyFull).

* `username` - (Optional) The username used to authenticate against the Azure Database for PostgreSQL server.

---

A `key_vault_password` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores the password.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Azure Database for PostgreSQL Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Azure Database for PostgreSQL Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Azure Database for PostgreSQL Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Azure Database for PostgreSQL Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Azure Database for PostgreSQL Linked Service.

## Import

Data Factory Azure Database for PostgreSQL Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_azure_postgresql.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataFactory/factories/factory1/linkedservices/linkedservice1
```
