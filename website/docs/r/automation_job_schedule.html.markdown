---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_job_schedule"
sidebar_current: "docs-azurerm-resource-automation-job-schedule"
description: |-
  Links an Automation Runbook and Schedule.
---

# azurerm_automation_job_schedule

Links an Automation Runbook and Schedule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-automation-account"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "tfex-automation-account"
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

resource "azurerm_automation_schedule" "example" {
  name                    = "tfex-automation-schedule"
  resource_group_name     = "${azurerm_resource_group.example.name}"
  automation_account_name = "${azurerm_automation_account.example.name}"
  frequency               = "Week"
  interval                = 1
  timezone                = "Central Europe Standard Time"
  start_time              = "2014-04-15T18:00:15+02:00"
  description             = "This is an example schedule"

  advanced_schedule {
    week_days = ["Friday"]
  }
}

resource "azurerm_automation_job_schedule" "example" {
  resource_group_name     = "${azurerm_resource_group.example.name}"
  automation_account_name = "${azurerm_automation_account.example.name}"
  schedule_name           = "${azurerm_automation_schedule.example.name}"
  runbook_name            = "${azurerm_automation_runbook.example.name}"

  parameters = {
    Connection         = "AzureRunAsConnection"
    VMCount            = 10
  }
}

```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which the Job Schedule is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the Automation Account in which the Job Schedule is created. Changing this forces a new resource to be created.

* `runbook_name` - (Required) The name of a Runbook to link to a Schedule. It needs to be in the same Automation Account as the Schedule and Job Schedule. Changing this forces a new resource to be created.

* `parameters` -  (Optional) A map of key/value pairs corresponding to the arguments that can be passed to the Runbook. Changing this forces a new resource to be created.

* `run_on` -  (Optional) Name of a Hybrid Worker Group the Runbook will be executed on. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Automation Schedule ID.

## Import

Automation Schedule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_job_schedule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/jobSchedules/10000000-1001-1001-1001-000000000001
```
