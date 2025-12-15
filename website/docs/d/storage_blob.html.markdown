---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_blob"
description: |-
  Gets information about an existing Storage Blob.
---

# Data Source: azurerm_storage_blob

Use this data source to access information about an existing Storage Blob.

## Example Usage

```hcl
data "azurerm_storage_blob" "example" {
  name                   = "example-blob-name"
  storage_account_name   = "example-storage-account-name"
  storage_container_name = "example-storage-container-name"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of the Blob.

* `storage_account_name` - The name of the Storage Account where the Container exists.

* `storage_container_name` - The name of the Storage Container where the Blob exists.

* `is_content_sensitive` - Whether the (human readable) content of this blob will be populated to `content` or `sensitive_content` (when set to `true`)? Defaults to `false`.

## Attributes Reference

* `id` - The ID of the storage blob.

* `url` - The URL of the storage blob.

* `type` - The type of the storage blob

* `access_tier` - The access tier of the storage blob.

* `content` - The content stored in the storage blob.

* `sensitive_content` - The sensitive content stored in the storage blob.


~> **Note:** `content`/`sensitive_content` is only populated when the `content_type` is one of `text/*`, `application/{json | ld+json | x-httpd-php | xhtml+xml | x-csh | x-sh | xml | atom+xml | x-sql | yaml}`.

* `content_type` - The content type of the storage blob.

* `content_md5` - The MD5 sum of the blob contents.

* `encryption_scope` - The encryption scope for this blob.

* `metadata` - A map of custom blob metadata.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Blob.
