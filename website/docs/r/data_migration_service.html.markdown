---
subcategory: "Data Migration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_migration_service"
sidebar_current: "docs-azurerm-resource-data-migration-service"
description: |-
  Manage Azure Service instance.
---

# azurerm_data_migration_service

Manage Azure Data Migration Service instance.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "northeurope"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name 				   = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_data_migration_service" "example" {
	name                = "example-dms"
	location            = azurerm_resource_group.example.location
	resource_group_name = azurerm_resource_group.example.name
	virtual_subnet_id   = azurerm_subnet.example.id
	sku_name            = "Standard_1vCores"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The Azure region of the operation. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the resource group. Changing this forces a new resource to be created.

* `name` - (Required) Name of the service. Changing this forces a new resource to be created.

* `virtual_subnet_id` - (Required) The ID of the virtual subnet resource to which the service should be joined. Changing this forces a new resource to be created.

* `sku_name` - (Required) The sku name of the service. Changing this forces a new resource to be created.

* `kind` - (Optional) The resource kind. Only 'Cloud' (the default) is supported. Changing this forces a new resource to be created.

* `delete_running_tasks` - (Optional) When destroy the resource, delete it even if it contains running tasks.

* `tags` - (Optional) Resource tags.

## Attributes Reference

The following attributes are exported:

* `type` - (Optional) The resource type chain (e.g. virtualMachines/extensions)

* `provisioning_state` - The resource's provisioning state

* `id` - Resource ID.

## Import

Data Migration Service can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_data_migration_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.DataMigration/services/data_migration_service1
```
