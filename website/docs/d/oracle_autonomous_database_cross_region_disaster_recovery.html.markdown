---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_autonomous_database_cross_region_disaster_recovery"
description: |-
  Gets information about an existing Cross Region Disaster Recovery Autonomous Database.
---

# Data Source: azurerm_oracle_autonomous_database_cross_region_disaster_recovery

Use this data source to access information about an existing peered Cross Region Disaster Recovery Autonomous Database.

## Example Usage

```hcl
data "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_autonomous_database_cross_region_disaster_recovery.example.id
}
```
## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Cross Region Disaster Recovery Autonomous Database.

* `resource_group_name` - (Required) The name of the Resource Group where the Cross Region Disaster Recovery Autonomous Database exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cross Region Disaster Recovery Autonomous Database.

* `actual_used_data_storage_size_in_tb` -  The current amount of storage in use for user and system data in terabytes.

* `allocated_storage_size_in_tb` - The amount of storage currently allocated for the database tables and billed for, rounded up. When auto-scaling is not enabled, this value is equal to the `data_storage_size_in_tb` value. You can compare this value to the `actual_used_data_storage_size_in_tb` value to determine if a manual shrink operation is appropriate for your allocated storage.

* `allowed_ip_addresses` - A list of IP addresses on the access control list.

* `auto_scaling_enabled` - Whether auto-scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - Whether auto-scaling is enabled for the Autonomous Database storage.

* `available_upgrade_versions` - A list of Oracle Database versions available for a database upgrade. If there are no version upgrades available, this list is empty.

* `backup_retention_period_in_days` - The backup retention period in days.

* `character_set` - The character set for the autonomous database.

* `compute_count` - The compute amount (CPUs) available to the database.

* `compute_model` - The compute model of the Autonomous Database.

* `connection_strings` - The connection string used to connect to the Autonomous Database.

* `cpu_core_count` - The number of CPU cores available to the database. When the ECPU is selected, the value for cpuCoreCount is 0.

-> **Note:** For Autonomous Databases on Dedicated Exadata infrastructure, the maximum number of cores is determined by the infrastructure shape. See [Characteristics of Infrastructure Shapes](https://docs.oracle.com/en/cloud/paas/autonomous-database/dedicated/adbde/index.html#GUID-944C9B72-CE8D-48EE-88FB-FDF2A8CB988B) for shape details.

* `customer_contacts` - A list of Customer's contact email addresses.

* `database_node_storage_size_in_gb` - The Database node storage size in, in Gigabytes.

* `data_storage_size_in_gb` - The quantity of data in the database in Gigabytes.

* `data_storage_size_in_tb` - The maximum storage that can be allocated for the database in terabytes.

* `database_version` - The Oracle Database version for Autonomous Database.

* `database_workload` -  The Autonomous Database workload type.

* `display_name` - The user-friendly name for the Cross Region Disaster Recovery Autonomous Database.

* `failed_data_recovery_in_seconds` - Indicates the number of seconds of data loss for Data Guard failover.

* `in_memory_area_in_gb` - The area assigned to In-Memory tables in Autonomous Database.

* `license_model` - The Oracle license model that applied to the Oracle Autonomous Database.

* `lifecycle_details` - Information about the current lifecycle state.

* `lifecycle_state` - The current state of the backup.

* `local_adg_auto_failover_max_data_loss_limit_in_seconds` -  Parameter that allows users to select an acceptable maximum data loss limit in seconds, up to which Automatic Failover will be triggered when necessary for a Local Autonomous Data Guard.

* `local_data_guard_enabled` - Indicates whether the Autonomous Database has local (in-region) Data Guard enabled. Not applicable to cross-region Autonomous Data Guard associations, or to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.

* `location` - The Azure Region where the autonomous database cloned from backup exists.

* `long_term_backup_schedule` - A `long_term_backup_schedule` block as defined below.

* `memory_per_oracle_compute_unit_in_gb` - The amount of memory in Gigabytes per ECPU or OCPU.

* `mtls_connection_required` - Specifies if the Autonomous Database requires mTLS connections.

* `national_character_set` - The national character set for the autonomous database.

* `next_long_term_backup_timestamp_in_utc` -  The timestamp when the next long-term backup would be created.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the autonomous database.

* `peer_database_ids` - The list of [OCIDs](https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of standby databases located in Autonomous Data Guard remote regions that are associated with the source database. Note that for Autonomous Database Serverless instances, standby databases located in the same region as the source primary database do not have OCIDs.

* `preview` - Indicates if the Autonomous Database version is a preview version.

* `preview_version_with_service_terms_accepted` -  Indicates if the Autonomous Database version is a preview version with service terms accepted.

* `private_endpoint_url` - The private endpoint for the resource.

* `private_endpoint_ip` - The private endpoint IP address for the resource.

* `private_endpoint_label` - The private endpoint label for the resource.

* `provisionable_cpus` - An array of CPU values that an Autonomous Database can be scaled to.

* `remote_data_guard_enabled` - Indicates whether the Cross Region Disaster Recovery Autonomous Database has Cross Region Data Guard enabled. Not applicable to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.

* `remote_disaster_recovery_type` - Type of recovery.

* `replicate_automatic_backups_enabled` - If true, 7 days worth of backups are replicated across regions for Cross-Region ADB or Backup-Based Disaster Recovery between Primary and Standby. If false, the backups taken on the Primary are not replicated to the Standby database.

* `source_id` - The Immutable Azure Resource ID of autonomous database for which cross region disaster recovery autonomous database was created.

* `source_location` - The Azure Region where source autonomous database for which cross region disaster recovery autonomous database is located

* `source_ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the autonomous database for which cross region disaster recovery autonomous database was created.

* `subnet_id` - The ID of the Azure Resource Manager subnet resource.

* `time_created_in_utc` - The date and time the Autonomous Database was created.

* `time_data_guard_role_changed_in_utc` - The timestamp the Autonomous Data Guard role was switched for the Autonomous Database. For databases that have standbys in both the primary Data Guard region and a remote Data Guard standby region, this is the latest timestamp of either the database using the "primary" role in the primary Data Guard region, or database located in the remote Data Guard standby region.

* `time_deletion_of_free_autonomous_database_in_utc` - The timestamp the Always Free database will be automatically deleted because of inactivity. If the database is in the STOPPED state and without activity until this time, it will be deleted.

* `time_local_data_guard_enabled_in_utc` - The timestamp that Autonomous Data Guard was enabled for an Autonomous Database where the standby was provisioned in the same region as the primary database.

* `time_maintenance_begin_in_utc` -  The timestamp when maintenance will begin.

* `time_maintenance_end_in_utc` -  The timestamp when maintenance will end.

* `time_of_last_failover_in_utc` - The timestamp of the last failover operation.

* `time_of_last_refresh_in_utc` - The timestamp when the last refresh happened.

* `time_of_last_refresh_point_in_utc` - The refresh point timestamp (UTC). The refresh point is the time to which the database was most recently refreshed. Data created after the refresh point is not included in the refresh.

* `time_of_last_switchover_in_utc` - The timestamp of the last switchover operation for the Autonomous Database.

* `time_reclamation_of_free_autonomous_database_in_utc` - The timestamp the Always Free database will be stopped because of inactivity. If this time is reached without any database activity, the database will automatically be put into the STOPPED state.

* `virtual_network_id` - The ID to an Azure Resource Manager virtual network resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cross Region Disaster Recovery Autonomous Database.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
