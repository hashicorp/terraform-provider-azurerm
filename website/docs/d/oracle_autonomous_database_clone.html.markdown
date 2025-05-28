---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_autonomous_database_clone"
description: |-
 Gets information about an existing Autonomous Database Clone.
---

# Data Source: azurerm_oracle_autonomous_database_clone

Use this data source to access information about an existing Autonomous Database Clone.

## Example Usage

```hcl
data "azurerm_oracle_autonomous_database_clone" "example" {
  name                = "existing-clone"
  resource_group_name = "existing"
}

output "clone_source_id" {
  value = data.azurerm_oracle_autonomous_database_clone.example.source_id
}

output "clone_type" {
  value = data.azurerm_oracle_autonomous_database_clone.example.clone_type
}

output "lifecycle_state" {
  value = data.azurerm_oracle_autonomous_database_clone.example.lifecycle_state
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Autonomous Database Clone.

* `resource_group_name` - (Required) The name of the Resource Group where the Autonomous Database Clone exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Autonomous Database Clone.

* `location` - The Azure Region where the Autonomous Database Clone exists.

* `tags` - A mapping of tags assigned to the Autonomous Database Clone.

### Clone-Specific Attributes

* `source_id` - The ID of the source Autonomous Database that was cloned.

* `clone_type` - The type of clone. Values include `Full` and `Metadata`.

* `source` - The source of the clone. Values include `None`, `Database`, `BackupFromId`, `BackupFromTimestamp`, `CloneToRefreshable`, `CrossRegionDataguard`, and `CrossRegionDisasterRecovery`.

* `is_refreshable_clone` - Indicates whether the clone is a refreshable clone.

* `refreshable_model` - The refreshable model for the clone. Values include `Automatic` and `Manual`.

* `refreshable_status` - The current refreshable status of the clone. Values include `Refreshing` and `NotRefreshing`.

* `is_reconnect_clone_enabled` - Indicates whether reconnect clone is enabled.

* `time_until_reconnect_clone_enabled` - The time until reconnect clone is enabled.

### Database Attributes

* `autonomous_database_id` - The database [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm).

* `backup_retention_period_in_days` - Retention period, in days, for backups.

* `character_set` - The character set for the autonomous database.

* `compute_count` - The compute amount (CPUs) available to the database.

* `compute_model` - The compute model of the Autonomous Database.

* `connection_strings` - List of connection strings for the Autonomous Database.

* `customer_contacts` - List of customer contact email addresses.

* `data_storage_size_in_gbs` - The quantity of data in the database, in gigabytes.

* `data_storage_size_in_tbs` - The maximum storage that can be allocated for the database, in terabytes.

* `db_version` - A valid Oracle Database version for Autonomous Database.

* `db_workload` - The Autonomous Database workload type.

* `display_name` - The user-friendly name for the Autonomous Database.

* `license_model` - The Oracle license model that applies to the Oracle Autonomous Database.

* `auto_scaling_enabled` - Indicates if auto scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - Indicates if auto scaling is enabled for the Autonomous Database storage.

* `mtls_connection_required` - Specifies if the Autonomous Database requires mTLS connections.

* `national_character_set` - The national character set for the autonomous database.

* `subnet_id` - The ID of the subnet the resource is associated with.

* `virtual_network_id` - The ID of the Virtual Network this Autonomous Database Clone is created in.

* `lifecycle_state` - The current state of the Autonomous Database Clone. Values include `Provisioning`, `Available`, `Stopping`, `Stopped`, `Starting`, `Terminating`, `Terminated`, `Unavailable`, `RestoreInProgress`, `RestoreFailed`, `BackupInProgress`, `ScaleInProgress`, `AvailableNeedsAttention`, `Updating`, `MaintenanceInProgress`, `Restarting`, `Recreating`, `RoleChangeInProgress`, `Upgrading`, `Inaccessible`, and `Standby`.

* `private_endpoint` - The private endpoint for the resource.

* `private_endpoint_ip` - The private endpoint IP address for the resource.

* `service_console_url` - The URL of the Service Console for the Autonomous Database.

* `sql_web_developer_url` - The URL of the SQL web developer.

* `time_created` - The date and time the Autonomous Database Clone was created.

* `oci_url` - The URL of the resource in the OCI console.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Autonomous Database Clone.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database`: 2024-06-01