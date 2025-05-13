---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_scaling_plan_host_pool_association"
description: |-
  Manages a Virtual Desktop Scaling Plan Host Pool Association.
---

# azurerm_virtual_desktop_scaling_plan_host_pool_association

Manages a Virtual Desktop Scaling Plan Host Pool Association.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}
provider "azuread" {}

resource "azurerm_resource_group" "example" {
  name     = "rg-example-virtualdesktop"
  location = "West Europe"
}


data "azuread_service_principal" "example" {
  # In some environments this will be "Azure Virtual Desktop"
  display_name = "Windows Virtual Desktop"
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = "Desktop Virtualization Power On Off Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
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
  friendly_name       = "Scaling Plan Test"
  description         = "Test Scaling Plan"
  time_zone           = "GMT Standard Time"
  schedule {
    name                                 = "Weekdays"
    days_of_week                         = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    ramp_up_start_time                   = "06:00"
    ramp_up_load_balancing_algorithm     = "BreadthFirst"
    ramp_up_minimum_hosts_percent        = 20
    ramp_up_capacity_threshold_percent   = 10
    peak_start_time                      = "09:00"
    peak_load_balancing_algorithm        = "BreadthFirst"
    ramp_down_start_time                 = "18:00"
    ramp_down_load_balancing_algorithm   = "BreadthFirst"
    ramp_down_minimum_hosts_percent      = 10
    ramp_down_force_logoff_users         = false
    ramp_down_wait_time_minutes          = 45
    ramp_down_notification_message       = "Please log of in the next 45 minutes..."
    ramp_down_capacity_threshold_percent = 5
    ramp_down_stop_hosts_when            = "ZeroSessions"
    off_peak_start_time                  = "22:00"
    off_peak_load_balancing_algorithm    = "BreadthFirst"
  }

  depends_on = [azurerm_role_assignment.example]
}


resource "azurerm_virtual_desktop_scaling_plan_host_pool_association" "example" {
  host_pool_id    = azurerm_virtual_desktop_host_pool.example.id
  scaling_plan_id = azurerm_virtual_desktop_scaling_plan.example.id
  enabled         = true
  depends_on      = [azurerm_role_assignment.example]
}


```

## Argument Reference

The following arguments are supported:

- `host_pool_id` - (Required) The resource ID for the Virtual Desktop Host Pool. Changing this forces a new resource to be created.

- `scaling_plan_id` - (Required) The resource ID for the Virtual Desktop Scaling Plan. Changing this forces a new resource to be created.

- `enabled` - (Required) Should the Scaling Plan be enabled on this Host Pool.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Virtual Desktop Scaling Plan Host Pool association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Virtual Desktop Scaling Plan Host Pool association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Scaling Plan Host Pool association.
* `update` - (Defaults to 1 hour) Used when updating the Virtual Desktop Scaling Plan Host Pool association.
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Desktop Scaling Plan Host Pool association.

## Import

Associations between Virtual Desktop Scaling Plans and Virtual Desktop Host Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_scaling_plan_host_pool_association.example "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/scalingPlans/plan1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/hostPools/myhostpool"
```

-> **Note:** This ID is specific to Terraform - and is of the format `{virtualDesktopScalingPlanID}|{virtualDesktopHostPoolID}`.
