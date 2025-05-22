---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_backup"
description: |-
    Manage Autonomous Database Backup.
---

# azurerm_oracle_autonomous_database_backup

Manages an Oracle Autonomous Database Backup in Azure.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}
resource "azurerm_oracle_autonomous_database" "example" {
  name                = "example-adb"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  # ... other autonomous database properties ...
}

resource "azurerm_oracle_autonomous_database_backup" "example" {
  name                     = "example-backup"
  resource_group_name      = azurerm_resource_group.example.name
  autonomous_database_name = azurerm_oracle_autonomous_database.example.name
  retention_period_in_days = 120
  backup_type              = "Full"

  # Optional: specify a custom display name
  display_name             = "My Database Backup"
}

```

## Arguments Reference
The following arguments are supported:

* `name` - (Required) The display name of the Autonomous Database Backup.

* `resource_group_name` - (Required) The name of the Resource Group where the Autonomous Database Backup exists.

* `autonomous_database_` - (Required) The name of the Autonomous Database that this backup is for.
* `retention_period_in_days` - (Required) (Updatable) The number of days to retain the backup. Must be between 90 and 3650 days.

* `backup_type` - (Optional) The type of backup to create. Possible values are Full, Incremental, and LongTerm. Defaults to Full.

* `display_name` - (Optional) The display name of the Autonomous Database Backup. If not specified, the name will be used.

## Attributes Reference
In addition to the Arguments listed above—the following Attributes are exported:

* `id` - The ID of the Autonomous Database Backup.

* `autonomous_database_id` - The OCID of the Autonomous Database OCID.

* `autonomous_database_backup_id`  - The backup OCID.

* `backup_size_in_tbs` - The size of the backup in terabytes.

* `db_version` - The Oracle Database version of the Autonomous Database at the time the backup was taken.

* `is_automatic` - Indicates whether the backup is user-initiated or automatic.

* `is_restorable` - Indicates whether the backup can be used to restore the Autonomous Database.

* `lifecycle_details` - Information about the current lifecycle state of the backup.

* `lifecycle_state` - The current state of the backup.

* `license_model` - The license model of the Autonomous Database at the time the backup was taken.

* `provisioning_state` - The current provisioning state of the Autonomous Database Backup.

* `time_available_til` - The date and time the backup will become unusable.

* `time_ended` - The date and time the backup was completed.

* `time_started` - The date and time the backup started.


## Timeouts
The `timeouts` block allows you to specify timeouts for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Autonomous Database Backup.
* `read` - (Defaults to 5 minutes) Used when retrieving the Autonomous Database Backup.
* `update` - (Defaults to 30 minutes) Used when updating the Autonomous Database Backup.
* `delete` - (Defaults to 30 minutes) Used when deleting the Autonomous Database Backup.
