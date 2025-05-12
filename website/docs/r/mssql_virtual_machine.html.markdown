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

This example provisions a brief Managed Microsoft SQL Virtual Machine. The detailed example of the `azurerm_mssql_virtual_machine` resource can be found in [the `./examples/mssql/mssqlvm` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/mssql/mssqlvm)

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

* `sql_license_type` - (Optional) The SQL Server license type. Possible values are `AHUB` (Azure Hybrid Benefit), `DR` (Disaster Recovery), and `PAYG` (Pay-As-You-Go). Changing this forces a new resource to be created.

* `auto_backup` - (Optional) An `auto_backup` block as defined below. This block can be added to an existing resource, but removing this block forces a new resource to be created.

* `auto_patching` - (Optional) An `auto_patching` block as defined below.

* `key_vault_credential` - (Optional) An `key_vault_credential` block as defined below.

* `r_services_enabled` - (Optional) Should R Services be enabled?

* `sql_connectivity_port` - (Optional) The SQL Server port. Defaults to `1433`.

* `sql_connectivity_type` - (Optional) The connectivity type used for this SQL Server. Possible values are `LOCAL`, `PRIVATE` and `PUBLIC`. Defaults to `PRIVATE`.

* `sql_connectivity_update_password` - (Optional) The SQL Server sysadmin login password.

* `sql_connectivity_update_username` - (Optional) The SQL Server sysadmin login to create.

* `sql_instance` - (Optional) A `sql_instance` block as defined below.

* `storage_configuration` - (Optional) An `storage_configuration` block as defined below.

* `assessment` - (Optional) An `assessment` block as defined below.

* `sql_virtual_machine_group_id` - (Optional) The ID of the SQL Virtual Machine Group that the SQL Virtual Machine belongs to.

* `wsfc_domain_credential` - (Optional) A `wsfc_domain_credential` block as defined below

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `auto_backup` block supports the following:


* `encryption_password` - (Optional) Encryption password to use. Setting a password will enable encryption.

* `manual_schedule` - (Optional) A `manual_schedule` block as documented below. When this block is present, the schedule type is set to `Manual`. Without this block, the schedule type is set to `Automated`.

* `retention_period_in_days` - (Required) Retention period of backups, in days. Valid values are from `1` to `30`.

* `storage_blob_endpoint` - (Required) Blob endpoint for the storage account where backups will be kept.

* `storage_account_access_key` - (Required) Access key for the storage account where backups will be kept.

* `system_databases_backup_enabled` - (Optional) Include or exclude system databases from auto backup.

---

The `manual_schedule` block supports the following:

* `full_backup_frequency` - (Required) Frequency of full backups. Valid values include `Daily` or `Weekly`.

* `full_backup_start_hour` - (Required) Start hour of a given day during which full backups can take place. Valid values are from `0` to `23`.

* `full_backup_window_in_hours` - (Required) Duration of the time window of a given day during which full backups can take place, in hours. Valid values are between `1` and `23`.

* `log_backup_frequency_in_minutes` - (Required) Frequency of log backups, in minutes. Valid values are from `5` to `60`.

* `days_of_week` - (Optional) A list of days on which backup can take place. Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`

~> **Note:** `days_of_week` can only be specified when `manual_schedule` is set to `Weekly`

---

The `auto_patching` block supports the following:

* `day_of_week` - (Required) The day of week to apply the patch on. Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`.

* `maintenance_window_starting_hour` - (Required) The Hour, in the Virtual Machine Time-Zone when the patching maintenance window should begin.

* `maintenance_window_duration_in_minutes` - (Required) The size of the Maintenance Window in minutes.

---

The `key_vault_credential` block supports the following:

* `name` - (Required) The credential name.

* `key_vault_url` - (Required) The Azure Key Vault url. Changing this forces a new resource to be created.

* `service_principal_name` - (Required) The service principal name to access key vault. Changing this forces a new resource to be created.

* `service_principal_secret` - (Required) The service principal name secret to access key vault. Changing this forces a new resource to be created.

---

The `sql_instance` block supports the following:

* `adhoc_workloads_optimization_enabled` - (Optional) Specifies if the SQL Server is optimized for adhoc workloads. Possible values are `true` and `false`. Defaults to `false`.

