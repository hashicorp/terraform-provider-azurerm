---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_schedule"
sidebar_current: "docs-azurerm-resource-automation-schedule"
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
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku {
    name = "Basic"
  }
}

resource "azurerm_automation_schedule" "example" {
  name                    = "tfex-automation-schedule"
  resource_group_name     = "${azurerm_resource_group.example.name}"
  automation_account_name = "${azurerm_automation_account.example.name}"
  frequency               = "Hour"
  interval                = 1
  timezone                = "Central Europe Standard Time"
  start_time              = "2014-04-15T18:00:15+02:00"
  description             = "This is an example schedule"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Schedule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Schedule is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Schedule is created. Changing this forces a new resource to be created.

* `frequency` - (Required) The frequency of the schedule. - can be either `OneTime`, `Day`, `Hour`, `Week`, or `Month`.

* `description` -  (Optional) A description for this Schedule.

* `interval` -  (Optional) The number of `frequency`s between runs. Only valid for `Day`, `Hour`, `Week`, or `Month` and defaults to `1`.

* `start_time` -  (Optional) Start time of the schedule. Must be at least five minutes in the future. Defaults to seven minutes in the future from the time the resource is created.

* `expiry_time` -  (Optional) The end time of the schedule.

* `timezone` - (Optional) The timezone of the start time. Defaults to `UTC`. For possible values see: https://msdn.microsoft.com/en-us/library/ms912391(v=winembedded.11).aspx

## Attributes Reference

The following attributes are exported:

* `id` - The Automation Schedule ID.

## Import

Automation Schedule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_schedule.schedule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/schedules/schedule1
```
