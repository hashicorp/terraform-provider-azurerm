---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_schedule"
sidebar_current: "docs-azurerm-resource-automation-schedule"
description: |-
  Creates a new Automation Schedule.
---

# azurerm\_automation\_schedule

Creates a new Automation Schedule.

## Example Usage

```
resource "azurerm_resource_group" "example" {
 name = "resourceGroup1"
 location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "account1"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku {
        name = "Free"
  }
}

resource "azurerm_automation_schedule" "example" {
  name                = "schedule1"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_automation_account.example.name}"
  frequency           = "OneTime"
  first_run { 
        "hour" = 20
        "minute" = 5
        "second" = 0
  }
  description         = "This is an example schedule"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Schedule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Schedule is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the automation account in which the Schedule is created. Changing this forces a new resource to be created.

* `description` -  (Optional) A description for this Schedule.

* `start_time` -  (Optional) Start time of the schedule. Must be at least five minutes in the future.

* `expiry_time` -  (Optional) The end time of the schedule.

* `interval` -  (Optional) The interval of the schedule. Must be set if the schedule is recurring. NOT YET SUPPORTED due to lack of SDK support.
 
* `frequency` - (Required) The frequency of the schedule. - can be either `OneTime`, `Day`, `Hour`, `Week`, or `Month`.

* `first_run` - (Optional) If an exact start time is not suitable, it can be used to make constraints for the first run. The start time will be calculated depending on these constraints. `start_time` will override this settings if defined.

`first_run` supports the following:

* `second` - (Optional) In which second should the schedule first triggered.
 
* `minute` - (Optional) In which minute should the schedule first triggered.

* `hour` - (Optional) In which hour should the schedule first triggered. Ignored if the frequency is `Hour`.

* `day_of_week` - (Optional) On which day of the week should the schedule first triggered. Ignored if the frequency is `OneTime`, `Hour` or `Day`. (0 - Sunday)

* `day_of_month` - (Optional) On which day of the month should the schedule first triggered. Ignored if the frequency is `Hour`, `Day` or `Week`.
 
## Attributes Reference

The following attributes are exported:

* `id` - The Automation Schedule ID.

## Import

Automation Schedule can be imported using the `resource id`, e.g.

```
terraform import azurerm_automation_schedule.schedule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/schedules/schedule1
```
