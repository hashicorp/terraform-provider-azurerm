---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_app_attach_package"
description: |-
  Lists Virtual Desktop App Attach Package resources.
---

# List resource: azurerm_virtual_desktop_app_attach_package

Lists Virtual Desktop App Attach Package resources.

## Example Usage

### List all Virtual Desktop App Attach Packages in the subscription

```hcl
list "azurerm_virtual_desktop_app_attach_package" "example" {
  provider = azurerm
  config {}
}
```

### List all Virtual Desktop App Attach Packages in a specific resource group

```hcl
list "azurerm_virtual_desktop_app_attach_package" "example" {
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
