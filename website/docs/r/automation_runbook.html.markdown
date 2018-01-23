---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_runbook"
sidebar_current: "docs-azurerm-resource-automation-runbook"
description: |-
  Creates a new Automation Runbook.
---

# azurerm\_automation\_runbook

Creates a new Automation Runbook.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
 name = "resourceGroup1"
 location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "account1"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku {
    name = "Basic"
  }
}

resource "azurerm_automation_runbook" "example" {
  name                = "Get-AzureVMTutorial"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_automation_account.example.name}"
  log_verbose         = "true"
  log_progress        = "true"
  description         = "This is an example runbook"
  runbook_type        = "PowerShellWorkflow"
  publish_content_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Runbook. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Runbook is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the automation account in which the Runbook is created. Changing this forces a new resource to be created.

* `runbook_type` - (Required) The type of the runbook - can be either `Graph`, `GraphPowerShell`, `GraphPowerShellWorkflow`, `PowerShellWorkflow`, `PowerShell` or `Script`.

* `log_progress` - (Required) Progress log option.

* `log_verbose` -  (Required) Verbose log option.

* `publish_content_link` - (Required) The published runbook content link.

* `description` -  (Optional) A description for this credential.

`publish_content_link` supports the following:

* `uri` - (Required) The uri of the runbook content.

## Attributes Reference

The following attributes are exported:

* `id` - The Automation Runbook ID.

## Import

Automation Runbooks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_runbook.Get-AzureVMTutorial /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/runbooks/Get-AzureVMTutorial
```
