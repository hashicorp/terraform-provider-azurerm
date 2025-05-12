---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_container"
description: |-
  Gets information about an existing Storage Container.
---

# Data Source: azurerm_storage_container

Use this data source to access information about an existing Storage Container.

## Example Usage

```hcl
data "azurerm_storage_account" "example" {
  name                = "exampleaccount"
  resource_group_name = "examples"
}

data "azurerm_storage_container" "example" {
  name               = "example-container-name"
  storage_account_id = data.azurerm_storage_account.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Container.

* `storage_account_name` - (Optional) The name of the Storage Account where the Container exists. This property is deprecated in favour of `storage_account_id`.

* `storage_account_id` - (Optional) The id of the Storage Account where the Container exists. This property will become Required in version 5.0 of the Provider.

~> **Note:** One of `storage_account_name` or `storage_account_id` must be specified. When specifying `storage_account_id` the resource will use the Resource Manager API, rather than the Data Plane API.

## Attributes Reference

* `container_access_type` - The Access Level configured for this Container.

* `default_encryption_scope` - The default encryption scope in use for blobs uploaded to this container.

* `encryption_scope_override_enabled` - Whether blobs are allowed to override the default encryption scope for this container.

* `has_immutability_policy` - Is there an Immutability Policy configured on this Storage Container?

* `has_legal_hold` - Is there a Legal Hold configured on this Storage Container?

* `metadata`  - A mapping of MetaData for this Container.

* `id` - The Resource Manager ID of this Storage Container.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Container.
