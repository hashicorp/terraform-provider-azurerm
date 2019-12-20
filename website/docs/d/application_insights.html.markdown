---
subcategory: "Application Insights"
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
data "azurerm_application_insights" "example" {
  name                = "production"
  resource_group_name = "networking"
}

output "application_insights_instrumentation_key" {
  value = "${data.azurerm_application_insights.example.instrumentation_key}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Application Insights component.
* `resource_group_name` - (Required) Specifies the name of the resource group the Application Insights component is located in.

## Attributes Reference

* `id` - The ID of the Virtual Machine.
* `app_id` - The App ID associated with this Application Insights component.
* `application_type` - The type of the component.
* `instrumentation_key` - The instrumentation key of the Application Insights component.
* `location` - The Azure location where the component exists.
* `tags` - Tags applied to the component.
