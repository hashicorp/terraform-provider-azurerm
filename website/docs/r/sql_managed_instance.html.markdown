---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_managed_instance"
description: |-
  Manages a SQL Azure Managed Instance.
---

# azurerm_sql_managed_instance

Manages a SQL Azure Managed Instance.

~> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "database-rg"
  location = "West Europe"
}

resource "azurerm_sql_managed_instance" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
  license_type                 = "BasePrice"
  subnet_id                    = "${azurerm_subnet.test.id}"
  sku_name                     = "GP_Gen5"
  vcores                       = 4
  storage_size_in_gb           = 32

  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_subnet_route_table_association.test,
  ]
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL Managed Instance. This needs to be globally unique within Azure.

* `resource_group_name` - (Required) The name of the resource group in which to create the SQL Server.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for the SQL Managed Instance. Valid values include `GP_Gen4`, `GP_Gen5`, `BC_Gen4`, `BC_Gen5`. Changing this forces a new resource to be created.

* `vcores` - (Required) Number of cores that should be assigned to your instance. Values can be `8`, `16`, or `24` if `sku_name` is `GP_Gen4`, or `8`, `16`, `24`, `32`, or `40` if `sku_name` is `GP_Gen5`.

* `storage_size_in_gb` - (Required) Maximum storage space for your instance. It should be a multiple of 32GB.

* `license_type` - (Required) What type of license the Managed Instance will use. Valid values include can be `PriceIncluded` or `BasePrice`.

* `administrator_login` - (Required) The administrator login name for the new server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx)

* `subnet_id` - (Required) The subnet resource id that the SQL Managed Instance will be associated with.

* `collation` - (Optional) Specifies how the SQL Managed Instance will be collated. Default value is `SQL_Latin1_General_CP1_CI_AS`. Changing this forces a new resource to be created.

* `public_data_endpoint_enabled` - (Optional) Is the public data endpoint enabled? Default value is `false`. 

* `minimum_tls_version` - (Optional) The Minimum TLS Version. Default value is `1.2` Valid values include `1.0`, `1.1`, `1.2`. 

* `proxy_override` - (Optional) Specified how the SQL Managed Instance will be accessed. Default value is `Default`. Valid values include `Default`, `Proxy`, and `Redirect`. 

* `timezone_id` - (Optional) The TimeZone ID that the SQL Managed Instance will be operating in. Default value is `UTC`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) Sku of the managed instance. Values can be GP_Gen4 or GP_Gen5.

## Attributes Reference

The following attributes are exported:

* `id` - The SQL Managed Instance ID.
* `fqdn` - The fully qualified domain name of the Azure Managed SQL Instance

## Import

SQL Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_managed_instance.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/managedInstances/myserver
```