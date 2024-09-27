---
subcategory: "Service Networking"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_load_balancer_subnet_association"
description: |-
  Manages an association between an Application Gateway for Containers and a Subnet.
---

# azurerm_application_load_balancer_subnet_association

Manages an association between an Application Gateway for Containers and a Subnet.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "westeurope"
}

resource "azurerm_application_load_balancer" "example" {
  name                = "example-alb"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
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
      name    = "Microsoft.ServiceNetworking/trafficControllers"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_application_load_balancer_subnet_association" "example" {
  name                         = "example"
  application_load_balancer_id = azurerm_application_load_balancer.example.id
  subnet_id                    = azurerm_subnet.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Application Gateway for Containers Association. Changing this forces a new resource to be created.

* `application_load_balancer_id` - (Required) The ID of the Application Gateway for Containers. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet which the Application Gateway for Containers associated to.

~> **Note:** The subnet to be used must have a delegation for  `Microsoft.ServiceNetworking/trafficControllers` as shown in the example above.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Application Gateway for Containers Association.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Application Gateway for Containers Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Gateway for Containers Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Gateway for Containers Association.
* `update` - (Defaults to 30 minutes) Used when updating the Application Gateway for Containers Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Gateway for Containers Association.

## Import

Application Gateway for Containers Associations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_load_balancer_subnet_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/alb1/associations/association1
```
