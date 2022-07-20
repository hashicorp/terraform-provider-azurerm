---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_snapshot_policy"
description: |-
  Manages a NetApp Snapshot Policy.
---

# azurerm_netapp_snapshot_policy

Manages a NetApp Snapshot Policy.

## NetApp Snapshot Policy Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resource-group-01"
  location = "East US"
}

resource "azurerm_netapp_account" "example" {
  name                = "netappaccount-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_snapshot_policy" "example" {
  name                = "snapshotpolicy-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  enabled             = true

  hourly_schedule {
    snapshots_to_keep = 4
    minute            = 15
  }

  daily_schedule {
    snapshots_to_keep = 2
    hour              = 20
    minute            = 15
  }

  weekly_schedule {
    snapshots_to_keep = 1
    days_of_week      = ["Monday", "Friday"]
    hour              = 23
    minute            = 0
  }

  monthly_schedule {
    snapshots_to_keep = 1
    days_of_month     = [1, 15, 20, 30]
    hour              = 5
    minute            = 45
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Snapshot Policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Snapshot Policy should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp Account in which the NetApp Snapshot Policy should be created. Changing this forces a new resource to be created.

* `enabled` - (Required) Defines that the NetApp Snapshot Policy is enabled or not.

* `hourly_schedule` - (Optional) Sets an hourly snapshot schedule. See details in below `hourly_schedule` block.

* `daily_schedule` - (Optional) Sets a daily snapshot schedule. See details in below `daily_schedule` block.
  
* `weekly_schedule` - (Optional) Sets a weekly snapshot schedule. See details in below `weekly_schedule` block.

* `monthly_schedule` - (Optional) Sets a monthly snapshot schedule. See details in below `monthly_schedule` block.

---

An `hourly_schedule` block supports the following:

* `snapshots_to_keep` - (Required) How many hourly snapshots to keep, valid range is from 0 to 255.

* `minute` - (Required) Minute of the hour that the snapshots will be created, valid range is from 0 to 59.

---

A `daily_schedule` block supports the following:

* `snapshots_to_keep` - (Required) How many hourly snapshots to keep, valid range is from 0 to 255.

* `hour` - (Required) Hour of the day that the snapshots will be created, valid range is from 0 to 23.

* `minute` - (Required) Minute of the hour that the snapshots will be created, valid range is from 0 to 59.

---

A `weekly_schedule` block supports the following:

* `snapshots_to_keep` - (Required) How many hourly snapshots to keep, valid range is from 0 to 255.

* `days_of_week` - (Required) List of the week days using English names when the snapshots will be created.

* `hour` - (Required) Hour of the day that the snapshots will be created, valid range is from 0 to 23.

* `minute` - (Required) Minute of the hour that the snapshots will be created, valid range is from 0 to 59.

---

A `weekly_schedule` block supports the following:

* `snapshots_to_keep` - (Required) How many hourly snapshots to keep, valid range is from 0 to 255.

* `monthly_schedule` - (Required) List of the days of the month when the snapshots will be created, valid range is from 1 to 30.

* `hour` - (Required) Hour of the day that the snapshots will be created, valid range is from 0 to 23.

* `minute` - (Required) Minute of the hour that the snapshots will be created, valid range is from 0 to 59.

---

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the NetApp Snapshot Policy.
* `update` - (Defaults to 30 minutes) Used when updating the NetApp Snapshot Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Snapshot Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the NetApp Snapshot Policy.

## Import

NetApp Snapshot Policy can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_snapshot_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/snapshotPolicies/snapshotpolicy1
```
