---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_job_schedule"
description: |-
  Links an Automation Runbook and Schedule.
---

# azurerm_automation_job_schedule

Links an Automation Runbook and Schedule.

~> **Note:** AzureRM provides this stand-alone [azurerm_automation_job_schedule](automation_job_schedule.html.markdown) and an inlined `job_schedule` property in [azurerm_runbook](automation_runbook.html.markdown) to manage the job schedules. You can only make use of one of these methods to manage a job schedule.

## Example Usage

This is an example of just the Job Schedule. A full example of the `azurerm_automation_job_schedule` resource can be found in [the `./examples/automation-account` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/automation-account)

```hcl
resource "azurerm_automation_job_schedule" "example" {
  resource_group_name     = "tf-rgr-automation"
  automation_account_name = "tf-automation-account"
  schedule_name           = "hour"
  runbook_name            = "Get-VirtualMachine"

  parameters = {
    resourcegroup = "tf-rgr-vm"
    vmname        = "TF-VM-01"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which the Job Schedule is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the Automation Account in which the Job Schedule is created. Changing this forces a new resource to be created.

* `runbook_name` - (Required) The name of a Runbook to link to a Schedule. It needs to be in the same Automation Account as the Schedule and Job Schedule. Changing this forces a new resource to be created.

* `schedule_name` - (Required) The name of the Schedule. Changing this forces a new resource to be created.

* `parameters` - (Optional) A map of key/value pairs corresponding to the arguments that can be passed to the Runbook. Changing this forces a new resource to be created.

-> **Note:** The parameter keys/names must strictly be in lowercase, even if this is not the case in the runbook. This is due to a limitation in Azure Automation where the parameter names are normalized. The values specified don't have this limitation.

* `run_on` - (Optional) Name of a Hybrid Worker Group the Runbook will be executed on. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Job Schedule. The format of the ID is `azurerm_automation_account.id|azurerm_automation_runbook.id`. There is an example in the [#Import](#import) part.

* `job_schedule_id` - The UUID identifying the Automation Job Schedule.

* `resource_manager_id` - The Resource Manager ID of the Automation Job Schedule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Job Schedule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Job Schedule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Job Schedule.

## Import

Automation Job Schedules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_job_schedule.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/schedules/schedule1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/runbooks/runbook1"
```
