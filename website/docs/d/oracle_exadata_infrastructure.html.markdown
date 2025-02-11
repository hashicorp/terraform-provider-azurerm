---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_exadata_infrastructure"
description: |-
  Gets information about an existing Cloud Exadata Infrastructure.
---

# Data Source: azurerm_oracle_exadata_infrastructure

Use this data source to access information about an existing Cloud Exadata Infrastructure.

## Example Usage

```hcl
data "azurerm_oracle_exadata_infrastructure" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_exadata_infrastructure.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Cloud Exadata Infrastructure.

* `resource_group_name` - (Required) The name of the Resource Group where the Cloud Exadata Infrastructure exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cloud Exadata Infrastructure.

* `activated_storage_count` - The requested number of additional storage servers activated for the Cloud Exadata Infrastructure.

* `additional_storage_count` - The requested number of additional storage servers for the Cloud Exadata Infrastructure.

* `available_storage_size_in_gbs` - The available storage can be allocated to the Cloud Exadata Infrastructure resource, in gigabytes (GB).

* `compute_count` - The number of compute servers for the Cloud Exadata Infrastructure.

* `custom_action_timeout_enabled` - If true, enables the configuration of a custom action timeout (waiting period) between database servers patching operations.

* `cpu_count` - The total number of CPU cores allocated.

* `customer_contacts` - A `customer_contacts` block as defined below.

* `data_storage_size_in_tbs` - The data storage size in terabytes of the DATA disk group.

* `db_node_storage_size_in_gbs` - The local node storage allocated in GBs.

* `db_server_version` - The software version of the database servers (dom0) in the Cloud Exadata Infrastructure.

* `display_name` - The user-friendly name for the Cloud Exadata Infrastructure resource. The name does not need to be unique.

* `estimated_patching_time` - A `estimated_patching_time` block as defined below.

* `last_maintenance_run_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the last maintenance run.

* `lifecycle_details` - Additional information about the current lifecycle state.

* `lifecycle_state` - Cloud Exadata Infrastructure lifecycle state.

* `location` - The Azure Region where the Cloud Exadata Infrastructure exists.

* `maintenance_window` - A `maintenance_window` block as defined below.

* `max_cpu_count` -  The total number of CPU cores available.

* `max_data_storage_in_tbs` - The total available DATA disk group size.

* `max_db_node_storage_size_in_gbs` - The total local node storage available in GBs.

* `max_memory_in_gbs` - The total memory available in GBs.

* `memory_size_in_gbs` - The memory allocated in GBs.

* `monthly_db_server_version` - The monthly software version of the database servers (dom0) in the Cloud Exadata Infrastructure.

* `monthly_patching_enabled` - If true, enables the monthly patching option.

* `monthly_storage_server_version` - The monthly software version of the storage servers (cells) in the Cloud Exadata Infrastructure.

* `next_maintenance_run_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the next maintenance run.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Cloud Exadata Infrastructure.

* `shape` - The model name of the Cloud Exadata Infrastructure resource.

* `storage_count` - The number of storage servers for the Cloud Exadata Infrastructure.

* `storage_server_version` - The software version of the storage servers (cells) in the Cloud Exadata Infrastructure.

* `tags` - A mapping of tags assigned to the Cloud Exadata Infrastructure.

* `time_created` - The date and time the Cloud Exadata Infrastructure resource was created.

* `total_storage_size_in_gbs` -  The total storage allocated to the Cloud Exadata Infrastructure resource, in gigabytes (GB).

* `zones` - The Cloud Exadata Infrastructure Azure zones.

---

A `estimated_patching_time` block exports the following:

* `estimated_db_server_patching_time` - The estimated time required in minutes for database server patching.

* `estimated_network_switches_patching_time` - The estimated time required in minutes for network switch patching.

* `estimated_storage_server_patching_time` - The estimated time required in minutes for storage server patching.

* `total_estimated_patching_time` - The estimated total time required in minutes for all patching operations.

---

A `maintenance_window` block exports the following:

* `days_of_week` - Days during the week when maintenance should be performed.

* `hours_of_day` - The window of hours during the day when maintenance should be performed.

* `lead_time_in_weeks` -  Lead time window allows user to set a lead time to prepare for a down time.

* `months` - A `months` block as defined below.

* `patching_mode` -  Cloud Exadata Infrastructure node patching method.

* `preference` - The maintenance window scheduling preference.

* `weeks_of_month` - Weeks during the month when maintenance should be performed.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cloud Exadata Infrastructure.
