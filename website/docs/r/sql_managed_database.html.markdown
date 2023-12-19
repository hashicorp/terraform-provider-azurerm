---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_managed_database"
description: |-
  Manages a SQL Azure Managed Database.
---

# azurerm_sql_managed_database

Manages a SQL Azure Managed Database.

!>**IMPORTANT:** To mitigate the possibility of accidental data loss it is highly recommended that you use the `prevent_destroy` lifecycle argument in your configuration file for this resource. For more information on the `prevent_destroy` lifecycle argument please see the [terraform documentation](https://developer.hashicorp.com/terraform/tutorials/state/resource-lifecycle#prevent-resource-deletion).

->**NOTE:** The `azurerm_sql_managed_database` resource is deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use the [`azurerm_mssql_managed_database`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_managed_database) resource instead.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "test-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "acctsub"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_sql_managed_instance" "example" {
  name                         = "example"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
  license_type                 = "BasePrice"
  subnet_id                    = azurerm_subnet.example.id
  sku_name                     = "GP_Gen5"
  vcores                       = 4
  storage_size_in_gb           = 32
}

resource "azurerm_sql_managed_database" "example" {
  sql_managed_instance_id = azurerm_sql_managed_instance.example.id
  name                    = "exampledatabase"
  location                = azurerm_resource_group.example.location

  # prevent the possibility of accidental data loss
  lifecycle {
    prevent_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL Managed Instance. Changing this forces a new resource to be created.

* `sql_managed_instance_id` - (Required) The SQL Managed Instance ID that this Managed Database will be associated with. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

---

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The SQL Managed Database ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Sql Managed Database.
* `create` - (Defaults to 24 hours) Used when creating the Sql Managed Database.
* `delete` - (Defaults to 24 hours) Used when deleting the Sql Managed Database.

## Import

SQL Managed Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_managed_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/managedInstances/myserver/databases/mydatabase
```
