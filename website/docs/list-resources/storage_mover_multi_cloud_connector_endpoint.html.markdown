---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_multi_cloud_connector_endpoint"
description: |-
  Lists Storage Mover Multi-Cloud Connector Endpoint resources.
---

# List resource: azurerm_storage_mover_multi_cloud_connector_endpoint

Lists Storage Mover Multi-Cloud Connector Endpoint resources.

## Example Usage

### List Multi-Cloud Connector Endpoints in a Storage Mover

```hcl
list "azurerm_storage_mover_multi_cloud_connector_endpoint" "example" {
  provider = azurerm
  config {
    storage_mover_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.StorageMover/storageMovers/example-mover"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `storage_mover_id` - (Required) The ID of the Storage Mover to query.
