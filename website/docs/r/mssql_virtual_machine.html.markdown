---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_machine"
description: |-
  Manages a Microsoft SQL Virtual Machine
---

# azurerm_mssql_virtual_machine

Manages a Microsoft SQL Virtual Machine

## Example Usage

This example provisions a brief Managed MsSql Virtual Machine. The detailed example of the `azurerm_mssql_virtual_machine` resource can be found in [the `./examples/mssql/mssqlvm` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/mssql/mssqlvm)

```hcl
data "azurerm_virtual_machine" "example" {
  name                = "example-vm"
  resource_group_name = "example-resources"
}

resource "azurerm_mssql_virtual_machine" "example" {
  virtual_machine_id               = data.azurerm_virtual_machine.example.id
  sql_license_type                 = "PAYG"
  r_services_enabled               = true
  sql_connectivity_port            = 1433
  sql_connectivity_type            = "PRIVATE"
  sql_connectivity_update_password = "Password1234!"
  sql_connectivity_update_username = "sqllogin"

  auto_patching {
    day_of_week                            = "Sunday"
    maintenance_window_duration_in_minutes = 60
    maintenance_window_starting_hour       = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `virtual_machine_id` - (Required) The ID of the Virtual Machine. Changing this forces a new resource to be created.

* `sql_license_type` - (Optional) The SQL Server license type. Possible values are `AHUB` (Azure Hybrid Benefit) and `PAYG` (Pay-As-You-Go). Changing this forces a new resource to be created.

* `auto_backup` (Optional) An `auto_backup` block as defined below.

* `auto_patching` - (Optional) An `auto_patching` block as defined below.

* `key_vault_credential` - (Optional) (Optional) An `key_vault_credential` block as defined below.

* `r_services_enabled` - (Optional) Should R Services be enabled?

* `sql_connectivity_port` - (Optional) The SQL Server port. Defaults to `1433`.

* `sql_connectivity_type` - (Optional) The connectivity type used for this SQL Server. Defaults to `PRIVATE`.

* `sql_connectivity_update_password` - (Optional) The SQL Server sysadmin login password.

* `sql_connectivity_update_username` - (Optional) The SQL Server sysadmin login to create.

* `storage_configuration` - (Optional) An `storage_configuration` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `auto_backup` block supports the following:

* `encryption_enabled` - (Optional) Enable or disable encryption for backups. Defaults to `false`.

* `encryption_password` - (Optional) Encryption password to use. Must be specified when encryption is enabled.

* `manual_schedule` - (Optional) A `manual_schedule` block as documented below. When this block is present, the schedule type is set to `Manual`. Without this block, the schedule type is set to `Automated`.

* `retention_period_in_days` - (Required) Retention period of backups, in days. Valid values are from `1` to `30`.

* `storage_blob_endpoint` - (Required) Blob endpoint for the storage account where backups will be kept.

* `storage_account_access_key` - (Required) Access key for the storage account where backups will be kept.

* `system_databases_backup_enabled` - (Optional) Include or exclude system databases from auto backup. Defaults to `false`.

---

The `manual_schedule` block supports the following:

* `full_backup_frequency` - (Optional) Frequency of full backups. Valid values include `Daily` or `Weekly`. Required when `backup_schedule_automated` is false.

* `full_backup_start_hour` - (Optional) Start hour of a given day during which full backups can take place. Valid values are from `0` to `23`. Required when `backup_schedule_automated` is false.

* `full_backup_window_in_hours` - (Optional) Duration of the time window of a given day during which full backups can take place, in hours. Valid values are between `1` and `23`. Required when `backup_schedule_automated` is false.

* `log_backup_frequency_in_minutes` - (Optional) Frequency of log backups, in minutes. Valid values are from `5` to `60`. Required when `backup_schedule_automated` is false.

---

The `auto_patching` block supports the following:

* `day_of_week` - (Required) The day of week to apply the patch on.

* `maintenance_window_starting_hour` - (Required) The Hour, in the Virtual Machine Time-Zone when the patching maintenance window should begin.

* `maintenance_window_duration_in_minutes` - (Required) The size of the Maintenance Window in minutes.

---

The `key_vault_credential` block supports the following:

* `name` - (Required) The credential name.

* `key_vault_url` - (Required) The azure Key Vault url. Changing this forces a new resource to be created.

* `service_principal_name` - (Required) The service principal name to access key vault. Changing this forces a new resource to be created.

* `service_principal_secret` - (Required) The service principal name secret to access key vault. Changing this forces a new resource to be created.

---

The `storage_configuration` block supports the following:

* `disk_type` - (Required) The type of disk configuration to apply to the SQL Server. Valid values include `NEW`, `EXTEND`, or `ADD`.

* `storage_workload_type` - (Required) The type of storage workload. Valid values include `GENERAL`, `OLTP`, or `DW`.

* `data_settings` - (Optional) An `storage_settings` as defined below.

* `log_settings` - (Optional) An `storage_settings` as defined below.

* `temp_db_settings` - (Optional) An `storage_settings` as defined below.

---

The `storage_settings` block supports the following:

* `default_file_path` - (Required) The SQL Server default path

* `luns` - (Required) A list of Logical Unit Numbers for the disks. 

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the SQL Virtual Machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the MSSQL Virtual Machine.
* `update` - (Defaults to 60 minutes) Used when updating the MSSQL Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the MSSQL Virtual Machine.
* `delete` - (Defaults to 60 minutes) Used when deleting the MSSQL Virtual Machine.


## Import

Sql Virtual Machines can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_mssql_virtual_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/example1
```
