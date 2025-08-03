---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_azure_blob_storage"
description: |-
  Manages a Linked Service (connection) between an Azure Blob Storage Account and Azure Data Factory.
---

# azurerm_data_factory_linked_service_azure_blob_storage

Manages a Linked Service (connection) between an Azure Blob Storage Account and Azure Data Factory.

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_storage_account" "example" {
  name                = "storageaccountname"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_linked_service_azure_blob_storage" "example" {
  name              = "example"
  data_factory_id   = azurerm_data_factory.example.id
  connection_string = data.azurerm_storage_account.example.primary_connection_string
}
```

## Example Usage with SAS URI and SAS Token

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_data_factory" "test" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_key_vault" "test" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_data_factory_linked_service_key_vault" "test" {
  name            = "linkkv"
  data_factory_id = azurerm_data_factory.test.id
  key_vault_id    = azurerm_key_vault.test.id
}

resource "azurerm_data_factory_linked_service_azure_blob_storage" "test" {
  name            = "example"
  data_factory_id = azurerm_data_factory.test.id

  sas_uri = "https://example.blob.core.windows.net"
  key_vault_sas_token {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}

resource "azurerm_data_factory_linked_service_azure_blob_storage" "test" {
  name            = "example"
  data_factory_id = azurerm_data_factory.test.id

  service_endpoint     = "https://example.blob.core.windows.net"
  service_principal_id = "00000000-0000-0000-0000-000000000000"
  tenant_id            = "00000000-0000-0000-0000-000000000000"
  service_principal_linked_key_vault_key {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.test.name
    secret_name         = "secret"
  }
}
```

## Argument Reference

The following supported arguments are common across all Azure Data Factory Linked Services:

* `name` - (Required) Specifies the name of the Data Factory Linked Service. Changing this forces a new resource to be created. Must be unique within a data factory. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `description` - (Optional) The description for the Data Factory Linked Service.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service.

The following supported arguments are specific to Azure Blob Storage Linked Service:

* `connection_string` - (Optional) The connection string. Conflicts with `connection_string_insecure`, `sas_uri` and `service_endpoint`.

* `connection_string_insecure` - (Optional) The connection string sent insecurely. Conflicts with `connection_string`, `sas_uri` and `service_endpoint`.

~> **Note:** `connection_string` uses the Azure [SecureString](https://learn.microsoft.com/en-us/dotnet/api/microsoft.azure.management.datafactory.models.securestring) to encrypt the contents within the REST payload sent to Azure whilst the `connection_string_insecure` is sent as a regular string. Both properties are still sent using SSL/HTTPS. At this time the portal will not decrypt Secure Strings so the `connection_string` property in the portal will show as `******` whilst `connection_string_insecure` will be viewable in the portal.

* `sas_uri` - (Optional) The SAS URI. Conflicts with `connection_string_insecure`, `connection_string` and `service_endpoint`.

* `key_vault_sas_token` - (Optional) A `key_vault_sas_token` block as defined below. Use this argument to store SAS Token in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service. A `sas_uri` is required.

---

A `key_vault_sas_token` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores the SAS token.

---

* `service_principal_linked_key_vault_key` - (Optional) A `service_principal_linked_key_vault_key` block as defined below. Use this argument to store Service Principal key in an existing Key Vault. It needs an existing Key Vault Data Factory Linked Service.

---

A `service_principal_linked_key_vault_key` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores the Service Principal key.

---

* `service_endpoint` - (Optional) The Service Endpoint. Conflicts with `connection_string`, `connection_string_insecure` and `sas_uri`.

* `use_managed_identity` - (Optional) Whether to use the Data Factory's managed identity to authenticate against the Azure Blob Storage account. Incompatible with `service_principal_id` and `service_principal_key`.

* `service_principal_id` - (Optional) The service principal id in which to authenticate against the Azure Blob Storage account.

* `service_principal_key` - (Optional) The service principal key in which to authenticate against the AAzure Blob Storage account.

* `storage_kind` - (Optional) Specify the kind of the storage account. Allowed values are `Storage`, `StorageV2`, `BlobStorage` and `BlockBlobStorage`.

* `tenant_id` - (Optional) The tenant id or name in which to authenticate against the Azure Blob Storage account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Linked Service.

## Import

Data Factory Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_azure_blob_storage.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
