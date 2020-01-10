---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_encryption_set"
description: |-
  Gets information about an existing Disk Encryption Set
---

# Data Source: azurerm_disk_encryption_set

Use this data source to access information about an existing Disk Encryption Set.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Disk Encryption Set exists.

* `resource_group_name` - (Required) The name of the Resource Group where the Disk Encryption Set exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Disk Encryption Set.

* `location` - The location where the Disk Encryption Set exists.

* `tags` - A mapping of tags assigned to the Disk Encryption Set.
