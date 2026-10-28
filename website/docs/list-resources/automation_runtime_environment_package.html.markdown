---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azurerm_automation_runtime_environment_package"
description: |-
    Lists Automation Runtime Environment Package resources.
---

# List resource: azurerm_automation_runtime_environment_package

Lists Automation Runtime Environment Package resources.

## Example Usage

### List all Automation Runtime Environment Packages in an Automation Runtime Environment

```hcl
list "azurerm_automation_runtime_environment_package" "example" {
  provider = azurerm
  config {
    automation_runtime_environment_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Automation/automationAccounts/example-account/runtimeEnvironments/example-environment"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `automation_runtime_environment_id` - (Required) The ID of the Automation Runtime Environment to query.
