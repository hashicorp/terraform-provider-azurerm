---
subcategory: "Azure VMware Solution"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vmware_private_cloud"
description: |-
  Lists Azure VMware Solution Private Cloud resources.
---

# List resource: azurerm_vmware_private_cloud

Lists Azure VMware Solution Private Cloud resources.

## Example Usage

### List all VMware Private Clouds in the subscription

```hcl
list "azurerm_vmware_private_cloud" "example" {
  provider = azurerm
  config {}
}
```

### List all VMware Private Clouds in a specific resource group

```hcl
list "azurerm_vmware_private_cloud" "example" {
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
