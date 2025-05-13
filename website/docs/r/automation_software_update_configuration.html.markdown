---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_software_update_configuration"
description: |-
  Manages an Automation Software Update Configuration.
---

# azurerm_automation_software_update_configuration

Manages an Automation Software Update Configuration.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

resource "azurerm_automation_account" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runbook" "example" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a example runbook for terraform acceptance example"
  runbook_type = "Python3"

  content = <<CONTENT
# Some example content
# for Terraform acceptance example
CONTENT
  tags = {
    ENV = "runbook_test"
  }
}

resource "azurerm_automation_software_update_configuration" "example" {
  name                  = "example"
  automation_account_id = azurerm_automation_account.example.id

  linux {
    classifications_included = "Security"
    excluded_packages        = ["apt"]
    included_packages        = ["vim"]
    reboot                   = "IfRequired"
  }

  pre_task {
    source = azurerm_automation_runbook.example.name
    parameters = {
      COMPUTER_NAME = "Foo"
    }
  }

  duration = "PT2H2M2S"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Automation. Changing this forces a new Automation to be created.

* `automation_account_id` - (Required) The ID of Automation Account to manage this Source Control. Changing this forces a new Automation Source Control to be created.

---

* `duration` - (Optional) Maximum time allowed for the software update configuration run. using format `PT[n]H[n]M[n]S` as per ISO8601. Defaults to `PT2H`.

* `linux` - (Optional) A `linux` block as defined below.

* `windows` - (Optional) A `windows` block as defined below.

~> **Note:** One of `linux` or `windows` must be specified.

* `virtual_machine_ids` - (Optional) Specifies a list of Azure Resource IDs of azure virtual machines.

* `non_azure_computer_names` - (Optional) Specifies a list of names of non-Azure machines for the software update configuration.

* `target` - (Optional) A `target` blocks as defined below.

* `post_task` - (Optional) A `post_task` blocks as defined below.

* `pre_task` - (Optional) A `pre_task` blocks as defined below.

* `schedule` - (Required) A `schedule` blocks as defined below.

---

A `linux` block supports the following:

* `classifications_included` - (Required) Specifies the list of update classifications included in the Software Update Configuration. Possible values are `Unclassified`, `Critical`, `Security` and `Other`.

* `excluded_packages` - (Optional) Specifies a list of packages to excluded from the Software Update Configuration.

* `included_packages` - (Optional) Specifies a list of packages to included from the Software Update Configuration.

* `reboot` - (Optional) Specifies the reboot settings after software update, possible values are `IfRequired`, `Never`, `RebootOnly` and `Always`. Defaults to `IfRequired`.

---

A `windows` block supports the following:

* `classifications_included` - (Required) Specifies the list of update classification. Possible values are `Unclassified`, `Critical`, `Security`, `UpdateRollup`, `FeaturePack`, `ServicePack`, `Definition`, `Tools` and `Updates`.

* `excluded_knowledge_base_numbers` - (Optional) Specifies a list of knowledge base numbers excluded.

* `included_knowledge_base_numbers` - (Optional) Specifies a list of knowledge base numbers included.

* `reboot` - (Optional) Specifies the reboot settings after software update, possible values are `IfRequired`, `Never`, `RebootOnly` and `Always`. Defaults to `IfRequired`.

---

A `target` block supports the following:

* `azure_query` - (Optional) One or more `azure_query` blocks as defined above.

* `non_azure_query` - (Optional) One or more `non_azure_query` blocks as defined above.

---

A `azure_query` block supports the following:

* `locations` - (Optional) Specifies a list of locations to scope the query to.

* `scope` - (Optional) Specifies a list of Subscription or Resource Group ARM Ids to query.

* `tag_filter` - (Optional) Specifies how the specified tags to filter VMs. Possible values are `Any` and `All`.

* `tags` - (Optional) A mapping of tags used for query filter. One or more `tags` block as defined below.

---

A `tags` block supports the following:

* `tag` - (Required) Specifies the name of the tag to filter.

* `values` - (Required) Specifies a list of values for this tag key.

---

A `non_azure_query` block supports the following:

* `function_alias` - (Optional) Specifies the Log Analytics save search name.

* `workspace_id` - (Optional) The workspace id for Log Analytics in which the saved search in.

---

A `pre_task` block supports the following:

* `parameters` - (Optional) Specifies a map of parameters for the task.

* `source` - (Optional) The name of the runbook for the pre task.

---

A `post_task` block supports the following:

* `parameters` - (Optional) Specifies a map of parameters for the task.

* `source` - (Optional) The name of the runbook for the post task.

---

A `schedule` block supports the following:

* `frequency` - (Required) The frequency of the schedule. - can be either `OneTime`, `Day`, `Hour`, `Week`, or `Month`.

* `is_enabled` - (Optional) Whether the schedule is enabled. Defaults to `true`.

* `description` - (Optional) A description for this Schedule.

* `interval` - (Optional) The number of `frequency`s between runs. Only valid when frequency is `Day`, `Hour`, `Week`, or `Month`.

* `start_time` - (Optional) Start time of the schedule. Must be at least five minutes in the future. Defaults to seven minutes in the future from the time the resource is created.

* `expiry_time` - (Optional) The end time of the schedule.

* `time_zone` - (Optional) The timezone of the start time. Defaults to `Etc/UTC`. For possible values see: <https://docs.microsoft.com/en-us/rest/api/maps/timezone/gettimezoneenumwindows>

* `advanced_week_days` - (Optional) List of days of the week that the job should execute on. Only valid when frequency is `Week`. Possible values include `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, and `Sunday`.

* `advanced_month_days` - (Optional) List of days of the month that the job should execute on. Must be between `1` and `31`. `-1` for last day of the month. Only valid when frequency is `Month`.

* `monthly_occurrence` - (Optional) List of `monthly_occurrence` blocks as defined below to specifies occurrences of days within a month. Only valid when frequency is `Month`. The `monthly_occurrence` block supports fields as defined below.

---

The `monthly_occurrence` block supports the following:

* `day` - (Required) Day of the occurrence. Must be one of `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`.

* `occurrence` - (Required) Occurrence of the week within the month. Must be between `1` and `4`. `-1` for last week within the month.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Software Update Configuration.

* `error_code` - The Error code when failed.

* `error_message` - The Error message indicating why the operation failed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
* `update` - (Defaults to 30 minutes) Used when updating the Automation.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation.

## Import

Automations Software Update Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_software_update_configuration.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/softwareUpdateConfigurations/suc1
```
