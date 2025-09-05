---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_autonomous_database_clone_from_database"
description: |-
  Gets information about an existing autonomous database clone from database.
---

# Data Source: azurerm_oracle_autonomous_database_clone_from_database

Use this data source to access information about an existing autonomous database clone from database.

## Example Usage

```hcl
data "azurerm_oracle_autonomous_database_clone_from_database" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_autonomous_database_clone_from_database.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this autonomous database clone from database.

* `resource_group_name` - (Required) The name of the Resource Group where the autonomous database cloned from database exists.

## Attributes Reference

In addition to the Arguments listed above—the following Attributes are exported: 

* `id` - The ID of the autonomous database cloned from database.

* `actual_used_data_storage_size_in_tb` - The current amount of storage in use for user and system data, in terabytes (TB).

* `allocated_storage_size_in_tb` - The amount of storage currently allocated for the database tables and billed for, rounded up. When auto-scaling is not enabled, this value is equal to the `dataStorageSizeInTBs` value. You can compare this value to the `actualUsedDataStorageSizeInTBs` value to determine if a manual shrink operation is appropriate for your allocated storage.

* `allowed_ip_addresses` - A list of IP addresses on the access control list.

* `auto_scaling_enabled` - Indicates if auto scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - Indicates if auto scaling is enabled for the Autonomous Database storage.

* `available_upgrade_versions` - A list of Oracle Database versions available for a database upgrade. If there are no version upgrades available, this list is empty.

* `backup_retention_period_in_days` - The backup retention period in days.

* `character_set` - The character set for the autonomous database.

* `compute_count` - The compute amount (CPUs) available to the database.

* `cpu_core_count` - The number of CPU cores to be made available to the database. When the ECPU is selected, the value for cpuCoreCount is 0. 

-> **Note:** For Autonomous Databases on Dedicated Exadata infrastructure, the maximum number of cores is determined by the infrastructure shape. See [Characteristics of Infrastructure Shapes](https://docs.oracle.com/en/cloud/paas/autonomous-database/dedicated/adbde/index.html#GUID-944C9B72-CE8D-48EE-88FB-FDF2A8CB988B) for shape details.

* `data_storage_size_in_gb` - The quantity of data in the database, in gigabytes.

* `data_storage_size_in_tb` - The maximum storage that can be allocated for the database, in terabytes.

* `database_version` - The Oracle Database version for Autonomous Database.

* `display_name` - The display name for the Autonomous Database.

* `failed_data_recovery_in_seconds` - Indicates the number of seconds of data loss for a Data Guard failover.

* `in_memory_area_in_gb` - The area assigned to In-Memory tables in Autonomous Database.

* `lifecycle_details` - Information about the current lifecycle state.

* `local_adg_auto_failover_max_data_loss_limit` - Parameter that allows users to select an acceptable maximum data loss limit in seconds, up to which Automatic Failover will be triggered when necessary for a Local Autonomous Data Guard

* `local_data_guard_enabled` -  Indicates whether the Autonomous Database has local (in-region) Data Guard enabled. Not applicable to cross-region Autonomous Data Guard associations, or to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.

* `location` - The Azure Region where the autonomous database cloned from database exists.

* `memory_per_oracle_compute_unit_in_gb` - The amount of memory in gigabytes per ECPU or OCPU.

* `mtls_connection_required` - Specifies if the Autonomous Database requires mTLS connections.

* `national_character_set` - The national character set for the autonomous database.

* `next_long_term_backup_time_stamp_in_utc` -  The timestamp when the next long-term backup would be created.

* `oci_url` -  The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the autonomous database.

* `peer_database_ids` - The list of [OCIDs](https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of standby databases located in Autonomous Data Guard remote regions that are associated with the source database. Note that for Autonomous Database Serverless instances, standby databases located in the same region as the source primary database do not have OCIDs.

* `preview` - Indicates if the Autonomous Database version is a preview version.

* `preview_version_with_service_terms_accepted` - Indicates if the Autonomous Database version is a preview version with service terms accepted.

* `private_endpoint_url` - The private endpoint for the resource.

* `private_endpoint_ip` - The private endpoint IP address for the resource.

* `private_endpoint_label` - The private endpoint label for the resource.

* `provisionable_cpus` - An array of CPU values that an Autonomous Database can be scaled to.

* `reconnect_clone_enabled` -Indicates whether reconnect clone is enabled.

* `refreshable_clone` - Indicates whether the clone is a refreshable clone.

* `refreshable_status` - The current refreshable status of the clone. Values include `Refreshing` and `NotRefreshing`.

* `remote_data_guard_enabled` -  Indicates whether the Autonomous Database has Cross Region Data Guard enabled. Not applicable to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.

* `service_console_url` - The URL of the Service Console for the Autonomous Database.

* `source_autonomous_database_id` - The ID of the source Autonomous Database from which this clone was created.

* `sql_web_developer_url` - The URL of the SQL web developer portal.

* `subnet_id` - The ID to an Azure Resource Manager subnet the resource is associated with.

* `supported_regions_to_clone_to` -The list of regions that support the creation of an Autonomous Database clone or an Autonomous Data Guard standby database.

* `tags` - A mapping of tags assigned to the autonomous database clone from database.

* `time_created_in_utc` - The timestamp the Autonomous Database was created.

* `time_data_guard_role_changed_utc` - The timestamp the Autonomous Data Guard role was switched for the Autonomous Database. For databases that have standbys in both the primary Data Guard region and a remote Data Guard standby region, this is the latest timestamp of either the database using the "primary" role in the primary Data Guard region, or database located in the remote Data Guard standby region.

* `time_deletion_of_free_autonomous_database_utc` - The timestamp the Always Free database will be automatically deleted because of inactivity. If the database is in the STOPPED state and without activity until this time, it will be deleted.

* `time_local_data_guard_enabled_on_utc` - The timestamp that Autonomous Data Guard was enabled for an Autonomous Database where the standby was provisioned in the same region as the primary database.

* `time_maintenance_begin_utc` -  The timestamp when maintenance will begin.

* `time_maintenance_end_utc` -  The timestamp when maintenance will end.

* `time_of_last_failover_utc` - The timestamp of the last failover operation.

* `time_of_last_refresh_utc` - The timestamp when last refresh happened.

* `time_of_last_refresh_point_utc` - The refresh point timestamp (UTC). The refresh point is the time to which the database was most recently refreshed. Data created after the refresh point is not included in the refresh.

* `time_of_last_switchover_utc` - The timestamp of the last switchover operation for the Autonomous Database.

* `time_reclamation_of_free_autonomous_database_utc` - The timestamp the Always Free database will be stopped because of inactivity. If this time is reached without any database activity, the database will automatically be put into the STOPPED state.

* `time_until_reconnect_in_utc` - The time until reconnect clone is enabled.

* `used_data_storage_size_in_gb` - The storage space consumed by Autonomous Database in GBs.

* `used_data_storage_size_in_tb` - The amount of storage that has been used, in terabytes.

* `virtual_network_id` - The ID to an Azure Resource Manager virtual network resource.

---

A `long_term_backup_schedule` block exports the following:

* `enabled` - A boolean value that indicates if long-term backup is enabled/disabled.

* `repeat_cadence` - The frequency for automated long-term backups.

* `retention_period_in_days` - The retention period in days for Autonomous database backup.

* `time_of_backup_in_utc` - The timestamp in which the backup would be made.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the autonomous database clone from database.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-03-01
