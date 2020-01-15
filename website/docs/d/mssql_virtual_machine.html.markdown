---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_machine"
description: |-
  Gets information about an existing MsSql Virtual Machine
---

# Data Source: azurerm_mssql_virtual_machine

Use this data source to access information about an existing MS Sql Virtual Machine.


## Example Usage

```hcl
data "azurerm_mssql_virtual_machine" "example" {
  resource_group_name = "example-resource-group-name"
  name                = "example-sql-virtual-machine"
}

output "mssql_virtual_machine_id" {
  value = "${data.azurerm_mssql_virtual_machine.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL virtual machine.

* `resource_group_name` - (Required) The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal.


## Attributes Reference

The following attributes are exported:

* `id` - The ARM resource id of the SQL virtual machine group this SQL virtual machine is or will be part of.

* `name` - The name of the SQL virtual machine, which is the same with the name of the Virtual Machine provided.

* `resource_group_name` - The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal. Changing this forces a new resource to be created.

* `location` - The resource location. The change of the location forces a new resource to be created.

* `virtual_machine_resource_id` - The ARM Resource id of underlying virtual machine created from SQL marketplace image.

* `sql_license_type` - The SQL Server license type. Possible values include: 'PAYG'(Pay As You Go), 'AHUB'(Azure Hybrid Benefit).Defaults to `PAYG`.

* `sql_sku` - The SQL Server edition type. Possible values include: 'Developer', 'Express', 'Standard', 'Enterprise', 'Web'. Defaults to `Developer`.

* `auto_patching` -  The `auto_patching_setting` block defined below.SQL Server Azure VMs can use Automated Patching to schedule a maintenance window for installing important windows and SQL Server updates automatically. Please refer [automated patching](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/sql/virtual-machines-windows-sql-automated-patching) for more information.

* `key_vault_credential` -  The `key_vault_credential_setting` block defined below. For more information, please refer to [virtual machines windows sql keyvault](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/sql/virtual-machines-windows-ps-sql-keyvault)

* `server_configuration` -  The `server_configurations_management_setting` block defined below.

* `tags` - Resource tags. Changing this forces a new resource to be created.

The `auto_patching_setting` block supports the following:

* `enable` -  Enable or disable autopatching on SQL virtual machine.

* `day_of_week` -  The day of week to apply the patch on. Defaults to `Monday`.

* `maintenance_window_starting_hour` -  The hour of the day when patching is initiated. Local VM time.

* `maintenance_window_duration_in_minutes` -  The duration of patching.

---

The `key_vault_credential_setting` block supports the following:

* `enable` -  Enable or disable key vault credential setting.

* `credential_name` -  The credential name.

* `azure_key_vault_url` -  The azure Key Vault url.

* `service_principal_name` -  The service principal name to access key vault.

---

The `server_configurations_management_setting` block supports the following:

* `sql_connectivity_type` -  The SQL Server connectivity option. Defaults to `LOCAL`.

* `sql_connectivity_port` -  The SQL Server port.

* `is_r_services_enabled` - Enable or disable R services (SQL 2016 onwards).Enables SQL Server Machine Learning Services (In-Database), allowing you to utilize advanced analytics within your SQL Server. SQL Server Machine Learning Services (In-Database) is only supported with SQL Server 2017 Enterprise.

