---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_disk"
sidebar_current: "docs-azurerm-datasource-managed-disk"
description: |-
  Get information about the specified managed disk.
---

# Data Source: azurerm_managed_disk

Use this data source to access the properties of an existing Azure Managed Disk.

## Example Usage

```hcl
data "azurerm_managed_disk" "example" {
  name                = "example-datadisk"
  resource_group_name = "example-resources"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Managed Disk.
* `resource_group_name` - (Required) Specifies the name of the resource group.


## Attributes Reference

* `storage_account_type` - The storage account type for the managed disk.
* `source_uri` - The source URI for the managed disk
* `source_resource_id` - ID of an existing managed disk that the current resource was created from.
* `os_type` - The operating system for managed disk. Valid values are `Linux` or `Windows`
* `disk_size_gb` - The size of the managed disk in gigabytes.
* `tags` - A mapping of tags assigned to the resource.
* `zones` - (Optional) A collection containing the availability zone the managed disk is allocated in.

-> **Please Note**: Availability Zones are [in Preview and only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview) - as such you must be opted into the Preview to use this functionality. You can [opt into the Availability Zones Preview in the Azure Portal](http://aka.ms/azenroll).
