---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_storage_share"
description: |-
  Gets information about an existing File Share.
---

# Data Source: azurerm_storage_share

Use this data source to access information about an existing File Share.

~> **Note on Authentication** Shared Key authentication will always be used for this data source, as AzureAD authentication is not supported by the Storage API for files.

## Example Usage

```hcl
data "azurerm_storage_share" "example" {
  name                 = "existing"
  storage_account_name = "existing"
}

output "id" {
  value = data.azurerm_storage_share.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the share.

* `storage_account_name` - (Required) The name of the storage account.

## Attributes Reference

* `id` - The ID of the File Share.

* `quota` - The quota of the File Share in GB.

* `metadata` - A map of custom file share metadata.

* `acl` - One or more acl blocks as defined below.

---

A `acl` block has the following attributes:

* `id` - The ID which should be used for this Shared Identifier.

* `access_policy` - An `access_policy` block as defined below.

---

A `access_policy` block has the following attributes:

* `permissions` - The permissions which should be associated with this Shared Identifier. Possible value is combination of `r` (read), `w` (write), `d` (delete), and `l` (list).

* `start` - The time at which this Access Policy should be valid from, in [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) format.

* `expiry` - The time at which this Access Policy should be valid until, in [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) format.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage.
