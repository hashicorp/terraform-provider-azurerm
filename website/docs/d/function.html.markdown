---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function"
sidebar_current: "docs-azurerm-datasource-function"
description: |-
  Gets information about an Azure Function hosted in an FunctionApp.
---

# Data Source: azurerm_function

Use this data source to access information about an Azure Function hosted in a Function App.

## Example Usage

```hcl
data "azurerm_functions_key" "test" {
  function_app_name   = "function_app_name_here"
  resource_group_name = "resource_group_here"
  function_name       = "name_of_function_you_want_key_for_here"
}

output "key" {
  value = "${data.azurerm_function.test.key}"
}
```

## Argument Reference

* `function_app_name` - (Required) Specifies the name of the Function App in which the Function is hosted.
* `function_name` - (Required) Specifies the name of the Function.
* `resource_group_name` - (Required) Specifies the name of the resource group the Function App is located.

## Attributes Reference

* `key` - The Function Key
* `trigger_url` - The full URL with Function Key used to trigger the function
