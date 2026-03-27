---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hpc_cache_nfs_target"
description: |-
  Lists NFS Target resources within a HPC Cache.
---

# List resource: azurerm_hpc_cache_nfs_target

Lists NFS Target resources within a HPC Cache.

!> **Note:** The `azurerm_hpc_cache_nfs_target` resource has been deprecated because the service is retiring on 2025-09-30. This resource will be removed in v5.0 of the AzureRM Provider. See https://aka.ms/hpccacheretirement for more information.

## Example Usage

```hcl
list "azurerm_hpc_cache_nfs_target" "example" {
  provider = azurerm

  config {
    cache_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.StorageCache/caches/example-cache"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cache_id` - (Required) The ID of the HPC Cache. The list operation returns NFS targets for this cache.