---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_clone_from_backup"
description: |-
  Manages a autonomous database clone from backup.
---

# azurerm_oracle_autonomous_database_clone_from_backup

Manages a autonomous database clone from backup.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_oracle_autonomous_database" "example" {
  name                             = "example"
  resource_group_name              = "example"
  location                         = "West Europe"
  subnet_id                        = "example"
  display_name                     = "example"
  db_workload                      = "example"
  mtls_connection_required         = false
  backup_retention_period_in_days  = 42
  compute_model                    = "example"
  data_storage_size_in_gbs         = 42
  auto_scaling_for_storage_enabled = false
  virtual_network_id               = "example"
  admin_password                   = "example"
  auto_scaling_enabled             = "example"
  character_set                    = "example"
  compute_count                    = 1.23456
  national_character_set           = "example"
  license_model                    = false
  db_version                       = "example"
  allowedIps                       = ""
}

resource "azurerm_oracle_autonomous_database_clone_from_backup" "example" {
  name                             = "example"
  resource_group_name              = azurerm_oracle_autonomous_database.example.resource_group_name
  location                         = "West Europe"
  source                           = "BackupFromTimestamp"
  mtls_connection_required         = false
  compute_model                    = "ECPU"
  clone_type                       = "Full"
  db_version                       = "19c"
  auto_scaling_enabled             = false
  backup_retention_period_in_days  = 42
  license_model                    = "LicenseIncluded"
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  source_autonomous_database_id    = azurerm_oracle_autonomous_database.example.id
  display_name                     = "example"
  auto_scaling_for_storage_enabled = false
  admin_password                   = "BEstrO0ng_#11"
  db_workload                      = "OLTP"
  data_storage_size_in_tb          = 42
  national_character_set           = "AL16UTF16"
  allowedIps                       = ""
}
```

## Arguments Reference

The following arguments are supported:

* `admin_password` - (Required) The password for the SYS, SYSTEM, and PDB Admin users. The password must be at least 12 characters long, and contain at least 1 uppercase, 1 lowercase, and 1 numeric character. It cannot contain the double quote symbol (") or the username "admin," regardless of casing.

* `auto_scaling_enabled` - (Required) Indicates if auto-scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - (Required) ndicates if auto-scaling is enabled for the Autonomous Database storage.

* `backup_retention_period_in_days` - (Required) Retention period, in days, for backups. Possible values are between 1 and 60. Changing these forces a new Autonomous Database Clone to be created.

* `character_set` - (Required) The character set for the autonomous database. Changing this forces a new Autonomous Database Clone to be created

* `clone_type` - (Required) The type of clone to create. Possible values are Full and Metadata. Changing this forces a new autonomous database clone from backup to be created.

* `compute_count` - (Required) The compute amount (CPUs) available to the database. Possible values are between 2.0 and 512.0.

* `compute_model` - (Required) he compute model of the Autonomous Database. Changing this forces a new autonomous database clone from backup to be created.

* `data_storage_size_in_tb` - (Required) The maximum storage that can be allocated for the database, in terabytes. Possible values are between 1 and 384.

* `db_version` - (Required) A valid Oracle Database version for Autonomous Database. Changing this forces a new autonomous database clone from backup to be created.

* `db_workload` - (Required) The Autonomous Database workload type. Possible values are OLTP, DW, APEX, and AJD. Changing this forces a new autonomous database clone from backup to be created.

    * OLTP: Indicates an Autonomous Transaction Processing database. 
    * DW: Indicates an Autonomous Data Warehouse database. 
    * AJD: Indicates an Autonomous JSON Database. 
    * APEX: Indicates an Autonomous Database with the Oracle APEX Application Development workload type.

* `display_name` - (Required) he user-friendly name for the Autonomous Database. Changing this forces a new autonomous database clone from backup to be created.

* `license_model` - (Required) he Oracle license model that applies to the Oracle Autonomous Database. Possible values are LicenseIncluded and BringYourOwnLicense. Changing this forces a new autonomous database clone from backup to be created.

* `location` - (Required) The Azure Region where the autonomous database clone from backup should exist. Changing this forces a new autonomous database clone from backup to be created.

* `mtls_connection_required` - (Required) Specifies if the Autonomous Database requires mTLS connections. Changing this forces a new autonomous database clone from backup to be created.

* `name` - (Required) The name which should be used for this autonomous database clone from backup. Changing this forces a new autonomous database clone from backup to be created.

* `national_character_set` - (Required) The national character set for the autonomous database. Changing this forces a new autonomous database clone from backup to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the autonomous database clone from backup should exist. Changing this forces a new autonomous database clone from backup to be created.

* `source` - (Required) The source of the clone. Possible value is BackupFromTimestamp. Changing this forces a new autonomous database clone from backup to be created.

* `source_autonomous_database_id` - (Required) The ID of the source Autonomous Database to clone from. Changing this forces a new autonomous database clone from backup to be created.

---

* `allowed_ips` - (Optional) Defines the network access type for the Autonomous Database. If the property is explicitly set to an empty list, it allows secure public access to the database from any IP address. If specific ACL (Access Control List) values are provided, access will be restricted to only the specified IP addresses.

* `backup_timestamp` - (Optional) The autonomous database backup time stamp to be used for a cloning autonomous database. Changing this forces a new autonomous database clone from backup to be created.

* `customer_contacts` - (Optional) Customer contact email addresses. Changing this forces a new autonomous database clone from backup to be created.

* `subnet_id` - (Optional) The ID of the subnet the resource is associated with. Changing this forces a new autonomous database clone from backup to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the autonomous database clone from backup.

* `use_latest_available_backup_time_stamp` - (Optional) setting this value to true will initiate cloning from latest backup time stamp. Changing this forces a new autonomous database clone from backup to be created.

* `virtual_network_id` - (Optional) The ID of the Virtual Network this Autonomous Database Clone should be created in. Changing this forces a new autonomous database clone from backup to be created.

## Attributes Reference

In addition to the Arguments listed above—the following Attributes are exported: 

* `id` - The ID of the autonomous database clone from backup.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the autonomous database clone from backup.
* `read` - (Defaults to 5 minutes) Used when retrieving the autonomous database clone from backup.
* `update` - (Defaults to 30 minutes) Used when updating the autonomous database clone from backup.
* `delete` - (Defaults to 30 minutes) Used when deleting the autonomous database clone from backup.

## Import

autonomous database clone from backups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_autonomous_database_clone_from_backup.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Oracle.Database/autonomousDatabases/example
```
