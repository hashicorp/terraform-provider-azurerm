---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_standard_web_test"
description: |-
    Lists Application Insights Standard Web Test resources.
---

# List resource: azurerm_application_insights_standard_web_test

Lists Application Insights Standard Web Test resources.

## Example Usage

### List all Application Insights Standard Web Tests in the subscription

```hcl
list "azurerm_application_insights_standard_web_test" "example" {
  provider = azurerm
  config {}
}
```

### List all Application Insights Standard Web Tests in a specific resource group

```hcl
list "azurerm_application_insights_standard_web_test" "example" {
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
