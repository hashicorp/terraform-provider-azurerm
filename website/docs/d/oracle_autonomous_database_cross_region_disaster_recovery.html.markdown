---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_autonomous_database"
description: |-
  Gets information about an existing Autonomous Database.
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

* `actual_used_data_storage_size_in_tbs` - The current amount of storage in use for user and system data, in terabytes (TB).

* `allocated_storage_size_in_tbs` - The amount of storage currently allocated for the database tables and billed for, rounded up. When auto-scaling is not enabled, this value is equal to the `dataStorageSizeInTBs` value. You can compare this value to the `actualUsedDataStorageSizeInTBs` value to determine if a manual shrink operation is appropriate for your allocated storage.

* `auto_scaling_enabled` - Indicates if auto scaling is enabled for the Cross Region Disaster Recovery Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - Indicates if auto scaling is enabled for the Cross Region Disaster Recovery Autonomous Database storage.

* `backup_retention_period_in_days` - Retention period, in days, for backups.

* `character_set` - The character set for the Cross Region Disaster Recovery Autonomous Database.

* `compute_count` - The compute amount (CPUs) available to the cross region disaster recovery database.
* 
* `compute_model` - The compute model of the Cross Region Disaster Recovery Autonomous Database.

* `database_type` - The type of this Cross Region Disaster Recovery Autonomous Database.

* `db_workload` - The Autonomous Database workload type. The following values are possible:
    * OLTP - indicates an Autonomous Transaction Processing database
    * DW - indicates an Autonomous Data Warehouse database
    * AJD - indicates an Autonomous JSON Database
    * APEX - indicates an Autonomous Database with the Oracle APEX Application Development workload type.

* `display_name` - The user-friendly name for the Cross Region Disaster Recovery Autonomous Database. The name does not have to be unique.

* `id` - The Immutable Azure Resource ID of the Cross Region Disaster Recovery Autonomous Database.

* `license_model` - The Oracle license model that applies to the Oracle Cross Region Disaster Recovery Autonomous Database.

* `lifecycle_details` - Information about the current lifecycle state.

* `local_data_guard_enabled` - Indicates whether the Cross Region Disaster Recovery Autonomous Database has local (in-region) Data Guard enabled. Not applicable to cross-region Autonomous Data Guard associations, or to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.

* `location` - The Azure Region where the Cross Region Disaster Recovery Autonomous Database exists.

* `mtls_connection_required` - Specifies if the Cross Region Disaster Recovery Autonomous Database requires mTLS connections.

* `name` - The name of this Cross Region Disaster Recovery Autonomous Database.

* `national_character_set` - The national character set for the Cross Region Disaster Recovery Autonomous Database.  The default is AL16UTF16. Allowed values are: AL16UTF16 or UTF8.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the autonomous database.

* `peer_db_ids` - The list of Immutable Azure Resource IDs of Autonomous Databases peered with this Cross Region Disaster Recovery Autonomous Database.

* `private_endpoint` - The private endpoint for the resource.

* `private_endpoint_ip` - The private endpoint Ip address for the resource.

* `private_endpoint_label` - The private endpoint label for the resource.

* `remote_data_guard_enabled` - Indicates whether the Cross Region Disaster Recovery Autonomous Database has Cross Region Data Guard enabled. Not applicable to Autonomous Databases using dedicated Exadata infrastructure or Exadata Cloud@Customer infrastructure.

* `remote_disaster_recovery_type` - Type of recovery. Value can be either `Adg` (Autonomous Data Guard) or `BackupBased`.

* `replicate_automatic_backups` - If true, 7 days worth of backups are replicated across regions for Cross-Region ADB or Backup-Based Disaster Recovery between Primary and Standby. If false, the backups taken on the Primary are not replicated to the Standby database.

* `source` - The source of the database. For cross region disaster recovery autonomous database the source is always "CrossRegionDisasterRecovery" as this type of autonomous database can be created only as disaster recovery solution for another autonomous database.

* `source_id` - The Immutable Azure Resource ID of autonomous database for which cross region disaster recovery autonomous database was created.

* `source_location` - The Azure Region where source autonomous database for which cross region disaster recovery autonomous database is located

* `source_ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the autonomous database for which cross region disaster recovery autonomous database was created.

* `subnet_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the subnet the resource is associated with.

* `supported_regions_to_clone_to` - The list of regions that support the creation Cross Region Disaster Recovery clone. List is always empty as Cross Region Disaster Recovery Autonomous Database cannot be cloned.

* `time_created` - The date and time the Autonomous Database was created.

* `time_data_guard_role_changed` - The date and time the Autonomous Data Guard role was switched for the Cross Region Disaster Recovery Autonomous Database. For databases that have standbys in both the primary Data Guard region and a remote Data Guard standby region, this is the latest timestamp of either the database using the "primary" role in the primary Data Guard region, or database located in the remote Data Guard standby region.

* `time_maintenance_begin` - The date and time when maintenance will begin.

* `time_maintenance_end` - The date and time when maintenance will end.

* `virtual_network_id` - The ID to an Azure Resource Manager vnet resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cross Region Disaster Recovery Autonomous Database.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database`: 2025-03-01