* `collation` - (Optional) Collation of the SQL Server. Defaults to `SQL_Latin1_General_CP1_CI_AS`. Changing this forces a new resource to be created.

* `instant_file_initialization_enabled` - (Optional) Specifies if Instant File Initialization is enabled for the SQL Server. Possible values are `true` and `false`. Defaults to `false`. Changing this forces a new resource to be created.

* `lock_pages_in_memory_enabled` - (Optional) Specifies if Lock Pages in Memory is enabled for the SQL Server. Possible values are `true` and `false`. Defaults to `false`. Changing this forces a new resource to be created.

* `max_dop` - (Optional) Maximum Degree of Parallelism of the SQL Server. Possible values are between `0` and `32767`. Defaults to `0`.

* `max_server_memory_mb` - (Optional) Maximum amount memory that SQL Server Memory Manager can allocate to the SQL Server process. Possible values are between `128` and `2147483647` Defaults to `2147483647`.

* `min_server_memory_mb` - (Optional) Minimum amount memory that SQL Server Memory Manager can allocate to the SQL Server process. Possible values are between `0` and `2147483647` Defaults to `0`.

~> **Note:** `max_server_memory_mb` must be greater than or equal to `min_server_memory_mb`

---

The `storage_configuration` block supports the following:

* `disk_type` - (Required) The type of disk configuration to apply to the SQL Server. Valid values include `NEW`, `EXTEND`, or `ADD`.

* `storage_workload_type` - (Required) The type of storage workload. Valid values include `GENERAL`, `OLTP`, or `DW`.

* `data_settings` - (Optional) A `storage_settings` block as defined below.

* `log_settings` - (Optional) A `storage_settings` block as defined below.

* `system_db_on_data_disk_enabled` - (Optional) Specifies whether to set system databases (except tempDb) location to newly created data storage. Possible values are `true` and `false`. Defaults to `false`.

* `temp_db_settings` - (Optional) An `temp_db_settings` block as defined below.

---

The `storage_settings` block supports the following:

* `default_file_path` - (Required) The SQL Server default path

* `luns` - (Required) A list of Logical Unit Numbers for the disks.

---

The `temp_db_settings` block supports the following:

* `default_file_path` - (Required) The SQL Server default path

* `luns` - (Required) A list of Logical Unit Numbers for the disks.

* `data_file_count` - (Optional) The SQL Server default file count. This value defaults to `8`

* `data_file_size_mb` - (Optional) The SQL Server default file size - This value defaults to `256`

* `data_file_growth_in_mb` - (Optional) The SQL Server default file size - This value defaults to `512`

* `log_file_size_mb` - (Optional) The SQL Server default file size - This value defaults to `256`

* `log_file_growth_mb` - (Optional) The SQL Server default file size - This value defaults to `512`

---

The `assessment` block supports the following:

* `enabled` - (Optional) Should Assessment be enabled? Defaults to `true`.

* `run_immediately` - (Optional) Should Assessment be run immediately? Defaults to `false`.

* `schedule` - (Optional) An `schedule` block as defined below.

---

The `schedule` block supports the following:

* `weekly_interval` - (Optional) How many weeks between assessment runs. Valid values are between `1` and `6`.

* `monthly_occurrence` - (Optional) How many months between assessment runs. Valid values are between `1` and `5`.

~> **Note:** Either one of `weekly_interval` or `monthly_occurrence` must be specified.

* `day_of_week` - (Required) What day of the week the assessment will be run. Possible values are `Friday`, `Monday`, `Saturday`, `Sunday`, `Thursday`, `Tuesday` and `Wednesday`.

* `start_time` - (Required) What time the assessment will be run. Must be in the format `HH:mm`.

---

The `wsfc_domain_credential` block supports the following:

* `cluster_bootstrap_account_password` - (Required) The account password used for creating cluster.

* `cluster_operator_account_password` - (Required) The account password used for operating cluster.

* `sql_service_account_password` - (Required) The account password under which SQL service will run on all participating SQL virtual machines in the cluster.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SQL Virtual Machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Microsoft SQL Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft SQL Virtual Machine.
* `update` - (Defaults to 1 hour) Used when updating the Microsoft SQL Virtual Machine.
* `delete` - (Defaults to 1 hour) Used when deleting the Microsoft SQL Virtual Machine.

## Import

Microsoft SQL Virtual Machines can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_virtual_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/example1
```
