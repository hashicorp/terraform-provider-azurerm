---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_storage_sync_group"
description: |-
  Gets information about an existing Storage Sync Group.
---

# Data Source: azurerm_storage_sync_group

Use this data source to access information about an existing Storage Sync Group.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_storage_sync_group" "example" {
  name            = "existing-ss-group"
  storage_sync_id = "existing-ss-id"
}

output "id" {
  value = data.azurerm_storage_sync_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Storage Sync Group.

* `storage_sync_id` - (Required) The resource ID of the Storage Sync where this Storage Sync Group is.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Sync Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Sync Group.
