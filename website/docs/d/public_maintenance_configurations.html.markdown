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
  location    = "West Europe"
  scope       = "SQLManagedInstance"
  recur_every = "Monday-Thursday"
}

output "name" {
  value = data.azurerm_public_maintenance_configurations.existing.configs[0].name
}
```

## Argument Reference

* `location` - The Azure location to filter the list of Public Maintenance Configurations against.

* `scope` - The scope to filter the list of Public Maintenance Configurations against. Possible values are `Extension`, `Host`, `InGuestPatch`, `OSImage`, `SQLDB` and `SQLManagedInstance`.

* `recur_every` - The recurring window to filter the list of Public Maintenance Configurations against. Possible values are `Monday-Thursday` and `Friday-Sunday`

## Attributes Reference

* `configs` - A `configs` block as defined below.

---

A `configs` block exports the following:

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Public Maintenance Configuration.
