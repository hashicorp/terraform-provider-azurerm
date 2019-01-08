---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights"
sidebar_current: "docs-azurerm-datasource-application-insights"
description: |-
  Gets information about an existing Application Insights component.
---

# Data Source: azurerm_application_insights

Use this data source to access information about an existing Application Insights component.

## Example Usage

```hcl
data "azurerm_application_insights" "test" {
  name                = "production"
  resource_group_name = "networking"
}

output "application_insights_instrumentation_key" {
  value = "${data.azurerm_application_insights.test.instrumentation_key}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Application Insights component.
* `resource_group_name` - (Required) Specifies the name of the resource group the Application Insights component is located in.

## Attributes Reference

* `id` - The ID of the Virtual Machine.
* `instrumentation_key` - The instrumentation key of the Application Insights component.
