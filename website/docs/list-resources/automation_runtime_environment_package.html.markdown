---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_runtime_environment_package"
description: |-
    Lists Automation Runtime Environment Package resources.
---

# List resource: azurerm_automation_runtime_environment_package

Lists Automation Runtime Environment Package resources.

## Example Usage

Lists Automation Runtime Environment Package resources.

~> **Note:** The Azure API doesn't currently return `content_uri`, `content_version`, `hash_algorithm` or `hash_value` for a package ([Azure/azure-rest-api-specs#41604](https://github.com/Azure/azure-rest-api-specs/issues/41604)), so these will be empty in list results. Since `content_uri` is required and forces a new resource, importing a listed package will plan a replacement until the API returns these values.

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
