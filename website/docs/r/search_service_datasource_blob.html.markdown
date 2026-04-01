---
subcategory: "Search"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_search_service_datasource_blob"
description: |-
  Manages a Search Service Blob Data Source.
---

# azurerm_search_service_datasource_blob

Manages a Search Service Blob Data Source.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_search_service" "example" {
  name                = "example-search-service"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "standard"
  authentication_failure_mode = "http403"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name               = "example-storage-container"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_search_service_datasource_blob" "example" {
  name              = "example-search-service-datasource-blob"
  search_service_id = azurerm_search_service.example.id
  container_name    = azurerm_storage_container.example.name
  connection_string = azurerm_storage_account.example.primary_connection_string
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Search Service Blob Data Source. Changing this forces a new resource to be created.

~> **Note:** The `name` must be 2-128 characters long, start and end with a lowercase letter or digit, and contain only lowercase letters, digits, or dashes.

* `search_service_id` - (Required) The ID of the Search Service in which this Blob Data Source should be created. Changing this forces a new resource to be created.

* `connection_string` - (Required) The connection string to the Azure Blob Storage account.

* `container_name` - (Required) The name of the Azure Blob Storage container from which to read data.

* `container_query` - (Optional) A query string that filters the set of blobs in the container.

* `description` - (Optional) A description for this Search Service Blob Data Source.

* `encryption_key` - (Optional) An `encryption_key` block as defined below. Changing this forces a new resource to be created.

* `soft_delete_column_name` - (Optional) The name of the column that indicates soft deletion of a blob.

* `soft_delete_marker_value` - (Optional) The value in the `soft_delete_column_name` column that indicates a soft-deleted blob.

~> **Note:** `soft_delete_column_name` is required when `soft_delete_marker_value` is specified.

---

An `encryption_key` block supports the following:

* `key_name` - (Required) The name of the Key Vault key used for encryption.

* `key_vault_uri` - (Required) The URI of the Key Vault. Must be an HTTPS URL.

* `application_id` - (Optional) The Application ID of the Azure Active Directory application used for authentication.

~> **Note:** `application_id` and `application_secret` must be specified together.

* `application_secret` - (Optional) The Application Secret of the Azure Active Directory application used for authentication.

* `key_version` - (Optional) The version of the Key Vault key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Search Service Blob Data Source.

* `etag` - The ETag of the Search Service Blob Data Source.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Search Service Blob Data Source.
* `read` - (Defaults to 5 minutes) Used when retrieving the Search Service Blob Data Source.
* `update` - (Defaults to 30 minutes) Used when updating the Search Service Blob Data Source.
* `delete` - (Defaults to 30 minutes) Used when deleting the Search Service Blob Data Source.

## Import

A Search Service Blob Data Source can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_search_service_datasource_blob.example https://searchservice1.search.windows.net/datasources/datasource1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Search` - 2025-05-01
