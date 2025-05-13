---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_runbook"
description: |-
  Manages a Automation Runbook.
---

# azurerm_automation_runbook

Manages a Automation Runbook.

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

  sku_name = "Basic"
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
```

## Example Usage - custom content

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "account1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "Basic"
}

data "local_file" "example" {
  filename = "${path.module}/example.ps1"
}

resource "azurerm_automation_runbook" "example" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  log_verbose             = "true"
  log_progress            = "true"
  description             = "This is an example runbook"
  runbook_type            = "PowerShell"

  content = data.local_file.example.content
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Runbook. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Runbook is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Runbook is created. Changing this forces a new resource to be created.

* `runbook_type` - (Required) The type of the runbook - can be either `Graph`, `GraphPowerShell`, `GraphPowerShellWorkflow`, `PowerShellWorkflow`, `PowerShell`, `PowerShell72`, `Python3`, `Python2` or `Script`. Changing this forces a new resource to be created.

* `log_progress` - (Required) Progress log option.

* `log_verbose` - (Required) Verbose log option.

* `publish_content_link` - (Optional) One `publish_content_link` block as defined below.

* `description` - (Optional) A description for this credential.

* `content` - (Optional) The desired content of the runbook.

~> **Note:** The Azure API requires a `publish_content_link` to be supplied even when specifying your own `content`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `log_activity_trace_level` - (Optional) Specifies the activity-level tracing options of the runbook, available only for Graphical runbooks. Possible values are `0` for None, `9` for Basic, and `15` for Detailed. Must turn on Verbose logging in order to see the tracing.

* `draft` - (Optional) A `draft` block as defined below.

* `job_schedule` - (Optional) One or more `job_schedule` block as defined below.

~> **Note:** AzureRM provides a stand-alone [azurerm_automation_job_schedule](automation_job_schedule.html.markdown) and this inlined `job_schedule` property to manage the job schedules. At this time you should choose one of them to manage the job schedule resources.

---

The `publish_content_link` block supports the following:

* `uri` - (Required) The URI of the runbook content.

* `version` - (Optional) Specifies the version of the content

* `hash` - (Optional) A `hash` block as defined below.

---

The `hash` block supports:

* `algorithm` - (Required) Specifies the hash algorithm used to hash the content.

* `value` - (Required) Specifies the expected hash value of the content.

---

The `draft` block supports:

* `edit_mode_enabled` - (Optional) Whether the draft in edit mode.

* `content_link` - (Optional) A `publish_content_link` block as defined above.

* `output_types` - (Optional) Specifies the output types of the runbook.

* `parameters` - (Optional) A list of `parameters` block as defined below.

---

The `parameters` block supports:

* `key` - (Required) The name of the parameter.

* `type` - (Required) Specifies the type of this parameter.

* `mandatory` - (Optional) Whether this parameter is mandatory.

* `position` - (Optional) Specifies the position of the parameter.

* `default_value` - (Optional) Specifies the default value of the parameter.

---

The `job_schedule` block supports:

* `schedule_name` - (Required) The name of the Schedule.

* `parameters` - (Optional) A map of key/value pairs corresponding to the arguments that can be passed to the Runbook.

-> **Note:** The parameter keys/names must strictly be in lowercase, even if this is not the case in the runbook. This is due to a limitation in Azure Automation where the parameter names are normalized. The values specified don't have this limitation.

* `run_on` - (Optional) Name of a Hybrid Worker Group the Runbook will be executed on.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Automation Runbook ID.

* `job_schedule` - One or more `job_schedule` block as defined below.

---

An `job_schedule` block exports the following:

* `job_schedule_id` - The UUID of automation runbook job schedule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Runbook.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Runbook.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Runbook.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Runbook.

## Import

Automation Runbooks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_runbook.Get-AzureVMTutorial /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/runbooks/Get-AzureVMTutorial
```
