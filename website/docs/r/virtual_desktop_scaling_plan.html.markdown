---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_scaling_plan"
description: |-
  Manages a Virtual Desktop Scaling Plan .
---

# azurerm_virtual_desktop_scaling_plan

Manages a Virtual Desktop Scaling Plan.

## Disclaimers

-> **Note:** Scaling Plans are currently in preview and are only supported in a limited number of regions. Both the Scaling Plan and any referenced Host Pools must be deployed in a supported region. [Autoscale (preview) for Azure Virtual Desktop host pools](https://docs.microsoft.com/azure/virtual-desktop/autoscale-scaling-plan).

-> **Note:** Scaling Plans require specific permissions to be granted to the Windows Virtual Desktop application before a 'host_pool' can be configured. [Required Permissions for Scaling Plans](https://docs.microsoft.com/azure/virtual-desktop/autoscale-scaling-plan#create-a-custom-rbac-role).

## Example Usage

```hcl

resource "random_uuid" "example" {
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_role_definition" "example" {
  name        = "AVD-AutoScale"
  scope       = azurerm_resource_group.example.id
  description = "AVD AutoScale Role"
  permissions {
    actions = [
      "Microsoft.Insights/eventtypes/values/read",
      "Microsoft.Compute/virtualMachines/deallocate/action",
      "Microsoft.Compute/virtualMachines/restart/action",
      "Microsoft.Compute/virtualMachines/powerOff/action",
      "Microsoft.Compute/virtualMachines/start/action",
      "Microsoft.Compute/virtualMachines/read",
      "Microsoft.DesktopVirtualization/hostpools/read",
      "Microsoft.DesktopVirtualization/hostpools/write",
      "Microsoft.DesktopVirtualization/hostpools/sessionhosts/read",
      "Microsoft.DesktopVirtualization/hostpools/sessionhosts/write",
      "Microsoft.DesktopVirtualization/hostpools/sessionhosts/usersessions/delete",
      "Microsoft.DesktopVirtualization/hostpools/sessionhosts/usersessions/read",
      "Microsoft.DesktopVirtualization/hostpools/sessionhosts/usersessions/sendMessage/action",
      "Microsoft.DesktopVirtualization/hostpools/sessionhosts/usersessions/read"
    ]
    not_actions = []
  }
  assignable_scopes = [
    azurerm_resource_group.example.id,
  ]
}

data "azuread_service_principal" "example" {
  display_name = "Azure Virtual Desktop"
}

resource "azurerm_role_assignment" "example" {
  name                             = random_uuid.example.result
  scope                            = azurerm_resource_group.example.id
  role_definition_id               = azurerm_role_definition.example.role_definition_resource_id
  principal_id                     = data.azuread_service_principal.example.id
  skip_service_principal_aad_check = true
}

resource "azurerm_virtual_desktop_host_pool" "example" {
  name                 = "example-hostpool"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"
}

resource "azurerm_virtual_desktop_scaling_plan" "example" {
  name                = "example-scaling-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  friendly_name       = "Scaling Plan Example"
  description         = "Example Scaling Plan"
  time_zone           = "GMT Standard Time"
  schedule {
    name                                 = "Weekdays"
    days_of_week                         = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    ramp_up_start_time                   = "05:00"
    ramp_up_load_balancing_algorithm     = "BreadthFirst"
    ramp_up_minimum_hosts_percent        = 20
    ramp_up_capacity_threshold_percent   = 10
    peak_start_time                      = "09:00"
    peak_load_balancing_algorithm        = "BreadthFirst"
    ramp_down_start_time                 = "19:00"
    ramp_down_load_balancing_algorithm   = "DepthFirst"
    ramp_down_minimum_hosts_percent      = 10
    ramp_down_force_logoff_users         = false
    ramp_down_wait_time_minutes          = 45
    ramp_down_notification_message       = "Please log off in the next 45 minutes..."
    ramp_down_capacity_threshold_percent = 5
    ramp_down_stop_hosts_when            = "ZeroSessions"
    off_peak_start_time                  = "22:00"
    off_peak_load_balancing_algorithm    = "DepthFirst"
  }
  host_pool {
    hostpool_id          = azurerm_virtual_desktop_host_pool.example.id
    scaling_plan_enabled = true
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Virtual Desktop Scaling Plan should exist. Changing this forces a new Virtual Desktop Scaling Plan to be created.

* `name` - (Required) The name which should be used for this Virtual Desktop Scaling Plan . Changing this forces a new Virtual Desktop Scaling Plan to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Desktop Scaling Plan should exist. Changing this forces a new Virtual Desktop Scaling Plan to be created.

* `schedule` - (Required) One or more `schedule` blocks as defined below.

* `host_pool` - (Optional) One or more `host_pool` blocks as defined below.

* `time_zone` - (Required) Specifies the Time Zone which should be used by the Scaling Plan for time based events, [the possible values are defined here](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/).

---

* `description` - (Optional) A description of the Scaling Plan.

* `exclusion_tag` - (Optional) The name of the tag associated with the VMs you want to exclude from autoscaling.

* `friendly_name` - (Optional) Friendly name of the Scaling Plan.

* `host_pool` - (Optional) One or more `host_pool` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Virtual Desktop Scaling Plan .

---

A `host_pool` block supports the following:

* `hostpool_id` - (Required) The ID of the HostPool to assign the Scaling Plan to.

* `scaling_plan_enabled` - (Required) Specifies if the scaling plan is enabled or disabled for the HostPool.

---

A `schedule` block supports the following:

* `days_of_week` - (Required) A list of Days of the Week on which this schedule will be used. Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, and `Sunday`

* `name` - (Required) The name of the schedule.

* `off_peak_load_balancing_algorithm` - (Required) The load Balancing Algorithm to use during Off-Peak Hours. Possible values are `DepthFirst` and `BreadthFirst`.

* `off_peak_start_time` - (Required) The time at which Off-Peak scaling will begin. This is also the end-time for the Ramp-Down period. The time must be specified in "HH:MM" format.

* `peak_load_balancing_algorithm` - (Required) The load Balancing Algorithm to use during Peak Hours. Possible values are `DepthFirst` and `BreadthFirst`.

* `peak_start_time` - (Required) The time at which Peak scaling will begin. This is also the end-time for the Ramp-Up period. The time must be specified in "HH:MM" format.

* `ramp_down_capacity_threshold_percent` - (Required) This is the value in percentage of used host pool capacity that will be considered to evaluate whether to turn on/off virtual machines during the ramp-down and off-peak hours. For example, if capacity threshold is specified as 60% and your total host pool capacity is 100 sessions, autoscale will turn on additional session hosts once the host pool exceeds a load of 60 sessions.

* `ramp_down_force_logoff_users` - (Required) Whether users will be forced to log-off session hosts once the `ramp_down_wait_time_minutes` value has been exceeded during the Ramp-Down period. Possible values are `true` and `false`.

* `ramp_down_load_balancing_algorithm` - (Required) The load Balancing Algorithm to use during the Ramp-Down period. Possible values are `DepthFirst` and `BreadthFirst`.

* `ramp_down_minimum_hosts_percent` - (Required) The minimum percentage of session host virtual machines that you would like to get to for ramp-down and off-peak hours. For example, if Minimum percentage of hosts is specified as 10% and total number of session hosts in your host pool is 10, autoscale will ensure a minimum of 1 session host is available to take user connections.

* `ramp_down_notification_message` - (Required) The notification message to send to users during Ramp-Down period when they are required to log-off.

* `ramp_down_start_time` - (Required) The time at which Ramp-Down scaling will begin. This is also the end-time for the Ramp-Up period. The time must be specified in "HH:MM" format.

* `ramp_down_stop_hosts_when` - (Required) Controls Session Host shutdown behaviour during Ramp-Down period. Session Hosts can either be shutdown when all sessions on the Session Host have ended, or when there are no Active sessions left on the Session Host. Possible values are `ZeroSessions` and `ZeroActiveSessions`.

* `ramp_down_wait_time_minutes` - (Required) The number of minutes during Ramp-Down period that autoscale will wait after setting the session host VMs to drain mode, notifying any currently signed in users to save their work before forcing the users to logoff. Once all user sessions on the session host VM have been logged off, Autoscale will shut down the VM.

* `ramp_up_load_balancing_algorithm` - (Required) The load Balancing Algorithm to use during the Ramp-Up period. Possible values are `DepthFirst` and `BreadthFirst`.

* `ramp_up_start_time` - (Required) The time at which Ramp-Up scaling will begin. This is also the end-time for the Ramp-Up period. The time must be specified in "HH:MM" format.

* `ramp_up_capacity_threshold_percent` - (Optional) This is the value of percentage of used host pool capacity that will be considered to evaluate whether to turn on/off virtual machines during the ramp-up and peak hours. For example, if capacity threshold is specified as `60%` and your total host pool capacity is `100` sessions, autoscale will turn on additional session hosts once the host pool exceeds a load of `60` sessions.

* `ramp_up_minimum_hosts_percent` - (Optional) Specifies the minimum percentage of session host virtual machines to start during ramp-up for peak hours. For example, if Minimum percentage of hosts is specified as `10%` and total number of session hosts in your host pool is `10`, autoscale will ensure a minimum of `1` session host is available to take user connections.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Desktop Scaling Plan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Virtual Desktop Scaling Plan.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Scaling Plan.
* `update` - (Defaults to 1 hour) Used when updating the Virtual Desktop Scaling Plan.
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Desktop Scaling Plan.

## Import

Virtual Desktop Scaling Plans can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_scaling_plan.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/scalingPlans/plan1
```
