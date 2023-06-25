---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_watcher"
description: |-
  Manages an Automation Watcher.
---

# azurerm_automation_watcher

Manages an Automation Wacher.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_hybrid_runbook_worker_group" "example" {
  name                    = "example"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
}

resource "azurerm_automation_runbook" "example" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  log_verbose             = "true"
  log_progress            = "true"
  description             = "This is an example runbook"
  runbook_type            = "PowerShellWorkflow"

  publish_content_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/c4935ffb69246a6058eb24f54640f53f69d3ac9f/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
  }
}

resource "azurerm_automation_watcher" "example" {
  name                           = "example"
  automation_account_id          = azurerm_automation_account.example.id
  location                       = "West Europe"
  script_name                    = azurerm_automation_runbook.example.name
  script_run_on                  = azurerm_automation_hybrid_runbook_worker_group.example.name
  description                    = "example-watcher desc"
  execution_frequency_in_seconds = 42

  tags = {
    "foo" = "bar"
  }

  script_parameters = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `automation_account_id` - (Required) The ID of Automation Account to manage this Watcher. Changing this forces a new Watcher to be created.

* `name` - (Required) The name which should be used for this Automation Watcher. Changing this forces a new Automation Watcher to be created.

* `execution_frequency_in_seconds` - (Required) Specify the frequency at which the watcher is invoked.

* `location` - (Required) The Azure Region where the Automation Watcher should exist. Changing this forces a new Automation Watcher to be created.

* `script_name` - (Required) Specify the name of an existing runbook this watcher is attached to. Changing this forces a new Automation to be created.

* `script_run_on` - (Required) Specify the name of the Hybrid work group the watcher will run on.

---

* `description` - (Optional) A description of this Automation Watcher.

* `etag` - (Optional) A string of etag assigned to this Automation Watcher.

* `script_parameters` - (Optional) Specifies a list of key-vaule parameters. Changing this forces a new Automation watcher to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Automation Watcher.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Watcher.

* `status` - The current status of the Automation Watcher.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Watcher.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Watcher.
* `update` - (Defaults to 10 minutes) Used when updating the Automation Watcher.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation Watcher.

## Import

Automation Watchers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_watcher.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/watchers/watch1
```
