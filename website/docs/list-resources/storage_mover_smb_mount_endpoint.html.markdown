---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_smb_mount_endpoint"
description: |-
  Lists Storage Mover SMB Mount Endpoint resources.
---

# List resource: azurerm_storage_mover_smb_mount_endpoint

Lists Storage Mover SMB Mount Endpoint resources.

## Example Usage

### List SMB Mount Endpoints in a Storage Mover

```hcl
list "azurerm_storage_mover_smb_mount_endpoint" "example" {
  provider = azurerm
  config {
    storage_mover_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.StorageMover/storageMovers/example-mover"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `storage_mover_id` - (Required) The ID of the Storage Mover to query.
