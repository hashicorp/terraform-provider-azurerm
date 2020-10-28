---
subcategory: "Database Migration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_database_migration_service"
description: |-
  Manage a Azure Database Migration Service.
---

# azurerm_database_migration_service

Manages a Azure Database Migration Service.

~> **NOTE:** Destroying a Database Migration Service will leave any outstanding tasks untouched. This is to avoid unexpectedly deleting any tasks managed outside of terraform.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

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
  name               = "example-subnet"
  virtual_network_id = azurerm_virtual_network.example.id
  address_prefixes   = ["10.0.1.0/24"]
}

resource "azurerm_database_migration_service" "example" {
  name                = "example-dms"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example.id
  sku_name            = "Standard_1vCores"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specify the name of the database migration service. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the resource group in which to create the database migration service. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the virtual subnet resource to which the database migration service should be joined. Changing this forces a new resource to be created.

* `sku_name` - (Required) The sku name of the database migration service. Possible values are `Premium_4vCores`, `Standard_1vCores`, `Standard_2vCores` and `Standard_4vCores`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assigned to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Database Migration Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API.

## Import

Database Migration Services can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_database_migration_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.DataMigration/services/database_migration_service1
```
