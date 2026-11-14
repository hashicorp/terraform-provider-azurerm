---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app_function_keys"
description: |-
  Gets the Function keys of an existing Function within a Function App.
---

# Ephemeral: azurerm_function_app_function_keys

~> **Note:** Ephemeral Resources are supported in Terraform 1.10 and later.

Use this to access the Function keys of an existing Function within a Function App, without storing them in the Terraform state.

Unlike the Host keys exposed by the `azurerm_function_app_host_keys` Data Source, Function keys only grant access to a single Function.

## Example Usage

```hcl
ephemeral "azurerm_function_app_function_keys" "example" {
  function_app_id = azurerm_linux_function_app.example.id
  function_name   = "HttpTrigger1"
}

resource "azurerm_key_vault_secret" "example" {
  name             = "function-key"
  key_vault_id     = azurerm_key_vault.example.id
  value_wo         = ephemeral.azurerm_function_app_function_keys.example.keys["default"]
  value_wo_version = 1
}
```

## Arguments Reference

The following arguments are supported:

* `function_app_id` - (Required) The ID of the Function App in which the Function resides.

* `function_name` - (Required) The name of the Function.

## Attributes Reference

The following attributes are exported:

* `default_key` - The value of the `default` Function key.

* `keys` - A map of all Function keys, where the map key is the name of the Function key and the value is the key itself.
