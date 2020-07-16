---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_schedule"
description: |-
  Manages a Automation Schedule.
---

# azurerm_automation_schedule

Manages a Automation Schedule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-automation-account"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "tfex-automation-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name = "Basic"
  }
}

resource "azurerm_automation_schedule" "example" {
  name                    = "tfex-automation-schedule"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  frequency               = "Week"
  interval                = 1
  timezone                = "Australia/Perth"
  start_time              = "2014-04-15T18:00:15+02:00"
  description             = "This is an example schedule"
  week_days               = ["Friday"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Schedule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Schedule is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Schedule is created. Changing this forces a new resource to be created.

* `frequency` - (Required) The frequency of the schedule. - can be either `OneTime`, `Day`, `Hour`, `Week`, or `Month`.

* `description` -  (Optional) A description for this Schedule.

* `interval` -  (Optional) The number of `frequency`s between runs. Only valid when frequency is `Day`, `Hour`, `Week`, or `Month` and defaults to `1`.

* `start_time` -  (Optional) Start time of the schedule. Must be at least five minutes in the future. Defaults to seven minutes in the future from the time the resource is created.

* `expiry_time` -  (Optional) The end time of the schedule.

* `timezone` - (Optional) The timezone of the start time. Defaults to `UTC`. For possible values see: https://s2.automation.ext.azure.com/api/Orchestrator/TimeZones?_=1594792230258

* `week_days` - (Optional) List of days of the week that the job should execute on. Only valid when frequency is `Week`.

* `month_days` - (Optional) List of days of the month that the job should execute on. Must be between `1` and `31`. `-1` for last day of the month. Only valid when frequency is `Month`.

* `monthly_occurrence` - (Optional) List of occurrences of days within a month. Only valid when frequency is `Month`. The `monthly_occurrence` block supports fields documented below.

---

The `monthly_occurrence` block supports:

* `day` - (Required) Day of the occurrence. Must be one of `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`.

* `occurrence` - (Required) Occurrence of the week within the month. Must be between `1` and `5`. `-1` for last week within the month.

## Attributes Reference

The following attributes are exported:

* `id` - The Automation Schedule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Schedule.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Schedule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Schedule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Schedule.

## Import

Automation Schedule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_schedule.schedule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/schedules/schedule1
```
