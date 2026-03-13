---
subcategory: "Azure VMware Solution"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vmware_netapp_volume_attachment"
description: |-
    Lists VMware NetApp Volume Attachment resources.
---

# List resource: azurerm_vmware_netapp_volume_attachment

Lists VMware NetApp Volume Attachment resources.

## Example Usage

### List all NetApp Volume Attachments in a VMware Cluster

```hcl
list "azurerm_vmware_netapp_volume_attachment" "example" {
  provider = azurerm
  config {
    vmware_cluster_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AVS/privateClouds/myPrivateCloud/clusters/Cluster-1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `vmware_cluster_id` - (Required) The ID of the VMware Cluster to query.
