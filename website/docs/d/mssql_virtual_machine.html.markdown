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

* `resource_group_name` - (Required) The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal. Changing this forces a new resource to be created.

* `location` - (Required) The resource location. The change of the location forces a new resource to be created.

* `virtual_machine_resource_id` - The ARM Resource id of underlying virtual machine created from SQL marketplace image.

* `sql_license_type` - The SQL Server license type. Possible values include: 'PAYG'(Pay As You Go), 'AHUB'(Azure Hybrid Benefit).Defaults to `PAYG`.

* `sql_sku` - The SQL Server edition type. Possible values include: 'Developer', 'Express', 'Standard', 'Enterprise', 'Web'. Defaults to `Developer`.

* `tags` - Resource tags. Changing this forces a new resource to be created.

