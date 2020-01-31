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

### Timeouts

~> **Note:** Custom Timeouts are available [as an opt-in Beta in version 1.43 of the Azure Provider](/docs/providers/azurerm/guides/2.0-beta.html) and will be enabled by default in version 2.0 of the Azure Provider.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Disk Encryption Set.
