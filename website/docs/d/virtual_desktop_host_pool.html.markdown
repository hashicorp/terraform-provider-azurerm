---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_host_pool"
description: |-
  Gets information about an existing Virtual Desktop Host Pool.
---

# Data Source: azurerm_virtual_desktop_host_pool

Use this data source to access information about an existing Virtual Desktop Host Pool.

## Example Usage

```hcl

data "azurerm_virtual_desktop_host_pool" "example" {
  name                = "example-pool"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Desktop Host Pool to retrieve.

* `resource_group_name` - (Required) The name of the resource group where the Virtual Desktop Host Pool exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Desktop Host Pool.

* `location` - The location/region where the Virtual Desktop Host Pool is located.

* `type` - The type of the Virtual Desktop Host Pool.

* `load_balancer_type` - The type of load balancing performed by the Host Pool
    
* `friendly_name` - The friendly name for the Virtual Desktop Host Pool.

* `description` - The description for the Virtual Desktop Host Pool.

* `validate_environment` - Returns `true` if the Host Pool is in Validation mode.

* `start_vm_on_connect` - Returns `true` if the Start VM on Connection Feature is enabled.

* `custom_rdp_properties` - The custom RDP properties string for the Virtual Desktop Host Pool.

* `personal_desktop_assignment_type` - The type of personal desktop assignment in use by the Host Pool

* `maximum_sessions_allowed` - The maximum number of users that can have concurrent sessions on a session host.

* `preferred_app_group_type` - The preferred Application Group type for the Virtual Desktop Host Pool.

* `scheduled_agent_updates` - A `scheduled_agent_updates` block as defined below.

* `tags` - A mapping of tags to assign to the resource.

---

A `scheduled_agent_updates` block exports the following:

* `enabled` - Are scheduled updates of the AVD agent components (RDAgent, Geneva Monitoring agent, and side-by-side stack) enabled on session hosts.
* `timezone` - The time zone in which the agent update schedule will apply.
* `use_session_host_timezone` - Specifies whether scheduled agent updates should be applied based on the timezone of the affected session host.
* `schedule` - A `schedule` block as defined below.

---

A `schedule` block exports the following:

* `day_of_week` - The day of the week on which agent updates should be performed.
* `hour_of_day` - The hour of day the update window should start.

---


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Host Pool.
