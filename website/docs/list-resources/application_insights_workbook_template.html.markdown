---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_workbook_template"
description: |-
    Lists Application Insights Workbook Template resources.
---

# List resource: azurerm_application_insights_workbook_template

Lists Application Insights Workbook Template resources.

## Example Usage

### List all Application Insights Workbook Templates in a specific resource group

```hcl
list "azurerm_application_insights_workbook_template" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Required) The name of the resource group to query.
