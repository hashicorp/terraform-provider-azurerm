# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-cdn-frontdoor-loadBalancer-privateLinkOrigin"
  location = "westeurope"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.5.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_public_ip" "example" {
  name                = "${var.prefix}-api"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "${var.prefix}-lb"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.example.name
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_private_link_service" "example" {
  name                = "${var.prefix}-privateLinkService"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  visibility_subscription_ids                 = [data.azurerm_client_config.current.subscription_id]
  load_balancer_frontend_ip_configuration_ids = [azurerm_lb.example.frontend_ip_configuration.0.id]

  nat_ip_configuration {
    name                       = "primary"
    private_ip_address         = "10.5.1.17"
    private_ip_address_version = "IPv4"
    subnet_id                  = azurerm_subnet.example.id
    primary                    = true
  }
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  depends_on = [azurerm_private_link_service.example]

  name                = "${var.prefix}-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Premium_AzureFrontDoor"

  response_timeout_seconds = 120

  tags = {
    environment = "example"
  }
}

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "${var.prefix}-origin-group"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  session_affinity_enabled = false

  restore_traffic_time_to_healed_or_new_endpoint_in_minutes = 10

  health_probe {
    interval_in_seconds = 100
    path                = "/"
    protocol            = "Http"
    request_type        = "HEAD"
  }

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }
}

resource "azurerm_cdn_frontdoor_endpoint" "example" {
  name                     = "${var.prefix}-endpoint"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  enabled                  = true
}

resource "azurerm_cdn_frontdoor_origin" "example" {
  name                          = "${var.prefix}-origin"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id
  enabled                       = true

  certificate_name_check_enabled = true # Required for Private Link
  host_name                      = azurerm_private_link_service.example.nat_ip_configuration.0.private_ip_address
  origin_host_header             = azurerm_private_link_service.example.nat_ip_configuration.0.private_ip_address
  priority                       = 1
  weight                         = 500

  private_link {
    request_message        = "Request access for CDN Frontdoor Private Link Origin Load Balancer Example"
    location               = azurerm_resource_group.example.location
    private_link_target_id = azurerm_private_link_service.example.id
  }
}
