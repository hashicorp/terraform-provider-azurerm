---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maintenance_configuration"
description: |-
  Manages a Maintenance Configuration.
---

# azurerm_maintenance_configuration

Manages a maintenance configuration.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_maintenance_configuration" "example" {
  name                = "example-mc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  scope               = "All"

  tags = {
    Env = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Maintenance Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Maintenance Configuration should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specified the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `scope` - (Optional) The scope of the Maintenance Configuration. Possible values are `All`, `Extension`, `Host`, `InGuestPatch`, `OSImage`, `SQLDB` or `SQLManagedInstance`. Defaults to `All`.

* `visibility` - (Optional) The visibility of the Maintenance Configuration. The only allowable value is `Custom`.

* `window` - (Optional) A `window` block as defined below.

* `properties` - (Optional) A mapping of properties to assign to the resource.

* `tags` - (Optional) A mapping of tags to assign to the resource. The key could not contain upper case letter.

---

A `window` block supports:

* `start_date_time` - (Required) Effective start date of the maintenance window in YYYY-MM-DD hh:mm format.

* `expiration_date_time` - (Optional) Effective expiration date of the maintenance window in YYYY-MM-DD hh:mm format.

* `duration` - (Optional) The duration of the maintenance window in HH:mm format.

* `time_zone` - (Required) The time zone for the maintenance window. A list of timezones can be obtained by executing [System.TimeZoneInfo]::GetSystemTimeZones() in PowerShell.

* `recur_every` - (Optional) The rate at which a maintenance window is expected to recur. The rate can be expressed as daily, weekly, or monthly schedules.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Maintenance Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Maintenance Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Maintenance Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Maintenance Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Maintenance Configuration.

## Import

Maintenance Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maintenance_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.maintenance/maintenanceconfigurations/example-mc
```
