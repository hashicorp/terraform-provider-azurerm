---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_autonomous_database_cross_region_disaster_recovery"
description: |-
  Gets information about an existing Cross Region Disaster Recovery Autonomous Database.
---

# Data Source: azurerm_oracle_autonomous_database_cross_region_disaster_recovery

Use this data source to access information about an existing Cross Region Disaster Recovery Autonomous Database.

## Example Usage

```hcl
data "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "example" {
  name                = "example"
  resource_group_name = "example-resources"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Cross Region Disaster Recovery Autonomous Database.

* `resource_group_name` - (Required) The name of the Resource Group where the Cross Region Disaster Recovery Autonomous Database exists.

## Attributes Reference

In addition to the Arguments listed above the following Attributes are exported:

* `id` - The ID of the Cross Region Disaster Recovery Autonomous Database.

* `actual_used_data_storage_size_in_tb` - The current amount of storage in use for user and system data in terabytes.

* `allocated_storage_size_in_tb` - The amount of storage currently allocated for the database tables and billed for, rounded up.

* `allowed_ip_addresses` - A list of IP addresses on the access control list.

* `auto_scaling_enabled` - Whether auto-scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - Whether auto-scaling is enabled for the Autonomous Database storage.

* `available_upgrade_versions` - A list of Oracle Database versions available for a database upgrade.

* `backup_retention_period_in_days` - The backup retention period in days.

* `character_set` - The character set for the Autonomous Database.

* `compute_count` - The compute amount available to the database.

* `cpu_core_count` - The number of CPU cores available to the database.

* `customer_contacts` - A list of customer contact email addresses.

* `data_storage_size_in_gb` - The quantity of data in the database in gigabytes.

* `data_storage_size_in_tb` - The maximum storage that can be allocated for the database in terabytes.

* `database_workload` - The Autonomous Database workload type.

* `database_type` - The database type.

* `database_version` - The Oracle Database version for the Autonomous Database.

* `display_name` - The user-friendly name for the Cross Region Disaster Recovery Autonomous Database.

* `failed_data_recovery_in_seconds` - Indicates the number of seconds of data loss for Data Guard failover.

* `in_memory_area_in_gb` - The area assigned to in-memory tables in the Autonomous Database.

* `lifecycle_details` - Information about the current lifecycle state.

* `local_adg_auto_failover_max_data_loss_limit` - The maximum data loss limit in seconds for local ADG automatic failover.

* `local_data_guard_enabled` - Indicates whether local (in-region) Data Guard is enabled.

* `license_model` - The Oracle license model that applies to the Autonomous Database.

* `location` - The Azure Region where the Cross Region Disaster Recovery Autonomous Database exists.

* `memory_per_oracle_compute_unit_in_gb` - The amount of memory in gigabytes per ECPU or OCPU.

* `mtls_connection_required` - Specifies if the Autonomous Database requires mTLS connections.

* `national_character_set` - The national character set for the Autonomous Database.

* `next_long_term_backup_timestamp` - The timestamp when the next long-term backup would be created.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The OCID of the Autonomous Database.

* `peer_database_ids` - The IDs of standby databases associated with the source database.

* `preview` - Indicates whether the Autonomous Database version is a preview version.

* `preview_version_with_service_terms_accepted` - Indicates whether the preview version terms were accepted.

* `private_endpoint_ip` - The private endpoint IP address for the resource.

* `private_endpoint_label` - The private endpoint label for the resource.

* `private_endpoint_url` - The private endpoint URL for the resource.

* `provisionable_cpus` - An array of CPU values that an Autonomous Database can be scaled to.

* `remote_data_guard_enabled` - Indicates whether cross-region Data Guard is enabled.

* `remote_disaster_recovery_type` - Type of recovery.

* `replicate_automatic_backups_enabled` - If true, backups are replicated across regions.

* `service_console_url` - The URL of the service console.

* `source` - The source type used for the Cross Region Disaster Recovery Autonomous Database.

* `source_autonomous_database_id` - The ID of the source (primary) Autonomous Database.

* `source_location` - The Azure Region where the source Autonomous Database is located.

* `source_ocid` - The OCID of the source Autonomous Database.

* `sql_web_developer_url` - The URL of SQL Web Developer.

* `subnet_id` - The ID of the subnet associated with this database.

* `tags` - A mapping of tags assigned to the resource.

* `time_created_in_utc` - The date and time the Autonomous Database was created.

* `time_data_guard_role_changed_in_utc` - The timestamp when the Data Guard role was last switched.

* `time_deletion_of_free_autonomous_database_in_utc` - The timestamp when an inactive Always Free database is deleted.

* `time_local_data_guard_enabled_in_utc` - The timestamp when local Data Guard was enabled.

* `time_maintenance_begin_in_utc` - The timestamp when maintenance begins.

* `time_maintenance_end_in_utc` - The timestamp when maintenance ends.

* `time_of_last_failover_in_utc` - The timestamp of the last failover operation.

* `time_of_last_refresh_in_utc` - The timestamp when the last refresh happened.

* `time_of_last_refresh_point_in_utc` - The refresh point timestamp.

* `time_of_last_switchover_in_utc` - The timestamp of the last switchover operation.

* `time_reclamation_of_free_autonomous_database_in_utc` - The timestamp when an inactive Always Free database is stopped.

* `used_data_storage_size_in_gb` - The amount of used data storage in gigabytes.

* `used_data_storage_size_in_tb` - The amount of used data storage in terabytes.

* `virtual_network_id` - The ID of the virtual network associated with this database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cross Region Disaster Recovery Autonomous Database.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
