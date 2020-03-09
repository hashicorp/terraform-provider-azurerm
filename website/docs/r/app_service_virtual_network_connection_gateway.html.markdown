---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_virtual_network_connection_gateway"
sidebar_current: "docs-azurerm-resource-app-service-virtual-network-connection-gateway"
description: |-
  Manages an App Service Virtual Network Connection Gateway.

---

# azurerm_app_service_virtual_network_connection_gateway

Manages an App Service Virtual Network Connection Gateway.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "example" {
  name                = "example"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "example-appservice"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  app_service_plan_id = "${azurerm_app_service_plan.example.id}"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-virtual-network"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "GatewaySubnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "example" {
  name                = "example"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  allocation_method = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "example" {
  name                = "example"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = "${azurerm_public_ip.example.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = "${azurerm_subnet.example.id}"
  }

  vpn_client_configuration {
    address_space = ["10.2.0.0/24"]
    vpn_client_protocols = ["SSTP"]
  }

  lifecycle {
      ignore_changes = [
        vpn_client_configuration.0.root_certificate,
      ]
  }
}

resource "azurerm_app_service_virtual_network_connection_gateway" "example" {
  app_service_name              = "${azurerm_app_service.example.name}"
  resource_group_name           = "${azurerm_resource_group.example.name}"
  virtual_network_id            = "${azurerm_virtual_network.example.id}"
  virtual_network_gateway_id 	= "${azurerm_virtual_network_gateway.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group which the app service belongs to.

* `app_service_name` - (Required) Specifies the name of the App Service.

* `virtual_network_id` - (Required) The Virtual Network's resource ID.

* `virtual_network_gateway_id` - (Required) The Virtual Network Gateway's resource ID, which should connect to the ***GatewaySubnet*** subnet of virtual network argument

---

## Attributes Reference

The following attributes are exported:

* `id` - The id of the App Service Virtual Network Connection.

* `name` - The name of the App Service Virtual Network Connection.

* `virtual_network_id` - The Virtual Network's resource ID.

* `certificate_blob` - A certificate file (.cer) blob containing the public key of the private key used to authenticate a Point-To-Site VPN connection.

* `certificate_thumbprint` - The client certificate thumbprint.

* `dns_servers` - DNS servers to be used by this Virtual Network. It is a list of IP addresses.

* `resync_required` - Is resync required?

* `routes` - (Array) One or more `route` block defined below.
  
---

A `route` block supports the following:
* `name` - Resource Name.
* `route_type` - The type of route this is: DEFAULT - By default, every app has routes to the local address ranges specified by RFC1918 INHERITED - Routes inherited from the real Virtual Network routes STATIC - Static route set on the app only. Valid values are `DEFAULT`, `INHERITED`, `STATIC`
* `start_address` The starting address for the route.
* `end_address` - The ending address for the route.
---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the App Services Virtual Network Connection Gateway.
* `update` - (Defaults to 60 minutes) Used when updating the App Services Virtual Network Connection Gateway.
* `read` - (Defaults to 10 minutes) Used when retrieving the App Services Virtual Network Connection Gateway.
* `delete` - (Defaults to 60 minutes) Used when deleting the App Services Virtual Network Connection Gateway.

## Import

App Services Virtual Network Connection Gateway can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_virtual_network_connection_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.Web/sites/example/virtualNetworkConnections/example-virtual-network
```