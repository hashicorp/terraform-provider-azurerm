---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_job_definition"
description: |-
  Lists Storage Mover Job Definition resources.
---

# List resource: azurerm_storage_mover_job_definition

Lists Storage Mover Job Definition resources.

## Example Usage

### List Job Definitions in a Storage Mover Project

```hcl
list "azurerm_storage_mover_job_definition" "example" {
  provider = azurerm
  config {
    storage_mover_project_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.StorageMover/storageMovers/example-mover/projects/example-project"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `storage_mover_project_id` - (Required) The ID of the Storage Mover Project to query.