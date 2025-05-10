---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_dsc_configuration"
description: |-
  Manages a Automation DSC Configuration.
---

# azurerm_automation_dsc_configuration

Manages a Automation DSC Configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "account1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_dsc_configuration" "example" {
  name                    = "test"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  location                = azurerm_resource_group.example.location
  content_embedded        = "configuration test {}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the DSC Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the DSC Configuration is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the DSC Configuration is created. Changing this forces a new resource to be created.

* `content_embedded` - (Required) The PowerShell DSC Configuration script.

* `location` - (Required) Must be the same location as the Automation Account. Changing this forces a new resource to be created.

* `log_verbose` - (Optional) Verbose log option.

* `description` - (Optional) Description to go with DSC Configuration.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation DSC Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation DSC Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation DSC Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Automation DSC Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation DSC Configuration.

## Import

Automation DSC Configuration's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_dsc_configuration.configuration1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/configurations/configuration1
```
