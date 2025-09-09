---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_autonomous_database_backup"
description: |-
    Gets information about an existing Autonomous Database Backup.
---

# Data Source: azurerm_oracle_autonomous_database_backup

Use this data source to access information about existing Autonomous Database Backups.

## Example Usage

```hcl
data "azurerm_oracle_autonomous_database_backup" "example" {
  autonomous_database_id = azurerm_oracle_autonomous_database.example.id
}

```

## Arguments Reference
The following arguments are supported:

* `autonomous_database_id` - The azureId of the Autonomous Database for we will be listing the backups.

## Attributes Reference
In addition to the Arguments listed above—the following Attributes are exported:

* `id` - The ID of the Autonomous Database Backup.

* `autonomous_database_ocid` - The OCID of the Autonomous Database OCID.

* `autonomous_database_backup_ocid`  - The backup OCID.

* `backup_size_in_tbs` - The size of the backup in terabytes.

* `backup_type` - The type of backup.

* `db_version` - The Oracle Database version of the Autonomous Database at the time the backup was taken.

* `display_name` - The user-friendly name of the backup.

* `is_automatic` - Indicates whether the backup is user-initiated or automatic.

* `is_restorable` - Indicates whether the backup can be used to restore the Autonomous Database.

* `lifecycle_details` - Information about the current lifecycle state of the backup.

* `lifecycle_state` - The current state of the backup.

* `location` `- The Azure Region where the Autonomous Database Backup exists.

* `provisioning_state` - The current provisioning state of the Autonomous Database Backup.

* `retention_period_in_days` - The retention period in days for the Autonomous Database Backup.

* `time_available_til` - The date and time the backup will become unusable.

* `time_ended` - The date and time the backup was completed.

* `time_started` - The date and time the backup started.


## Timeouts
The `timeouts` block allows you to specify timeouts for certain actions:

* `read` - (Defaults to 10 minutes) Used when retrieving the Autonomous Database Backup.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-03-01
