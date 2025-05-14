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

* `hourly_schedule` - (Optional) Sets an hourly snapshot schedule. A `hourly_schedule` block as defined below.

* `daily_schedule` - (Optional) Sets a daily snapshot schedule. A `daily_schedule` block as defined below.
  
* `weekly_schedule` - (Optional) Sets a weekly snapshot schedule. A `weekly_schedule` block as defined below.

* `monthly_schedule` - (Optional) Sets a monthly snapshot schedule. A `monthly_schedule` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

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

A `monthly_schedule` block supports the following:

* `snapshots_to_keep` - (Required) How many hourly snapshots to keep, valid range is from 0 to 255.

* `days_of_month` - (Required) List of the days of the month when the snapshots will be created, valid range is from 1 to 30.

* `hour` - (Required) Hour of the day that the snapshots will be created, valid range is from 0 to 23.

* `minute` - (Required) Minute of the hour that the snapshots will be created, valid range is from 0 to 59.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NetApp Snapshot.
  
* `name` - (Required) The name of the NetApp Snapshot Policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Snapshot Policy should be created.
  
* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp Account in which the NetApp Snapshot Policy was created. Changing this forces a new resource to be created.

* `enabled` - (Required) Defines that the NetApp Snapshot Policy is enabled or not.

* `hourly_schedule` - Hourly snapshot schedule.

* `daily_schedule` - Daily snapshot schedule.
  
* `weekly_schedule` - Weekly snapshot schedule.

* `monthly_schedule` - Monthly snapshot schedule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the NetApp Snapshot Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Snapshot Policy.
* `update` - (Defaults to 30 minutes) Used when updating the NetApp Snapshot Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the NetApp Snapshot Policy.

## Import

NetApp Snapshot Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_snapshot_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/snapshotPolicies/snapshotpolicy1
```
