---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_endpoint"
sidebar_current: "docs-azurerm-resource-private-endpoint"
description: |-
  Manages a Azure Private Link Endpoint instance.
---

# azurerm_private_link_endpoint

Manages a Azure `Private Link Endpoint` instance.

Azure `Private Link Endpoint` is a network interface that connects you privately and securely to a service powered by `Azure Private Link`. `Private Link Endpoint` uses a private IP address from your VNet, effectively bringing the service into your VNet. The service could be an Azure service such as Azure Storage, SQL, etc. or your own `Private Link Service`.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exampleRG"
  location = "East US"
}

resource "azurerm_virtual_network" "example" {
  name                = "examplevnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "exampleSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.1.0/24"

  disable_private_link_service_network_policies = true
  disable_private_link_endpoint_network_policies     = true
}

resource "azurerm_public_ip" "example" {
  name                = "examplePip"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "exampleLb"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.example.name
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_private_link_service" "example" {
  name                = "examplepls"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  nat_ip_configuration {
    name      = azurerm_public_ip.example.name
    subnet_id = azurerm_subnet.example.id
  }

  load_balancer_frontend_ip_configuration_ids = [azurerm_lb.test.frontend_ip_configuration.0.id]
}

resource "azurerm_private_link_endpoint" "example" {
  name                = "examplepe"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example.id

  tags = {
    env = "example"
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the Name of the `Private Link Endpoint`. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the `Private Link Endpoint` exists.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet from which the private IP addresses will be allocated.

* `tags` - (Optional) A mapping of tags assigned to the resource. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Prviate Link Endpoint.

* `network_interface_ids` - Displays an list of network interface IDs that have been created for this `Private Link Endpoint`.


## Import

The `Private Link Endpoint` can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_private_link_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/privateEndpoints/example-private-link-endpoint
```
