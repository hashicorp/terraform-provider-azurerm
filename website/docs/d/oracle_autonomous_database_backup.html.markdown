---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_autonomous_database_backup"
description: |-
    Gets information about an existing Autonomous Database Backup.
---

# Data Source: azurerm_oracle_autonomous_database_backup

Use this data source to access information about an existing Autonomous Database Backup.

## Example Usage

```hcl
data "azurerm_oracle_autonomous_database_backup" "example" {
  name                     = "existing"
  resource_group_name      = "existing"
  autonomous_database_name = "existingadb"
}

output "id" {
  value = data.azurerm_oracle_autonomous_database_backup.example.id
}
```

## Arguments Reference
The following arguments are supported:

* `name` - (Required) The display name of the Autonomous Database Backup.

* `resource_group_name` - (Required) The name of the Resource Group where the Autonomous Database Backup exists.

* `autonomous_database_name` - (Required) The name of the Autonomous Database that this backup is for.

## Attributes Reference
In addition to the Arguments listed above—the following Attributes are exported:

* `id` - The ID of the Autonomous Database Backup.

* `autonomous_database_id` - The OCID of the Autonomous Database OCID.

* `autonomous_database_backup_id`  - The backup OCID.

* `backup_size_in_tbs` - The size of the backup in terabytes.

* `backup_type` - The type of backup.

* `db_version` - The Oracle Database version of the Autonomous Database at the time the backup was taken.

* `display_name` - The user-friendly name of the backup.

* `is_automatic` - Indicates whether the backup is user-initiated or automatic.

* `is_restorable` - Indicates whether the backup can be used to restore the Autonomous Database.

* `lifecycle_details` - Information about the current lifecycle state of the backup.

* `lifecycle_state` - The current state of the backup.

* `location` `- The Azure Region where the Autonomous Database Backup exists.

* `license_model` - The license model of the Autonomous Database at the time the backup was taken.

* `provisioning_state` - The current provisioning state of the Autonomous Database Backup.

* `retention_period_in_days` - The retention period in days for the Autonomous Database Backup.

* `time_available_til` - The date and time the backup will become unusable.

* `time_ended` - The date and time the backup was completed.

* `time_started` - The date and time the backup started.


## Timeouts
The `timeouts` block allows you to specify timeouts for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Autonomous Database Backup.
