---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_database"
description: |-
  Manages an Azure SQL Azure Managed Database.
---

# azurerm_mssql_managed_database

Manages an Azure SQL Azure Managed Database for a SQL Managed Instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_mssql_managed_instance" "example" {
  name                = "managedsqlinstance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.example.id
  vcores             = 4

  administrator_login          = "msadministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_managed_database" "example" {
  name                = "example"
  managed_instance_id = azurerm_mssql_managed_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Database to create. Changing this forces a new resource to be created.

* `managed_instance_id` - (Required) The ID of the Azure SQL Managed Instance on which to create this Managed Database. Changing this forces a new resource to be created.

* `long_term_retention_policy` - (Optional) A `long_term_retention_policy` block as defined below.

* `short_term_retention_days` - (Optional) The backup retention period in days. This is how many days Point-in-Time Restore will be supported.

---

A `long_term_retention_policy` block supports the following:

* `weekly_retention` - (Optional) The weekly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 520 weeks. e.g. `P1Y`, `P1M`, `P1W` or `P7D`.
* `monthly_retention` - (Optional) The monthly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 120 months. e.g. `P1Y`, `P1M`, `P4W` or `P30D`.
* `yearly_retention` - (Optional) The yearly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 10 years. e.g. `P1Y`, `P12M`, `P52W` or `P365D`.
* `week_of_year` - (Optional) The week of year to take the yearly backup. Value has to be between `1` and `52`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Azure SQL Managed Database ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mssql Managed Database.
* `create` - (Defaults to 30 minutes) Used when creating the Mssql Managed Database.
* `update` - (Defaults to 30 minutes) Used when updating the Mssql Managed Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Mssql Managed Database.

## Import

SQL Managed Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/managedInstances/myserver/databases/mydatabase
```
