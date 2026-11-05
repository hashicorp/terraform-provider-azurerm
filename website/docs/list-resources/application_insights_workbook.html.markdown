---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_workbook"
description: |-
    Lists Application Insights Workbook resources.
---

# List resource: azurerm_application_insights_workbook

Lists Application Insights Workbook resources.

## Example Usage

### List all Application Insights Workbooks in the subscription

```hcl
list "azurerm_application_insights_workbook" "example" {
  provider = azurerm
  config {}
}
```

### List all Application Insights Workbooks in a specific resource group

```hcl
list "azurerm_application_insights_workbook" "example" {
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
