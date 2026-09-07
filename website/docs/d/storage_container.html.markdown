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

## Arguments Reference

The following arguments are supported:

* `name` - The name of the Container.

* `storage_account_id` - (Required) The ID of the Storage Account where the Container exists.

## Attributes Reference

* `container_access_type` - The Access Level configured for this Container.

* `default_encryption_scope` - The default encryption scope in use for blobs uploaded to this container.

* `encryption_scope_override_enabled` - Whether blobs are allowed to override the default encryption scope for this container.

* `has_immutability_policy` - Is there an Immutability Policy configured on this Storage Container?

* `has_legal_hold` - Is there a Legal Hold configured on this Storage Container?

* `metadata`  - A mapping of MetaData for this Container.

* `id` - The Resource Manager ID of this Storage Container.

* `url` - The data plane URL of the Storage Container in the format of `<storage blob endpoint>/<container name>`. E.g. `https://example.blob.core.windows.net/mycontainer`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Container.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Storage` - 2025-08-01
