---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_runtime_environment"
description: |-
  Manages a Automation Runtime Environment.
---

# azurerm_automation_runtime_environment

Manages a Automation Runtime Environment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "%[2]s"
}

resource "azurerm_automation_account" "example" {
  name                = "accexample"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runtime_environment" "example" {
  name                    = "powershell_environment_custom_config"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name

  runtime_language = "PowerShell"
  runtime_version  = "7.2"

  location    = azurerm_resource_group.example.location
  description = "Thats my test subscription"

  runtime_default_packages = {
    "az"        = "11.2.0"
    "azure cli" = "2.56.0"
  }

  tags = {
    key = "foo"
  }
}

```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which the Automation Runtime Environment is created. Changing this forces a new resource to be created.

* `name` - (Required) The name for the Automation Runtime Environment. Changing this forces a new Automation Runtime Environment to be created.

* `location` - (Required) The location where the Automation Runtime Environment is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Automation Runtime Environment is created. Changing this forces a new resource to be created.

* `runtime_language` - (Required) The programming language used by the Automation Runtime Environment. Possible values are `Python` and `PowerShell`. Changing this forces a new Automation Runtime Environment to be created.

* `runtime_version` - (Required) The version of the runtime environment. Changing this forces a new Automation Runtime Environment to be created.

---
* `description` - (Optional) A description of the Automation Runtime Environment.

* `runtime_default_packages` - (Optional) A mapping of default packages to be installed in the Automation Runtime Environment. The default packages can only be used with PowerShell runtime environments. Removing packages will force a new Automation Runtime Environment, adding new packages will update the existing Automation Runtime Environment.

* `tags` - (Optional) A mapping of tags which should be assigned to the Automation Runtime Environment.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Runtime Environment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Automation Runtime Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Runtime Environment.
* `update` - (Defaults to 10 minutes) Used when updating the Automation Runtime Environment.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation Runtime Environment.

## Import

Automation Runtime Environments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_runtime_environment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/runtimeEnvironments/env1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Automation` - 2024-10-23
