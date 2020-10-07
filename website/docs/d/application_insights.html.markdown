---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights"
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
  value = data.azurerm_application_insights.example.instrumentation_key
}
```

## Argument Reference

* `name` - Specifies the name of the Application Insights component.
* `resource_group_name` - Specifies the name of the resource group the Application Insights component is located in.

## Attributes Reference

* `id` - The ID of the Virtual Machine.
* `app_id` - The App ID associated with this Application Insights component.
* `application_type` - The type of the component.
* `instrumentation_key` - The instrumentation key of the Application Insights component.
* `connection_string` - The connection string of the Application Insights component. (Sensitive)
* `location` - The Azure location where the component exists.
* `retention_in_days` - The retention period in days.
* `tags` - Tags applied to the component.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Application Insights component.
