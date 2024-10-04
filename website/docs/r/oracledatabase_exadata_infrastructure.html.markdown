---
subcategory: "Oracle Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracledatabase_exadata_infrastructure"
description: |-
  Manages a Exadata Infrastructure.
---

# azurerm_oracledatabase_exadata_infrastructure

Manages a Exadata Infrastructure.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_oracledatabase_exadata_infrastructure" "example" {
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

* `compute_count` - (Required) The number of compute servers for the cloud Exadata infrastructure.

* `display_name` - (Required) The user-friendly name for the cloud Exadata infrastructure resource. The name does not need to be unique.

* `location` - (Required) The Azure Region where the Exadata Infrastructure should exist. Changing this forces a new Exadata Infrastructure to be created.

* `name` - (Required) The name which should be used for this Exadata Infrastructure.

* `resource_group_name` - (Required) The name of the Resource Group where the Exadata Infrastructure should exist.

* `shape` - (Required) The shape of the cloud Exadata infrastructure resource.

* `storage_count` - (Required) The number of storage servers for the cloud Exadata infrastructure.

* `zones` - (Required) CloudExadataInfrastructure zones.

---

* `customer_contacts` - (Optional) The email address used by Oracle to send notifications regarding databases and infrastructure.

* `maintenance_window` - (Optional) One or more `maintenance_window` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Exadata Infrastructure.

---

A `maintenance_window` block supports the following:

* `days_of_week` - (Optional) Days during the week when maintenance should be performed. Valid values are: `0` - represents time slot `0:00 - 3:59 UTC - 4` - represents time slot `4:00 - 7:59 UTC - 8` - represents time slot 8:00 - 11:59 UTC - 12 - represents time slot 12:00 - 15:59 UTC - 16 - represents time slot 16:00 - 19:59 UTC - 20 - represents time slot `20:00 - 23:59 UTC`

* `hours_of_day` - (Optional) The window of hours during the day when maintenance should be performed. The window is a 4 hour slot.

* `lead_time_in_weeks` - (Optional) Lead time window allows user to set a lead time to prepare for a down time. The lead time is in weeks and valid value is between `1` to `4`.

* `months` - (Optional) Months during the year when maintenance should be performed.

* `patching_mode` - (Optional) Cloud Exadata infrastructure node patching method, either "ROLLING" or "NONROLLING". Default value is ROLLING. IMPORTANT: Non-rolling infrastructure patching involves system down time. See [Oracle-Managed Infrastructure Maintenance Updates](https://docs.cloud.oracle.com/iaas/Content/Database/Concepts/examaintenance.htm#Oracle) for more information.

* `preference` - (Optional) The maintenance window scheduling preference.

* `weeks_of_month` - (Optional) Weeks during the month when maintenance should be performed. Weeks start on the 1st, 8th, 15th, and 22nd days of the month, and have a duration of 7 days. Weeks start and end based on calendar dates, not days of the week. For example, to allow maintenance during the 2nd week of the month (from the 8th day to the 14th day of the month), use the value 2. Maintenance cannot be scheduled for the fifth week of months that contain more than 28 days. Note that this parameter works in conjunction with the daysOfWeek and hoursOfDay parameters to allow you to specify specific days of the week and hours that maintenance will be performed.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Exadata Infrastructure.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Exadata Infrastructure.
* `read` - (Defaults to 5 minutes) Used when retrieving the Exadata Infrastructure.
* `update` - (Defaults to 30 minutes) Used when updating the Exadata Infrastructure.
* `delete` - (Defaults to 30 minutes) Used when deleting the Exadata Infrastructure.

## Import

Exadata Infrastructures can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracledatabase_exadata_infrastructure.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup/providers/Oracle.Database/cloudExadataInfrastructures/cloudExadataInfrastructures1
```
