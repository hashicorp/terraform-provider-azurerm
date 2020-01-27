---
subcategory: "Data Migration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_migration_service"
description: |-
  Manage Azure Service instance.
---

# azurerm_data_migration_service

Manages a Azure Data Migration Service.

~> **NOTE on destroy behavior of Data Migration Service:** Destroy a Data Migration Service will leave any outstanding tasks untouched. This is to avoid unexpectedly delete any tasks managed out of terraform.

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
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_data_migration_service" "example" {
  name                = "example-dms"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example.id
  sku_name            = "Standard_1vCores"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specify the name of the data migration service. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the resource group in which to create the data migration service. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the virtual subnet resource to which the data migration service should be joined. Changing this forces a new resource to be created.

* `sku_name` - (Required) The sku name of the data migration service. Possible values are `Premium_4vCores`, `Standard_1vCores`, `Standard_2vCores` and `Standard_4vCores`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assigned to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Data Migration Service.

## Import

Data Migration Services can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_data_migration_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.DataMigration/services/data_migration_service1
```
