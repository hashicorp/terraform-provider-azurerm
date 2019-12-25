subcategory: "MSSQLVM"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_machine"
sidebar_current: "docs-azurerm-datasource-mssql-virtual-machine"
description: |-
  Gets information about an existing MsSql Virtual Machine
---

# Data Source: azurerm_mssql_virtual_machine

Use this data source to access information about an existing MS Sql Virtual Machine.


## Example Usage

```hcl
data "azurerm_sql_virtual_machine" "example" {
  resource_group           = "example-resource-group"
  name                     = "example-sql-virtual-machine"
}

output "sql_virtual_machine_id" {
  value = "${data.azurerm_sql_virtual_machine.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `resource_group` - (Required) Name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal.

* `name` - (Required) Name of the SQL virtual machine.


## Attributes Reference

The following attributes are exported:

* `resource_group` - (Required) Name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal. Changing this forces a new resource to be created.

* `location` - (Required) Resource location. Changing this forces a new resource to be created.

* `virtual_machine_resource_id` - (Required) ARM Resource id of underlying virtual machine created from SQL marketplace image.

* `sql_virtual_machine_group_resource_id` - (Optional) ARM resource id of the SQL virtual machine group this SQL virtual machine is or will be part of.

* `name` - (Computed) Name of the SQL virtual machine, which is the same with the name of the Virtual Machine provided.

* `sql_server_license_type` - (Optional) SQL Server license type. Defaults to `PAYG`.Possible values include: 'PAYG'(Pay As You Go), 'AHUB'(Azure Hybrid Benefit).

* `sql_image_sku` - (Optional) SQL Server edition type. Defaults to `Developer`.Possible values include: 'Developer', 'Express', 'Standard', 'Enterprise', 'Web'.

* `auto_patching_settings` - (Optional) One `auto_patching_setting` block defined below.SQL Server Azure VMs can use Automated Patching to schedule a maintenance window for installing important windows and SQL Server updates automatically. Please refer [automated-patching](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/sql/virtual-machines-windows-sql-automated-patching) for more information.

* `key_vault_credential_settings` - (Optional) One `key_vault_credential_setting` block defined below. For more information, please refer to [virtual machines windows sql keyvault](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/sql/virtual-machines-windows-ps-sql-keyvault)

* `server_configurations_management_settings` - (Optional) One `server_configurations_management_setting` block defined below.

* `storage_configuration_settings` - (Optional) One `storage_configuration_setting` block defined below.Customize performance, size, and workload type to optimize storage for this virtual machine. For optimal performance, separate drives will be created for data and log storage by default. [Learn more about SQL Server best performance practices](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/sql/virtual-machines-windows-sql-performance).

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.


The `auto_patching_setting` block supports the following:

* `enable` - (Optional) Enable or disable autopatching on SQL virtual machine.

* `day_of_week` - (Optional) Day of week to apply the patch on. Defaults to `Monday`.

* `maintenance_window_starting_hour` - (Optional) Hour of the day when patching is initiated. Local VM time.

* `maintenance_window_duration` - (Optional) Duration of patching.


---

The `key_vault_credential_setting` block supports the following:

* `enable` - (Optional) Enable or disable key vault credential setting.

* `credential_name` - (Optional) Credential name.

* `azure_key_vault_url` - (Optional) Azure Key Vault url.

* `service_principal_name` - (Optional) Service principal name to access key vault.

* `service_principal_secret` - (Optional) Service principal name secret to access key vault.

---

The `server_configurations_management_setting` block supports the following:

* `sql_connectivity_type` - (Optional) SQL Server connectivity option. Defaults to `LOCAL`.

* `sql_connectivity_port` - (Optional) SQL Server port.

* `sql_connectivity_auth_update_user_name` - (Optional) SQL Server sysadmin login to create.

* `sql_connectivity_auth_update_password` - (Optional) SQL Server sysadmin login password.

* `is_r_services_enabled` - (Optional) Enable or disable R services (SQL 2016 onwards).Enables SQL Server Machine Learning Services (In-Database), allowing you to utilize advanced analytics within your SQL Server. SQL Server Machine Learning Services (In-Database) is only supported with SQL Server 2017 Enterprise.

---

The `storage_configuration_setting` block supports the following:

* `storage_workload_type` - (Optional) Storage workload type. Defaults to `GENERAL`.Possible values include: 'GENERAL', 'OLTP'(Transactional processing), 'DW'(Data warehousing).

* `sql_data_luns` - (Optional) Logical Unit Numbers for the disks.

* `sql_data_default_file_path` - (Optional) SQL Server default file path

* `sql_log_luns` - (Optional) Logical Unit Numbers for the disks.

* `sql_log_default_file_path` - (Optional) SQL Server default file path

* `sql_temp_db_luns` - (Optional) Logical Unit Numbers for the disks.

* `sql_temp_db_default_file_path` - (Optional) SQL Server default file path



