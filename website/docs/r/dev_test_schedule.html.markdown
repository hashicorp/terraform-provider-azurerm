---
subcategory: "Dev Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_schedule"
description: |-
    Manages automated startup and shutdown schedules for Azure Dev Test Lab.
---

# azurerm_dev_test_schedule

Manages automated startup and shutdown schedules for Azure Dev Test Lab.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_dev_test_lab" "example" {
  name                = "YourDevTestLab"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_dev_test_schedule" "example" {
  name                = "LabVmAutoStart"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  lab_name            = azurerm_dev_test_lab.example.name

  weekly_recurrence {
    time      = "1100"
    week_days = ["Monday", "Tuesday"]
  }

  time_zone_id = "Pacific Standard Time"
  task_type    = "LabVmsStartupTask"

  notification_settings {
  }

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the dev test lab schedule. Valid value for name depends on the `task_type`. For instance for task_type `LabVmsStartupTask` the name needs to be `LabVmAutoStart`.

* `location` - (Required) The location where the schedule is created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the dev test lab schedule. Changing this forces a new resource to be created.

* `lab_name` - (Required) The name of the dev test lab. Changing this forces a new resource to be created.

* `status` - The status of this schedule. Possible values are `Enabled` and `Disabled`. Defaults to `Disabled`.

* `task_type` - (Required) The task type of the schedule. Possible values include `LabVmsShutdownTask` and `LabVmAutoStart`.

* `time_zone_id` - (Required) The time zone ID (e.g. Pacific Standard time).

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `weekly_recurrence` - block supports the following:

* `time` - The time when the schedule takes effect.

* `week_days` -  A list of days that this schedule takes effect . Possible values include `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`.

---

A `daily_recurrence` - block supports the following:

* `time` - The time each day when the schedule takes effect.

---

A `notification_settings` - (Required)  - block supports the following:

* `status` - The status of the notification. Possible values are `Enabled` and `Disabled`. Defaults to `Disabled`

* `time_in_minutes` - Time in minutes before event at which notification will be sent.

* `webhook_url` - The webhook URL to which the notification will be sent.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the DevTest Schedule.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DevTest Schedule.
* `update` - (Defaults to 30 minutes) Used when updating the DevTest Schedule.
* `read` - (Defaults to 5 minutes) Used when retrieving the DevTest Schedule.
* `delete` - (Defaults to 30 minutes) Used when deleting the DevTest Schedule.

## Import

DevTest Schedule's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dev_test_schedule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DevTestLab/labs/myDevTestLab/schedules/labvmautostart
```
