---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_snapshot_policy"
description: |-
  Gets information about an existing NetApp Snapshot Policy
---

# Data Source: azurerm_netapp_snapshot_policy

Uses this data source to access information about an existing NetApp Snapshot Policy.

## NetApp Snapshot Policy Usage

```hcl
data "azurerm_netapp_snapshot_policy" "example" {
  resource_group_name = "acctestRG"
  account_name        = "acctestnetappaccount"
  name                = "example-snapshot-policy"
}

output "id" {
  value = data.azurerm_netapp_snapshot_policy.example.id
}

output "name" {
  value = data.azurerm_netapp_snapshot_policy.example.name
}

output "enabled" {
  value = data.azurerm_netapp_snapshot_policy.example.enabled
}

output "hourly_schedule" {
  value = data.azurerm_netapp_snapshot_policy.example.hourly_schedule
}

output "daily_schedule" {
  value = data.azurerm_netapp_snapshot_policy.example.daily_schedule
}

output "weekly_schedule" {
  value = data.azurerm_netapp_snapshot_policy.example.weekly_schedule
}

output "monthly_schedule" {
  value = data.azurerm_netapp_snapshot_policy.example.monthly_schedule
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the NetApp Snapshot Policy.

* `resource_group_name` - The Name of the Resource Group where the NetApp Snapshot Policy exists.

* `account_name` - The name of the NetApp account where the NetApp Snapshot Policy exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NetApp Snapshot.
  
* `name` - The name of the NetApp Snapshot Policy.

* `resource_group_name` - The name of the resource group where the NetApp Snapshot Policy should be created.
  
* `location` - Specifies the supported Azure location where the resource exists.

* `account_name` - The name of the NetApp Account in which the NetApp Snapshot Policy was created.

* `enabled` - Defines that the NetApp Snapshot Policy is enabled or not.

* `hourly_schedule` - Hourly snapshot schedule.

* `daily_schedule` - Daily snapshot schedule.
  
* `weekly_schedule` - Weekly snapshot schedule.

* `monthly_schedule` - Monthly snapshot schedule.

---

An `hourly_schedule` block exports the following:

* `snapshots_to_keep` - How many hourly snapshots to keep.

* `minute` - Minute of the hour that the snapshots will be created.

---

A `daily_schedule` block exports the following:

* `snapshots_to_keep` - How many hourly snapshots to keep.

* `hour` - Hour of the day that the snapshots will be created.

* `minute` - Minute of the hour that the snapshots will be created.

---

A `weekly_schedule` block supports the following:

* `snapshots_to_keep` - How many hourly snapshots to keep.

* `days_of_week` - List of the week days using English names when the snapshots will be created.

* `hour` - Hour of the day that the snapshots will be created.

* `minute` - Minute of the hour that the snapshots will be created.

---

A `weekly_schedule` block supports the following:

* `snapshots_to_keep` -  How many hourly snapshots to keep.

* `monthly_schedule` -  List of the days of the month when the snapshots will be created.

* `hour` -  Hour of the day that the snapshots will be created.

* `minute` -  Minute of the hour that the snapshots will be created.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Volume.
