---
subcategory: "Database Migration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_database_migration_project"
description: |-
  Manage Azure Database Migration Project instance.
---

# azurerm_database_migration_project

Manage a Azure Database Migration Project.

~> **NOTE:** Destroying a Database Migration Project will leave any outstanding tasks untouched. This is to avoid unexpectedly deleting any tasks managed outside of terraform.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_database_migration_service" "example" {
  name                = "example-dbms"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  virtual_subnet_id   = azurerm_subnet.example.id
  sku_name            = "Standard_1vCores"
}

resource "azurerm_database_migration_project" "example" {
  name                = "example-dbms-project"
  service_name        = azurerm_database_migration_service.example.name
  resource_group_name = azurerm_resource_group.example.name
  location            = zurerm_resource_group.example.location
  source_platform     = "SQL"
  target_platform     = "SQLDB"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specify the name of the database migration project. Changing this forces a new resource to be created.

* `service_name` - (Required) Name of the database migration service where resource belongs to. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the resource group in which to create the database migration project. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `source_platform` - (Required) The platform type of the migration source. Currently only support: `SQL`(on-premises SQL Server). Changing this forces a new resource to be created.

* `target_platform` - (Required) The platform type of the migration target. Currently only support: `SQLDB`(Azure SQL Database). Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assigned to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Database Migration Project.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API.

## Import

Database Migration Projects can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_database_migration_project.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.DataMigration/services/example-dms/projects/project1
```
