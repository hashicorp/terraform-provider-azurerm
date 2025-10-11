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
  autonomous_database_id   = azurerm_oracle_autonomous_database.example.id
  retention_period_in_days = 120
  backup_type              = "Full"
}

```

## Argument Reference
The following arguments are supported:

* `name` - (Required) The display name of the Autonomous Database Backup. Changing this forces a new resource to be created.

* `autonomous_database_id` - (Required) The azureId of the Autonomous Database that this backup is for. Changing this forces a new resource to be created.

* `retention_period_in_days` - (Required) (Updatable) The number of days to retain the backup. Must be between 90 and 3650 days.

* `type` - (Optional) The type of backup to create.Currently, only `LongTerm` backup operations are supported through the Oracle database At azure service.

## Attribute Reference
In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Autonomous Database Backup.

## Timeouts
The `timeouts` block allows you to specify timeouts for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Autonomous Database Backup.
* `read` - (Defaults to 5 minutes) Used when retrieving the Autonomous Database Backup.
* `update` - (Defaults to 30 minutes) Used when updating the Autonomous Database Backup.
* `delete` - (Defaults to 30 minutes) Used when deleting the Autonomous Database Backup.

## Import

Autonomous Database Backups can be imported using the `id`, e.g.

```shell
terraform import azurerm_oracle_autonomous_database_backup.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup/providers/Oracle.Database/autonomousDatabases/autonomousDatabase1/autonomousDatabaseBackups/autonomousDatabaseBackup1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-03-01
