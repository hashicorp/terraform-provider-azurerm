---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_maintenance_configurations"
description: |-
  Get information about existing Public Maintenance Configurations.
---

# Data Source: azurerm_public_maintenance_configurations

Use this data source to access information about existing Public Maintenance Configurations.

## Example Usage

```hcl
data "azurerm_public_maintenance_configurations" "existing" {
  location_filter    = "westeurope"
  scope_filter       = "SQLManagedInstance"
  recur_every_filter = "weekMondayToThursday"
}

output "name" {
  value = data.azurerm_public_maintenance_configurations.existing.public_maintenance_configurations[0].name
}
```

## Argument Reference

* `location_filter` - The Azure location to filter the list of Public Maintenance Configurations against.

* `scope_filter` - The scope to filter the list of Public Maintenance Configurations against. Possible values are `All`, `Extension`, `Host`, `InGuestPatch`, `OSImage`, `SQLDB` and `SQLManagedInstance`.

* `recur_every_filter` - The recurring window to filter the list of Public Maintenance Configurations against. Possible values are `weekMondayToThursday` and `weekFridayToSunday`

## Attributes Reference

* `public_maintenance_configurations` - A `public_maintenance_configurations` block as defined below.

---

A `public_maintenance_configurations` block exports the following:

* `name` - The name of the Public Maintenance Configuration.

* `id` - The id of the Public Maintenance Configuration.

* `location` - The Azure location of the Public Maintenance Configuration.

* `description` - A description of the Public Maintenance Configuration.

* `duration` - The duration of the Public Maintenance Configuration window.

* `maintenance_scope` - The scope of the Public Maintenance Configuration.

* `time_zone` - The time zone for the maintenance window.

* `recur_every` - The rate at which a maintenance window is expected to recur.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Public Maintenance Configuration.
