---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_function_app_function_keys"
description: |-
  Gets the Function keys of an existing Function within a Function App.
---

# Data Source: azurerm_function_app_function_keys

Use this data source to access the Function keys of an existing Function within a Function App.

Unlike the Host keys exposed by the `azurerm_function_app_host_keys` Data Source, Function keys only grant access to a single Function.

~> **Note:** All attributes, including the key values, will be stored in the raw state as plain-text. Consider using the `azurerm_function_app_function_keys` Ephemeral Resource instead to avoid persisting the keys in state. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_function_app_function_keys" "example" {
  function_app_id = azurerm_linux_function_app.example.id
  function_name   = "HttpTrigger1"
}

output "default_key" {
  value     = data.azurerm_function_app_function_keys.example.default_key
  sensitive = true
}

output "custom_key" {
  value     = data.azurerm_function_app_function_keys.example.keys["my-custom-key"]
  sensitive = true
}
```

## Arguments Reference

The following arguments are supported:

* `function_app_id` - (Required) The ID of the Function App in which the Function resides.

* `function_name` - (Required) The name of the Function.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Function.

* `default_key` - The value of the `default` Function key.

* `keys` - A map of all Function keys, where the map key is the name of the Function key and the value is the key itself.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Function keys.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Web` - 2023-12-01
