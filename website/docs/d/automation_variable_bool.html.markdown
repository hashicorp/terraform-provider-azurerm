---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_variable_bool"
sidebar_current: "docs-azurerm-datasource-automation-variable-bool"
description: |-
  Gets information about an existing Automation Bool Variable
---

# Data Source: azurerm_automation_variable_bool

Use this data source to access information about an existing Automation Bool Variable.


## Example Usage

```hcl
data "azurerm_automation_variable_bool" "example" {
  name                    = "tfex-example-var"
  resource_group_name     = "tfex-example-rg"
  automation_account_name = "tfex-example-account"
}

output "variable_id" {
  value = "${data.azurerm_automation_variable_bool.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Automation Variable.

* `resource_group_name` - (Required) The Name of the Resource Group where the automation account exists.

* `automation_account_name` - (Required) The name of the automation account in which the Automation Variable exists.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Automation Variable.

* `description` - The description of the Automation Variable.

* `encrypted` - Specifies if the Automation Variable is encrypted. Defaults to `false`.

* `value` - The value of the Automation Variable as a `boolean`.
