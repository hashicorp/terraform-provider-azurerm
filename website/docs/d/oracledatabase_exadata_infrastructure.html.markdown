---
subcategory: "Oracle Database"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracledatabase_exadata_infrastructure"
description: |-
  Gets information about an existing Exadata Infrastructure.
---

# Data Source: azurerm_oracledatabase_exadata_infrastructure

Use this data source to access information about an existing Exadata Infrastructure.

## Example Usage

```hcl
data "azurerm_oracledatabase_exadata_infrastructure" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracledatabase_exadata_infrastructure.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Exadata Infrastructure.

* `resource_group_name` - (Required) The name of the Resource Group where the Exadata Infrastructure exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Exadata Infrastructure.

* `activated_storage_count` - The requested number of additional storage servers activated for the Exadata infrastructure.

* `additional_storage_count` - The requested number of additional storage servers for the Exadata infrastructure.

* `available_storage_size_in_gbs` - The available storage can be allocated to the cloud Exadata infrastructure resource, in gigabytes (GB).

* `compute_count` - The number of compute servers for the cloud Exadata infrastructure.

* `cpu_count` - The total number of CPU cores allocated.

* `customer_contacts` - A `customer_contacts` block as defined below.

* `data_storage_size_in_tbs` - Size, in terabytes, of the DATA disk group.

* `db_node_storage_size_in_gbs` - The local node storage allocated in GBs.

* `db_server_version` - The software version of the database servers (dom0) in the cloud Exadata infrastructure. Example: 20.1.15

* `display_name` - The user-friendly name for the cloud Exadata infrastructure resource. The name does not need to be unique.

* `estimated_patching_time` - A `estimated_patching_time` block as defined below.

* `last_maintenance_run_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the last maintenance run.

* `lifecycle_details` - Additional information about the current lifecycle state.

* `lifecycle_state` - CloudExadataInfrastructure lifecycle state.

* `location` - The Azure Region where the Exadata Infrastructure exists.

* `maintenance_window` - A `maintenance_window` block as defined below.

* `max_cpu_count` -  The total number of CPU cores available.

* `max_data_storage_in_tbs` - The total available DATA disk group size.

* `max_db_node_storage_size_in_gbs` - The total local node storage available in GBs.

* `max_memory_in_gbs` - The total memory available in GBs.

* `memory_size_in_gbs` - The memory allocated in GBs.

* `monthly_db_server_version` - The monthly software version of the database servers (dom0) in the cloud Exadata infrastructure. Example: 20.1.15

* `monthly_storage_server_version` - The monthly software version of the storage servers (cells) in the cloud Exadata infrastructure. Example: 20.1.15

* `next_maintenance_run_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the next maintenance run.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Exadata infrastructure.

* `provisioning_state` - CloudExadataInfrastructure provisioning state

* `shape` - The model name of the cloud Exadata infrastructure resource.

* `storage_count` - The number of storage servers for the cloud Exadata infrastructure.

* `storage_server_version` - The software version of the storage servers (cells) in the Exadata infrastructure. Example: 20.1.15

* `system_data` - A `system_data` block as defined below.

* `tags` - A mapping of tags assigned to the Exadata Infrastructure.

* `time_created` - The date and time the cloud Exadata infrastructure resource was created.

* `total_storage_size_in_gbs` -  The total storage allocated to the cloud Exadata infrastructure resource, in gigabytes (GB).

* `zones` - The Exadata infrastructure Azure zones.

---

A `estimated_patching_time` block exports the following:

* `estimated_db_server_patching_time` - The estimated time required in minutes for database server patching.

* `estimated_network_switches_patching_time` - The estimated time required in minutes for network switch patching.

* `estimated_storage_server_patching_time` - The estimated time required in minutes for storage server patching.

* `total_estimated_patching_time` - The estimated total time required in minutes for all patching operations.

---

A `maintenance_window` block exports the following:

* `days_of_week` - Days during the week when maintenance should be performed.

* `hours_of_day` - The window of hours during the day when maintenance should be performed. The window is a 4 hour slot. Valid values are: 0 - represents time slot 0:00 - 3:59 UTC - 4 - represents time slot 4:00 - 7:59 UTC - 8 - represents time slot 8:00 - 11:59 UTC - 12 - represents time slot 12:00 - 15:59 UTC - 16 - represents time slot 16:00 - 19:59 UTC - 20 - represents time slot 20:00 - 23:59 UTC

* `lead_time_in_weeks` -  Lead time window allows user to set a lead time to prepare for a down time. The lead time is in weeks and valid value is between 1 to 4.

* `months` - A `months` block as defined below.

* `patching_mode` -  Cloud Exadata infrastructure node patching method, either "ROLLING" or "NONROLLING".

* `preference` - The maintenance window scheduling preference.

* `weeks_of_month` - Weeks during the month when maintenance should be performed. Weeks start on the 1st, 8th, 15th, and 22nd days of the month, and have a duration of 7 days. Weeks start and end based on calendar dates, not days of the week. For example, to allow maintenance during the 2nd week of the month (from the 8th day to the 14th day of the month), use the value 2. Maintenance cannot be scheduled for the fifth week of months that contain more than 28 days. Note that this parameter works in conjunction with the daysOfWeek and hoursOfDay parameters to allow you to specify specific days of the week and hours that maintenance will be performed.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Exadata Infrastructure.
