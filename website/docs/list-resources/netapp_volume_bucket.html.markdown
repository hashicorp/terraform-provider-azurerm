---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_bucket"
description: |-
  Lists NetApp Files Volume Bucket resources.
---

# List resource: azurerm_netapp_volume_bucket

Lists NetApp Files Volume Bucket resources attached to a given NetApp Volume.

## Example Usage

### List Buckets on a NetApp Volume

```hcl
list "azurerm_netapp_volume_bucket" "example" {
  provider = azurerm
  config {
    volume_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.NetApp/netAppAccounts/example-account/capacityPools/example-pool/volumes/example-volume"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `volume_id` - (Required) The ID of the parent NetApp Volume to query for Buckets.
