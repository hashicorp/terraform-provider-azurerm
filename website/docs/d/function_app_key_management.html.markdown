---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app_key_management"
sidebar_current: "docs-azurerm-datasource-function-app-x"
description: |-
  Get information about a FunctionApp Key Management
---

# Data Source: azurerm_function_app

Use this data source to obtain information about a FunctionApp Key Management API.

## Example Usage

```hcl
data "azurerm_function_app_key_management" "test" {
  name                = "some-function-app-name"
  function_name       = "some-function-name"
  resource_group_name = "some-resource-group"
}

output "function_name_keys" {
  value = "${data.azurerm_function_app_key_management.test.function_keys["default"]}"
}
```

## Argument Reference

* `name` - (Required) The name of the FunctionApp.

* `resource_group_name` - (Required) The name of the Resource Group where the Function App exists.

* `function_name` - (Optional) The name of the Azure Function 

## Attributes Reference

* `function_keys` - A key-value pair of key management API function keys for the functions.

* `host_keys` - A key-value pair of key management API host keys for the function