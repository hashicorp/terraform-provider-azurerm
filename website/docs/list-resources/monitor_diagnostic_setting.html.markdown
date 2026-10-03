---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_diagnostic_setting"
description: |-
  Lists Diagnostic Settings for an Azure Resource.
---

# List resource: azurerm_monitor_diagnostic_setting

Lists Diagnostic Settings for an Azure Resource.

## Example Usage

### List all Diagnostic Settings on an Azure Resource

```hcl
list "azurerm_monitor_diagnostic_setting" "example" {
  provider = azurerm
  config {
    target_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.KeyVault/vaults/vault1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `target_resource_id` - (Required) The ID of the Target Azure Resource to query Diagnostic Settings for.
