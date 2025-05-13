---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_exadata_infrastructure"
description: |-
  Manages a Cloud Exadata Infrastructure.
---

# azurerm_oracle_exadata_infrastructure

Manages a Cloud Exadata Infrastructure.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_oracle_exadata_infrastructure" "example" {
  name                = "example-exadata-infra"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  zones               = ["1"]
  display_name        = "example-exadata-infra"
  storage_count       = 3
  compute_count       = 2
  shape               = "Exadata.X9M"
}
```

## Arguments Reference

The following arguments are supported:

* `compute_count` - (Required) The number of compute servers for the Cloud Exadata Infrastructure. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `display_name` - (Required) The user-friendly name for the Cloud Exadata Infrastructure resource. The name does not need to be unique. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `location` - (Required) The Azure Region where the Cloud Exadata Infrastructure should exist. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `name` - (Required) The name which should be used for this Cloud Exadata Infrastructure. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the ODB@A Infrastructure should exist. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `shape` - (Required) The shape of the ODB@A infrastructure resource. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `storage_count` - (Required) The number of storage servers for the Cloud Exadata Infrastructure. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `zones` - (Required) Cloud Exadata Infrastructure zones. Changing this forces a new Cloud Exadata Infrastructure to be created.

---

* `customer_contacts` - (Optional) The email address used by Oracle to send notifications regarding databases and infrastructure. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `maintenance_window` - (Optional) One or more `maintenance_window` blocks as defined below. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Cloud Exadata Infrastructure.

---

A `maintenance_window` block supports the following:

* `days_of_week` - (Optional) Days during the week when maintenance should be performed. Valid values are: `0` - represents time slot `0:00 - 3:59 UTC - 4` - represents time slot `4:00 - 7:59 UTC - 8` - represents time slot 8:00 - 11:59 UTC - 12 - represents time slot 12:00 - 15:59 UTC - 16 - represents time slot 16:00 - 19:59 UTC - 20 - represents time slot `20:00 - 23:59 UTC`. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `hours_of_day` - (Optional) The window of hours during the day when maintenance should be performed. The window is a 4 hour slot. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `lead_time_in_weeks` - (Optional) Lead time window allows user to set a lead time to prepare for a down time. The lead time is in weeks and valid value is between `1` to `4`. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `months` - (Optional) Months during the year when maintenance should be performed. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `patching_mode` - (Optional) Cloud Exadata Infrastructure node patching method, either `ROLLING` or `NONROLLING`. Default value is `ROLLING`. IMPORTANT: Non-rolling infrastructure patching involves system down time. See [Oracle-Managed Infrastructure Maintenance Updates](https://docs.cloud.oracle.com/iaas/Content/Database/Concepts/examaintenance.htm#Oracle) for more information. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `preference` - (Optional) The maintenance window scheduling preference. Changing this forces a new Cloud Exadata Infrastructure to be created.

* `weeks_of_month` - (Optional) Weeks during the month when maintenance should be performed. Weeks start on the 1st, 8th, 15th, and 22nd days of the month, and have a duration of 7 days. Weeks start and end based on calendar dates, not days of the week. For example, to allow maintenance during the 2nd week of the month (from the 8th day to the 14th day of the month), use the value 2. Maintenance cannot be scheduled for the fifth week of months that contain more than 28 days. Note that this parameter works in conjunction with the daysOfWeek and hoursOfDay parameters to allow you to specify specific days of the week and hours that maintenance will be performed. Changing this forces a new Cloud Exadata Infrastructure to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cloud Exadata Infrastructure.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Cloud Exadata Infrastructure.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cloud Exadata Infrastructure.
* `update` - (Defaults to 30 minutes) Used when updating the Cloud Exadata Infrastructure.
* `delete` - (Defaults to 1 hour) Used when deleting the Cloud Exadata Infrastructure.

## Import

Cloud Exadata Infrastructures can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_exadata_infrastructure.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup/providers/Oracle.Database/cloudExadataInfrastructures/cloudExadataInfrastructures1
```
