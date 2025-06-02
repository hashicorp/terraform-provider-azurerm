---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_variable_string"
description: |-
  Manages a string variable in Azure Automation.
---

# azurerm_automation_variable_string

Manages a string variable in Azure Automation

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-example-rg"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "tfex-example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_variable_string" "example" {
  name                    = "tfex-example-var"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  value                   = "Hello, Terraform Basic Test."
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Automation Variable. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Automation Variable. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Variable is created. Changing this forces a new resource to be created.

* `description` - (Optional) The description of the Automation Variable.

* `encrypted` - (Optional) Specifies if the Automation Variable is encrypted. Defaults to `false`.

* `value` - (Optional) The value of the Automation Variable as a `string`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Variable.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation String Variable.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation String Variable.
* `update` - (Defaults to 30 minutes) Used when updating the Automation String Variable.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation String Variable.

## Import

Automation String Variable can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_variable_string.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tfex-example-rg/providers/Microsoft.Automation/automationAccounts/tfex-example-account/variables/tfex-example-var
```
