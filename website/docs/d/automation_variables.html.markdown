---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_variables"
description: |-
  Gets all variables in an Automation Account
---

# Data Source: azurerm_automation_variables

Use this data source to get all variables in an Automation Account.

## Example Usage

```hcl
data "azurerm_automation_variables" "example" {
  resource_group_name     = "tfex-example-rg"
  automation_account_name = "tfex-example-account"
}

output "string_vars" {
  value = data.azurerm_automation_variable_string.example.string
}
```

## Argument Reference

The following arguments are supported:

- `resource_group_name` - The Name of the Resource Group where the automation account exists.

- `automation_account_name` - The name of the automation account in which the automation variables exist.

## Attributes Reference

In addition to the argument listed above, the following attributes are exported:

- `bool` - One or more `variable` blocks as defined below for each boolean variable.

- `datetime` - One or more `variable` blocks as defined below for each boolean variable.

- `encrypted` - One or more `variable` blocks as defined below for each boolean variable.

- `int` - One or more `variable` blocks as defined below for each boolean variable.

- `null` - One or more `variable` blocks as defined below for each boolean variable.

- `string` - One or more `variable` blocks as defined below for each boolean variable.

---

A `variable` block exports the following attributes:

- `name` - The name of the Automation Variable.

- `description` - The description of the Automation Variable.

- `encrypted` - Specifies if the Automation Variable is encrypted.

- `value` - The value of the Automation Variable.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

- `read` - (Defaults to 5 minutes) Used when retrieving the Automation String Variable.
