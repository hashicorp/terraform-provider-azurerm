---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_autonomous_database"
description: |-
  Gets information about an existing Autonomous Database.
---

# Data Source: azurerm_oracle_autonomous_database

Use this data source to access information about an existing Autonomous Database.

## Example Usage

```hcl
data "azurerm_oracle_autonomous_database" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_autonomous_database.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Autonomous Database.

* `resource_group_name` - (Required) The name of the Resource Group where the Autonomous Database exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Autonomous Database.

* `actual_used_data_storage_size_in_tbs` - The current amount of storage in use for user and system data, in terabytes (TB).

* `allocated_storage_size_in_tbs` - The amount of storage currently allocated for the database tables and billed for, rounded up. When auto-scaling is not enabled, this value is equal to the `dataStorageSizeInTBs` value. You can compare this value to the `actualUsedDataStorageSizeInTBs` value to determine if a manual shrink operation is appropriate for your allocated storage.

* `autonomous_database_id` - The database [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm).

* `available_upgrade_versions` - List of Oracle Database versions available for a database upgrade. If there are no version upgrades available, this list is empty.

* `backup_retention_period_in_days` - Retention period, in days, for backups.

* `character_set` - The character set for the autonomous database. 

* `compute_count` - The compute amount (CPUs) available to the database.

* `cpu_core_count` - The number of CPU cores to be made available to the database. When the ECPU is selected, the value for cpuCoreCount is 0. For Autonomous Database on Dedicated Exadata infrastructure, the maximum number of cores is determined by the infrastructure shape. See [Characteristics of Infrastructure Shapes](https://www.oracle.com/pls/topic/lookup?ctx=en/cloud/paas/autonomous-database&id=ATPFG-GUID-B0F033C1-CC5A-42F0-B2E7-3CECFEDA1FD1) for shape details.

* `data_storage_size_in_gbs` - The quantity of data in the database, in gigabytes.

* `data_storage_size_in_tbs` - The maximum storage that can be allocated for the database, in terabytes.

* `db_node_storage_size_in_gbs` - The DB node storage size in, in gigabytes.

* `db_version` - A valid Oracle Database version for Autonomous Database.

* `display_name` - The user-friendly name for the Autonomous Database. The name does not have to be unique.

* `failed_data_recovery_in_seconds` - Indicates the number of seconds of data loss for a Data Guard failover.

* `in_memory_area_in_gbs` - The area assigned to In-Memory tables in Autonomous Database.

* `auto_scaling_enabled` - Indicates if auto scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - Indicates if auto scaling is enabled for the Autonomous Database storage.

* `local_data_guard_enabled` - Indicates whether the Autonomous Database has local (in-region) Data Guard enabled. Not applicable to cross-region Autonomous Data Guard associations, or to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.

* `mtls_connection_required` - Specifies if the Autonomous Database requires mTLS connections.

* `preview` - Indicates if the Autonomous Database version is a preview version.

* `preview_version_with_service_terms_accepted` - Indicates if the Autonomous Database version is a preview version with service terms accepted.

* `remote_data_guard_enabled` - Indicates whether the Autonomous Database has Cross Region Data Guard enabled. Not applicable to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.

* `key_history_entry` - Key History Entry.

* `lifecycle_details` - Information about the current lifecycle state.

* `local_adg_auto_failover_max_data_loss_limit` - Parameter that allows users to select an acceptable maximum data loss limit in seconds, up to which Automatic Failover will be triggered when necessary for a Local Autonomous Data Guard

* `location` - The Azure Region where the Autonomous Database exists.

* `memory_per_oracle_compute_unit_in_gbs` - The amount of memory (in GBs) enabled per ECPU or OCPU.

* `national_character_set` - The national character set for the autonomous database.  The default is AL16UTF16. Allowed values are: AL16UTF16 or UTF8.

* `next_long_term_backup_time_stamp` - The date and time when the next long-term backup would be created.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the autonomous database.

* `peer_db_ids` - The list of [OCIDs](https://docs.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of standby databases located in Autonomous Data Guard remote regions that are associated with the source database. Note that for Autonomous Database Serverless instances, standby databases located in the same region as the source primary database do not have OCIDs.

* `private_endpoint` - The private endpoint for the resource.

* `private_endpoint_ip` - The private endpoint Ip address for the resource.

* `private_endpoint_label` - The private endpoint label for the resource.

* `provisionable_cpus` - An array of CPU values that an Autonomous Database can be scaled to.

* `service_console_url` - The URL of the Service Console for the Autonomous Database.

* `sql_web_developer_url` - The URL of the SQL web developer.

* `subnet_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the subnet the resource is associated with.

* `supported_regions_to_clone_to` - The list of regions that support the creation of an Autonomous Database clone or an Autonomous Data Guard standby database.

* `tags` - A mapping of tags assigned to the Autonomous Database.

* `time_created` - The date and time the Autonomous Database was created.

* `time_data_guard_role_changed` - The date and time the Autonomous Data Guard role was switched for the Autonomous Database. For databases that have standbys in both the primary Data Guard region and a remote Data Guard standby region, this is the latest timestamp of either the database using the "primary" role in the primary Data Guard region, or database located in the remote Data Guard standby region.

* `time_deletion_of_free_autonomous_database` - The date and time the Always Free database will be automatically deleted because of inactivity. If the database is in the STOPPED state and without activity until this time, it will be deleted.

* `time_local_data_guard_enabled_on` - The date and time that Autonomous Data Guard was enabled for an Autonomous Database where the standby was provisioned in the same region as the primary database.

* `time_maintenance_begin` - The date and time when maintenance will begin.

* `time_maintenance_end` - The date and time when maintenance will end.

* `time_of_last_failover` - The timestamp of the last failover operation.

* `time_of_last_refresh` - The date and time when last refresh happened.

* `time_of_last_refresh_point` - The refresh point timestamp (UTC). The refresh point is the time to which the database was most recently refreshed. Data created after the refresh point is not included in the refresh.

* `time_of_last_switchover` - The timestamp of the last switchover operation for the Autonomous Database.

* `time_reclamation_of_free_autonomous_database` - The date and time the Always Free database will be stopped because of inactivity. If this time is reached without any database activity, the database will automatically be put into the STOPPED state.

* `used_data_storage_size_in_gbs` - The storage space consumed by Autonomous Database in GBs.

* `used_data_storage_size_in_tbs` - The amount of storage that has been used, in terabytes.

* `virtual_network_id` - The ID to an Azure Resource Manager vnet resource.

* `allowed_ips` - The client IP access control list (ACL). This feature is available for [Autonomous Database Serverless] (https://docs.oracle.com/en/cloud/paas/autonomous-database/index.html) and on Exadata Cloud@Customer. Only clients connecting from an IP address included in the ACL may access the Autonomous Database instance. If `arePrimaryWhitelistedIpsUsed` is 'TRUE' then Autonomous Database uses this primary's IP access control list (ACL) for the disaster recovery peer called `standbywhitelistedips`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Autonomous Database.
