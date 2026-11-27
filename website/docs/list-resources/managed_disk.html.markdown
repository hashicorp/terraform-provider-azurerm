---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_disk"
description: |-
    Lists Managed Disk resources.
---

# List resource: azurerm_managed_disk

Lists Managed Disk resources.

## Example Usage

### List all Managed Disks in the subscription

```hcl
list "azurerm_managed_disk" "example" {
  provider = azurerm
  config {}
}
```

### List all Managed Disks in a specific resource group

```hcl
list "azurerm_managed_disk" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
