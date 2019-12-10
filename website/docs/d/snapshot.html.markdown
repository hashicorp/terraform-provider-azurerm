---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_snapshot"
sidebar_current: "docs-azurerm-datasource-snapshot"
description: |-
  Get information about an existing Snapshot
---

# Data Source: azurerm_snapshot

Use this data source to access information about an existing Snapshot.

## Example Usage

```hcl
data "azurerm_snapshot" "example" {
  name                = "my-snapshot"
  resource_group_name = "my-resource-group"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Snapshot.

* `resource_group_name` - (Required) Specifies the name of the resource group the Snapshot is located in.

## Attributes Reference

* `id` - The ID of the Snapshot.

* `create_option` - How the snapshot was created.

* `source_uri` - The URI to a Managed or Unmanaged Disk.

* `source_resource_id` - The reference to an existing snapshot.

* `storage_account_id` - The ID of an storage account.

* `disk_size_gb` - The size of the Snapshotted Disk in GB.
