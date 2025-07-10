---
subcategory: "Qumulo"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_qumulo_file_system"
description: |-
  Manages an Azure Native Qumulo Scalable File System.
---

# azurerm_qumulo_file_system

Manages an Azure Native Qumulo Scalable File System.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
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

  delegation {
    name = "delegation"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Qumulo.Storage/fileSystems"
    }
  }
}

resource "azurerm_qumulo_file_system" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  admin_password      = ")^X#ZX#JRyIY}t9"
  availability_zone   = "1"
  delegated_subnet_id = azurerm_subnet.example.id
  storage_sku         = "Standard"
  email               = "test@test.com"
  tags = {
    environment = "test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Azure Native Qumulo Scalable File System resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Azure Native Qumulo Scalable File System should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Native Qumulo Scalable File System should exist. Changing this forces a new resource to be created.

* `admin_password` - (Required) The initial administrator password of the Azure Native Qumulo Scalable File System. Changing this forces a new resource to be created.

* `email` - (Required) The email address used for the Azure Native Qumulo Scalable File System. Changing this forces a new resource to be created.

* `storage_sku` - (Required) The storage Sku. Possible values are `Cold_LRS`, `Hot_LRS` and `Hot_ZRS`. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The delegated subnet ID for Vnet injection. Changing this forces a new resource to be created.

* `zone` - (Required) The Availability Zone in which the Azure Native Qumulo Scalable File system is located. Changing this forces a new resource to be created.

* `offer_id` - (Optional) Specifies the marketplace offer ID. Defaults to `qumulo-saas-mpp`. Changing this forces a new resource to be created.

* `plan_id` - (Optional) Specifies the marketplace plan ID. Defaults to `azure-native-qumulo-v3`. Changing this forces a new resource to be created.

* `publisher_id` - (Optional) Specifies the marketplace publisher ID. Defaults to `qumulo1584033880660`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the File System.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the File System.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Azure Native Qumulo Scalable File System.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Native Qumulo Scalable File System.
* `update` - (Defaults to 1 hour) Used when updating the Azure Native Qumulo Scalable File System.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Native Qumulo Scalable File System.

## Import

An existing File System can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_qumulo_file_system.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Qumulo.Storage/fileSystems/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Qumulo.Storage`: 2024-06-19
