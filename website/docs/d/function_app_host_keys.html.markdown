---
subcategory: 'App Service (Web Apps)'
layout: 'azurerm'
page_title: 'Azure Resource Manager: azurerm_function_app_host_keys'
description: |-
  Gets HostKeys of an existing Function App.
---

# Data Source: azurerm_function_app_host_keys

Use this data source to fetch the Hostkeys of an existing Function App

## Example Usage

```hcl
data "azurerm_function_app_host_keys" "example" {
  name                = "test-azure-functions"
  resource_group_name = azurerm_resource_group.example.name
}
```

~> **Note:** The unencrypted value of FunctionKeys, MasterKey and SystemKeys will be stored in the raw state as plain-text.

## Argument Reference

The following arguments are supported:

- `name` - The name of the Function App resource.

- `resource_group_name` - The name of the Resource Group where the Function App exists.

## Attributes Reference

The following arguments are supported:

- `function_keys` - A key-value pair of Host level function keys.

- `system_keys` - A key-value pair of system keys

- `master_key` - Secret key
