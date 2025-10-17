---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_blob_content"
description: |-
  Gets the content of an existing Storage Blob.
---

# Data Source: azurerm_storage_blob

Use this data source to access the content of an existing Storage Blob.

~> **Note:** The content of the blob will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

~> **Note:** The maxiumum blob size is limited to 2MiB to prevent large blobs from breaking state files or crashing due to insufficient memory.

## Example Usage

```hcl
data "azurerm_storage_blob_content" "example" {
  name                   = "example-blob-name"
  storage_account_name   = "example-storage-account-name"
  storage_container_name = "example-storage-container-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Blob.

* `storage_account_name` - The name of the Storage Account where the Container exists.

* `storage_container_name` - The name of the Storage Container where the Blob exists.

## Attributes Reference

* `id` - The ID of the storage blob.

* `content` - The content of the storage blob, encoded in base64. The [function](https://developer.hashicorp.com/terraform/language/functions/base64decode) `base64decode` can be used to obtain the original string.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Blob.
