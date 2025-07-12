---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_clone"
description: |-
  Manages an Oracle Autonomous Database Clone.
---

# azurerm_oracle_autonomous_database_clone

Manages an Oracle Autonomous Database Clone.

## Example Usage

```hcl
resource "azurerm_oracle_autonomous_database_clone" "example" {
  name                = "example-clone-db"
  resource_group_name = "example-rg"
  location            = "East US"

  # Clone-specific configuration
  source_id      = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-rg/providers/Oracle.Database/autonomousDatabases/source-db"
  clone_type     = "Full"
  source         = "Database"
  data_base_type = "Clone"

  # Optional clone features
  is_refreshable_clone = true
  refreshable_model    = "Manual"

  # Database configuration
  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "Example Clone Database"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/example-vnet/subnets/oracle-subnet"
  virtual_network_id               = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/example-vnet"
  customer_contacts                = ["admin@example.com"]
  tags = {
    Environment = "Development"
    Purpose     = "Clone"
  }
}

### Clone from Specific Backup Timestamp
resource "azurerm_oracle_autonomous_database_clone" "specific_backup_clone" {
  name                = "example-specific-backup-clone"
  resource_group_name = "example-rg"
  location            = "East US"
  # Backup timestamp clone configuration
  source_id      = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-rg/providers/Oracle.Database/autonomousDatabases/source-db"
  clone_type     = "Full"
  source         = "BackupFromTimestamp"
  data_base_type = "CloneFromBackupTimestamp"
  # Specific timestamp for point-in-time recovery
  # or set  use_latest_available_backup_time_stamp = true to clone from the latest backup
  timestamp = "2024-12-01T12:00:00Z"
  # Database configuration (all required fields)
  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "Point-in-time Clone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/example-vnet/subnets/oracle-subnet"
  virtual_network_id               = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/example-vnet"
  tags = {
    Environment = "Production"
    Purpose     = "BackupClone"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Autonomous Database Clone. Changing this forces a new Autonomous Database Clone to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Autonomous Database Clone should exist. Changing this forces a new Autonomous Database Clone to be created.

* `location` - (Required) The Azure Region where the Autonomous Database Clone should exist. Changing this forces a new Autonomous Database Clone to be created.

### Clone  Specific Arguments

* `source_id` - (Required) The ID of the source Autonomous Database to clone from. Changing this forces a new Autonomous Database Clone to be created.

* `clone_type` - (Required) The type of clone to create. Possible values are `Full` and `Metadata`. Changing this forces a new Autonomous Database Clone to be created.

* `source` - (Required) The source of the clone. Possible values are  `Database`, `BackupFromId`, `BackupFromTimestamp`, `CloneToRefreshable`. Changing this forces a new Autonomous Database Clone to be created.

* `data_base_type` - (Required) The database type for the clone. Possible values are `Clone` and `CloneFromBackupTimestamp`. Changing this forces a new Autonomous Database Clone to be created.

### Clone from database Configuration (Optional)

* `refreshable_model` - (Optional) The refreshable model for the clone. Possible values are `Automatic` and `Manual`. 

* `time_until_reconnect_clone_enabled` - (Optional) The time until reconnect clone is enabled. Must be in RFC3339 format.

### Clone from Backup Timestamp specific Configuration (Optional)

* `timestamp` - (Optional) The timestamp specified for the point-in-time clone of the source Autonomous Database. The timestamp must be in the past and in RFC3339 format. Only applicable when `data_base_type` is `CloneFromBackupTimestamp`. Changing this forces a new Autonomous Database Clone to be created.

* `use_latest_available_backup_time_stamp` - (Optional) Clone from latest available backup timestamp. Only applicable when `data_base_type` is `CloneFromBackupTimestamp`. Changing this forces a new Autonomous Database Clone to be created.

### Database Configuration Arguments

* `admin_password` - (Required) The password for the SYS, SYSTEM, and PDB Admin users. The password must be at least 12 characters long, and contain at least 1 uppercase, 1 lowercase, and 1 numeric character. It cannot contain the double quote symbol (") or the username "admin", regardless of casing. Changing this forces a new Autonomous Database Clone to be created.

* `backup_retention_period_in_days` - (Required) Retention period, in days, for backups. Possible values are between `1` and `60`. Changing this forces a new Autonomous Database Clone to be created.

* `character_set` - (Required) The character set for the autonomous database. Changing this forces a new Autonomous Database Clone to be created.

* `compute_count` - (Required) The compute amount (CPUs) available to the database. Possible values are between `2.0` and `512.0`.

* `compute_model` - (Required) The compute model of the Autonomous Database. Changing this forces a new Autonomous Database Clone to be created.

* `data_storage_size_in_tbs` - (Required) The maximum storage that can be allocated for the database, in terabytes. Possible values are between `1` and `384`.

* `db_version` - (Required) A valid Oracle Database version for Autonomous Database. Changing this forces a new Autonomous Database Clone to be created.

* `db_workload` - (Required) The Autonomous Database workload type. Possible values are `OLTP` and `DW`, `APEX`, `AJD`. Changing this forces a new Autonomous Database Clone to be created.
    * OLTP - indicates an Autonomous Transaction Processing database
    * DW - indicates an Autonomous Data Warehouse database
    * AJD - indicates an Autonomous JSON Database
    * APEX - indicates an Autonomous Database with the Oracle APEX Application Development workload type. 

~> **Note:** To clone the database with a different `db_workload` type, please refer to the documentation [here](https://docs.public.oneportal.content.oci.oraclecloud.com/en-us/iaas/autonomous-database-serverless/doc/autonomous-clone-cross-workload-type.html#GUID-527A712D-FF82-498B-AB35-8A1623E36EDD) for correct configuration steps.


* `display_name` - (Required) The user-friendly name for the Autonomous Database. Changing this forces a new Autonomous Database Clone to be created.

* `license_model` - (Required) The Oracle license model that applies to the Oracle Autonomous Database. Possible values are `LicenseIncluded` and `BringYourOwnLicense`. Changing this forces a new Autonomous Database Clone to be created.

* `auto_scaling_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database storage.

* `mtls_connection_required` - (Required) Specifies if the Autonomous Database requires mTLS connections. Changing this forces a new Autonomous Database Clone to be created.

* `national_character_set` - (Required) The national character set for the autonomous database. Changing this forces a new Autonomous Database Clone to be created.

* `subnet_id` - (Required) The ID of the subnet the resource is associated with. Changing this forces a new Autonomous Database Clone to be created.

* `virtual_network_id` - (Required) The ID of the Virtual Network this Autonomous Database Clone should be created in. Changing this forces a new Autonomous Database Clone to be created.

* `customer_contacts` - (Optional) Customer contact email addresses. Changing this forces a new Autonomous Database Clone to be created.

* `tags` - (Optional) A mapping of tags to assign to the Autonomous Database Clone.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Autonomous Database Clone.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 120 minutes) Used when creating the Autonomous Database Clone.
* `read` - (Defaults to 5 minutes) Used when retrieving the Autonomous Database Clone.
* `delete` - (Defaults to 30 minutes) Used when deleting the Autonomous Database Clone.

## Import

Autonomous Database Clones can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_autonomous_database_clone.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Oracle.Database/autonomousDatabases/adb1
```

## API Providers

* `Oracle.Database`: 2025-03-01
