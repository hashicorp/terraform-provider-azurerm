---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_storage_sync"
description: |-
  Gets information about an existing Storage Sync.
---

# Data Source: azurerm_storage_sync

Use this data source to access information about an existing Storage Sync.

## Example Usage

```hcl
data "azurerm_storage_sync" "example" {
  name                = "existingStorageSyncName"
  resource_group_name = "existingResGroup"
}

output "id" {
  value = data.azurerm_storage_sync.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Storage Sync.

* `resource_group_name` - (Required) The name of the Resource Group where the Storage Sync exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Sync.

* `incoming_traffic_policy` - Incoming traffic policy.

* `location` - The Azure Region where the Storage Sync exists.

* `tags` - A mapping of tags assigned to the Storage Sync.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Sync.
