---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_host_pool"
description: |-
  Manages a Virtual Desktop Host Pool.
---

# azurerm_virtual_desktop_host_pool

Manages a Virtual Desktop Host Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_desktop_host_pool" "example" {
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  name                     = "pooleddepthfirst"
  friendly_name            = "pooleddepthfirst"
  validate_environment     = true
  start_vm_on_connect      = true
  custom_rdp_properties    = "audiocapturemode:i:1;audiomode:i:0;"
  description              = "Acceptance Test: A pooled host pool - pooleddepthfirst"
  type                     = "Pooled"
  maximum_sessions_allowed = 50
  load_balancer_type       = "DepthFirst"
  scheduled_agent_updates {
    enabled = true
    schedule {
      day_of_week = "Saturday"
      hour_of_day = 2
    }
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Desktop Host Pool. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Virtual Desktop Host Pool. Changing this forces a new resource to be created.

* `location` - (Required) The location/region where the Virtual Desktop Host Pool is located. Changing this forces a new resource to be created.

* `type` - (Required) The type of the Virtual Desktop Host Pool. Valid options are `Personal` or `Pooled`. Changing the type forces a new resource to be created.

* `load_balancer_type` - (Required) `BreadthFirst` load balancing distributes new user sessions across all available session hosts in the host pool. Possible values are `BreadthFirst`, `DepthFirst` and `Persistent`.
    `DepthFirst` load balancing distributes new user sessions to an available session host with the highest number of connections but has not reached its maximum session limit threshold.
    `Persistent` should be used if the host pool type is `Personal`

* `friendly_name` - (Optional) A friendly name for the Virtual Desktop Host Pool.

* `description` - (Optional) A description for the Virtual Desktop Host Pool.

* `validate_environment` - (Optional) Allows you to test service changes before they are deployed to production. Defaults to `false`.

* `start_vm_on_connect` - (Optional) Enables or disables the Start VM on Connection Feature. Defaults to `false`.

* `custom_rdp_properties` - (Optional) A valid custom RDP properties string for the Virtual Desktop Host Pool, available properties can be [found in this article](https://docs.microsoft.com/windows-server/remote/remote-desktop-services/clients/rdp-files).

* `personal_desktop_assignment_type` - (Optional) `Automatic` assignment – The service will select an available host and assign it to an user. Possible values are `Automatic` and `Direct`. `Direct` Assignment – Admin selects a specific host to assign to an user. Changing this forces a new resource to be created.

~> **Note:** `personal_desktop_assignment_type` is required if the `type` of your Virtual Desktop Host Pool is `Personal`

* `public_network_access` - (Optional) Whether public network access is allowed for the Virtual Desktop Host Pool. Possible values are `Enabled`, `Disabled`, `EnabledForClientsOnly` and `EnabledForSessionHostsOnly`. Defaults to `Enabled`.

* `maximum_sessions_allowed` - (Optional) A valid integer value from 0 to 999999 for the maximum number of users that have concurrent sessions on a session host.
    Should only be set if the `type` of your Virtual Desktop Host Pool is `Pooled`.

* `preferred_app_group_type` - (Optional) Option to specify the preferred Application Group type for the Virtual Desktop Host Pool. Valid options are `None`, `Desktop` or `RailApplications`. Default is `Desktop`.

* `scheduled_agent_updates` - (Optional) A `scheduled_agent_updates` block as defined below. This enables control of when Agent Updates will be applied to Session Hosts.

* `vm_template` - (Optional) A VM template for session hosts configuration within hostpool. This is a JSON string.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `scheduled_agent_updates` block supports the following:

* `enabled` - (Optional) Enables or disables scheduled updates of the AVD agent components (RDAgent, Geneva Monitoring agent, and side-by-side stack) on session hosts. If this is enabled then up to two `schedule` blocks must be defined. Default is `false`.

~> **Note:** if `enabled` is set to `true` then at least one and a maximum of two `schedule` blocks must be provided.

* `timezone` - (Optional) Specifies the time zone in which the agent update schedule will apply, [the possible values are defined here](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/). If `use_session_host_timezone` is enabled then it will override this setting. Default is `UTC`
* `use_session_host_timezone` - (Optional) Specifies whether scheduled agent updates should be applied based on the timezone of the affected session host. If configured then this setting overrides `timezone`. Default is `false`.
* `schedule` - (Optional) A `schedule` block as defined below. A maximum of two blocks can be added.

---

A `schedule` block supports the following:

* `day_of_week` - (Required) The day of the week on which agent updates should be performed. Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, and `Sunday`
* `hour_of_day` - (Required) The hour of day the update window should start. The update is a 2 hour period following the hour provided. The value should be provided as a number between 0 and 23, with 0 being midnight and 23 being 11pm. A leading zero should not be used.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Desktop Host Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Virtual Desktop Host Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Host Pool.
* `update` - (Defaults to 1 hour) Used when updating the Virtual Desktop Host Pool.
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Desktop Host Pool.

## Import

Virtual Desktop Host Pools can be imported using the `resource id`, e.g.

```text
terraform import azurerm_virtual_desktop_host_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/hostPools/myhostpool
```
