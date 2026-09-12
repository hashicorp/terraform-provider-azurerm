# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_network_security_group" "example" {
  name                = "example-security-group"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]

  subnet {
    name             = "${var.prefix}-subnet"
    address_prefixes = ["10.0.1.0/24"]
    security_group   = azurerm_network_security_group.example.id
    delegation {
      name = "delegation"
      service_delegation {
        name    = "Microsoft.App/environments"
        actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
      }
    }
  }
}

resource "azurerm_container_app_environment" "example" {
  name                           = "${var.prefix}-environment"
  location                       = azurerm_resource_group.example.location
  resource_group_name            = azurerm_resource_group.example.name
  infrastructure_subnet_id       = azurerm_virtual_network.example.subnet.*.id[0]
  internal_load_balancer_enabled = true

  lifecycle {
    ignore_changes = [
      infrastructure_resource_group_name
    ]
  }
}

resource "azurerm_container_app" "example" {
  name                         = "${var.prefix}-app"
  resource_group_name          = azurerm_resource_group.example.name
  container_app_environment_id = azurerm_container_app_environment.example.id
  revision_mode                = "Single"

  ingress {
    external_enabled = true # allow access from outside the container app environment (e.g. from APIM)
    target_port      = 80
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }

  template {
    container {
      name   = "examplecontainerapp"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}

# Create a private DNS zone and link it to the virtual network to allow other resources in the 
# same virtual network (e.g. APIM) to resolve the container app's fqdn to its internal IP address.
resource "azurerm_private_dns_zone" "example" {
  name                = azurerm_container_app.example.ingress[0].fqdn
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "example" {
  name                  = "${var.prefix}-dnsvnetlink"
  resource_group_name   = azurerm_resource_group.example.name
  private_dns_zone_name = azurerm_private_dns_zone.example.name
  virtual_network_id    = azurerm_virtual_network.example.id
}

resource "azurerm_private_dns_a_record" "example_wildcard" {
  name                = "*"
  resource_group_name = azurerm_resource_group.example.name
  zone_name           = azurerm_private_dns_zone.example.name
  ttl                 = 3600
  records             = [azurerm_container_app_environment.example.static_ip_address]
}

resource "azurerm_private_dns_a_record" "example_naked" {
  name                = "@"
  zone_name           = azurerm_private_dns_zone.example.name
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 3600
  records             = [azurerm_container_app_environment.example.static_ip_address]
}
