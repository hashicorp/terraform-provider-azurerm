---
subcategory: "Qumulo"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_qumulo_file_system"
description: |-
  Manages a File System.
---

# azurerm_qumulo_file_system

Manages a Qumulo File System.

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
  initial_capacity    = 21
  marketplace_plan_id = "qumulo-on-azure-v1%%gmz7xq9ge3py%%P1M"
  storage_sku         = "Standard"
  user_email_address  = "test@test.com"
  tags = {
    environment = "test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Qumulo File System resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Qumulo File System should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Qumulo File System should exist. Changing this forces a new resource to be created.
 
* `admin_password` - (Required) Initial administrator password of the Qumulo File System. Changing this forces a new resource to be created.

* `delegated_subnet_id` - (Required) Delegated subnet ID for Vnet injection. Changing this forces a new resource to be created.

* `initial_capacity` - (Required) Storage capacity in TB. Changing this forces a new resource to be created.

* `marketplace_plan_id` - (Required) Specifies the marketplace plan ID. Changing this forces a new resource to be created.

* `storage_sku` - (Required) Storage Sku. Possible values are `Performance` and `Standard`. Changing this forces a new resource to be created.

* `user_email_address` - (Required) The email address. Changing this forces a new resource to be created.

* `availability_zone` - (Optional) Availability zone. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the File System.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the File System.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 180 minutes) Used when creating this File System.
* `delete` - (Defaults to 120 minutes) Used when deleting this File System.
* `read` - (Defaults to 5 minutes) Used when retrieving this File System.
* `update` - (Defaults to 120 minutes) Used when updating this File System.

## Import

An existing File System can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_qumulo_file_system.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Qumulo.Storage/fileSystems/example
```
