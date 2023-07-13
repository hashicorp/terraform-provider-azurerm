# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "gateway" {
  name                 = "gateway"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_subnet" "endpoint" {
  name                 = "endpoint"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  private_endpoint_network_policies_enabled = false
}

resource "azurerm_public_ip" "example" {
  name                = "${var.prefix}-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

locals {
  private_link_configuration_name        = "private_link"
  private_frontend_ip_configuration_name = "private"
}

resource "azurerm_application_gateway" "example" {
  name                = "${var.prefix}-gateway"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "Standard_v2"
    tier     = "Standard_v2"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "gateway"
    subnet_id = azurerm_subnet.gateway.id
  }

  frontend_port {
    name = "frontend"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "public"
    public_ip_address_id = azurerm_public_ip.example.id
  }

  frontend_ip_configuration {
    name                            = local.private_frontend_ip_configuration_name
    subnet_id                       = azurerm_subnet.gateway.id
    private_ip_address_allocation   = "Static"
    private_ip_address              = "10.0.1.10"
    private_link_configuration_name = local.private_link_configuration_name
  }

  private_link_configuration {
    name = local.private_link_configuration_name
    ip_configuration {
      name                          = "primary"
      subnet_id                     = azurerm_subnet.gateway.id
      private_ip_address_allocation = "Dynamic"
      primary                       = true
    }
  }

  backend_address_pool {
    name = "backend"
  }

  backend_http_settings {
    name                  = "settings"
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "listener"
    frontend_ip_configuration_name = "private"
    frontend_port_name             = "frontend"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "rule"
    rule_type                  = "Basic"
    http_listener_name         = "listener"
    backend_address_pool_name  = "backend"
    backend_http_settings_name = "settings"
  }
}

resource "azurerm_private_endpoint" "example" {
  name                = "${var.prefix}-pe"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "tfex-appgateway-connection"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_application_gateway.example.id
    subresource_names = [
      local.private_frontend_ip_configuration_name,
    ]
  }
}
