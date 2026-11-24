---
subcategory: "Search"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_search_service_datasource_data_lake_storage_gen2"
description: |-
  Manages a Search Service Data Lake Storage Gen2 Data Source.
---

# azurerm_search_service_datasource_data_lake_storage_gen2

Manages a Search Service Data Lake Storage Gen2 Data Source.

## Example Usage

This example uses the system-assigned managed identity of the Search Service to access a hierarchical namespace-enabled Storage Account.

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_search_service" "example" {
  name                        = "example-search-service"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  sku                         = "standard"
  authentication_failure_mode = "http403"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_role_assignment" "search_service_contributor" {
  scope                = azurerm_search_service.example.id
  role_definition_name = "Search Service Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example-storage-data-lake-gen2-filesystem"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_role_assignment" "storage_blob_data_reader" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Reader"
  principal_id         = azurerm_search_service.example.identity[0].principal_id
}

resource "azurerm_search_service_datasource_data_lake_storage_gen2" "example" {
  name                         = "example-search-service-datasource-data-lake-storage-gen2"
  search_service_endpoint      = azurerm_search_service.example.endpoint
  connection_string_wo         = "ResourceId=${azurerm_storage_account.example.id};"
  connection_string_wo_version = 1
  filesystem_name              = azurerm_storage_data_lake_gen2_filesystem.example.name

  depends_on = [
    azurerm_role_assignment.search_service_contributor,
    azurerm_role_assignment.storage_blob_data_reader,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Search Service Data Lake Storage Gen2 Data Source. Changing this forces a new resource to be created.

~> **Note:** The `name` must be 2-128 characters long, start and end with a lowercase letter or digit, and contain only lowercase letters, digits, or dashes.

* `search_service_endpoint` - (Required) The endpoint of the Search Service in which this Data Source should be created. Changing this forces a new resource to be created.

~> **Note:** The `search_service_endpoint` must use HTTPS and must not end with a trailing `/`.

* `connection_string_wo` - (Required, Write-Only) The connection string used to access the Data Lake Storage Gen2 Storage Account.

-> **Note:** The connection string can use a Storage Account resource ID for managed identity authentication, an account key, or a supported shared access signature.

* `connection_string_wo_version` - (Required) An integer value used to trigger an update for `connection_string_wo`.

~> **Note:** `connection_string_wo_version` must be `1` or greater and must be changed when updating `connection_string_wo`. The connection string must be provided when changing this version.

* `filesystem_name` - (Required) The name of the Data Lake Storage Gen2 filesystem from which to read data.

~> **Note:** The `filesystem_name` must be 3-63 characters long, contain only lowercase letters, digits, or hyphens, and must not begin with a hyphen. The underlying Storage Account must have hierarchical namespace enabled.

* `container_query` - (Optional) The directory path within the filesystem to index.

* `description` - (Optional) A description of the Search Service Data Lake Storage Gen2 Data Source.

* `encryption_key` - (Optional) An `encryption_key` block as defined below. Changing this forces a new resource to be created.

* `soft_delete_column_name` - (Optional) The name of the custom metadata property that indicates soft deletion of a file.

* `soft_delete_marker_value` - (Optional) The value of the custom metadata property that indicates a soft-deleted file.

~> **Note:** `soft_delete_column_name` and `soft_delete_marker_value` must be specified together.

---

An `encryption_key` block supports the following:

* `key_vault_key_id` - (Required) The ID of the Key Vault Key used to encrypt the Data Source definition.

* `client_id` - (Optional) The Client ID of the Microsoft Entra ID application used to access the Key Vault Key.

~> **Note:** `client_id` must be a UUID. The `client_id` and `client_secret` arguments must be specified together.

* `client_secret` - (Optional) The Client Secret of the Microsoft Entra ID application used to access the Key Vault Key.

-> **Note:** If `client_id` and `client_secret` are not specified, the system-assigned managed identity of the Search Service is used to access the Key Vault Key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Search Service Data Lake Storage Gen2 Data Source.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Search Service Data Lake Storage Gen2 Data Source.
* `read` - (Defaults to 5 minutes) Used when retrieving the Search Service Data Lake Storage Gen2 Data Source.
* `update` - (Defaults to 30 minutes) Used when updating the Search Service Data Lake Storage Gen2 Data Source.
* `delete` - (Defaults to 30 minutes) Used when deleting the Search Service Data Lake Storage Gen2 Data Source.

## Import

A Search Service Data Lake Storage Gen2 Data Source can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_search_service_datasource_data_lake_storage_gen2.example https://searchservice1.search.windows.net/datasources/datasource1
```

-> **Note:** Only Data Sources of type `adlsgen2` can be imported. After importing, configure `connection_string_wo` and `connection_string_wo_version` before applying changes.
