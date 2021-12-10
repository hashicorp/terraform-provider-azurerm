---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maintenance_configuration"
description: |-
  Get information about an existing Maintenance Configuration.
---

# Data Source: azurerm_maintenance_configuration

Use this data source to access information about an existing Maintenance Configuration.

## Example Usage

```hcl
data "azurerm_maintenance_configuration" "existing" {
  name                = "example-mc"
  resource_group_name = "example-resources"
}

output "id" {
  value = azurerm_maintenance_configuration.existing.id
}
```

## Argument Reference

* `name` - Specifies the name of the Maintenance Configuration.

* `resource_group_name` - Specifies the name of the Resource Group where this Maintenance Configuration exists.

## Attributes Reference

* `location` - The Azure location where the resource exists.

* `scope` - The scope of the Maintenance Configuration.

* `visibility` - The visibility of the Maintenance Configuration.

* `window` - A `window` block as defined below.

* `properties` - The properties assigned to the resource.

* `tags` - A mapping of tags assigned to the resource.

---

A `window` block exports the following:

* `start_date_time` - Effective start date of the maintenance window.

* `expiration_date_time` - Effective expiration date of the maintenance window.

* `duration` - The duration of the maintenance window.

* `time_zone` - The time zone for the maintenance window.

* `recur_every` The rate at which a maintenance window is expected to recur.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Maintenance Configuration.
