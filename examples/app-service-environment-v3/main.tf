terraform {
  required_providers {
    azurerm = {
      version = ">=2.76.0"
    }
    random = {
      version = "3.1.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = var.ase_resource_group_name
  location = var.location
}

data "azurerm_subnet" "snet-exist" {
  count                = var.use_existing_vnet_and_subnet ? 1 : 0
  name                 = var.subnet_name
  virtual_network_name = var.virtual_network_name
  resource_group_name  = var.vnet_resource_group_name
}

resource "azurerm_network_security_group" "nsg" {
  count               = var.use_existing_vnet_and_subnet ? 0 : 1
  name                = var.network_security_group_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_virtual_network" "vnet" {
  count               = var.use_existing_vnet_and_subnet ? 0 : 1
  name                = var.virtual_network_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  address_space       = var.vnet_address_prefixes
}

resource "azurerm_subnet" "snet" {
  count                = var.use_existing_vnet_and_subnet ? 0 : 1
  name                 = var.subnet_name
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet[0].name
  address_prefixes     = var.subnet_address_prefixes

  delegation {
    name = "Microsoft.Web.hostingEnvironments"
    service_delegation {
      name    = "Microsoft.Web/hostingEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet" "snet-delegation" {
  count                = var.use_existing_vnet_and_subnet ? 1 : 0
  name                 = var.subnet_name
  virtual_network_name = var.virtual_network_name
  resource_group_name  = var.vnet_resource_group_name
  address_prefixes     = data.azurerm_subnet.snet-exist[0].address_prefixes

  delegation {
    name = "Microsoft.Web.hostingEnvironments"
    service_delegation {
      name    = "Microsoft.Web/hostingEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "nsg-association" {
  count                     = var.use_existing_vnet_and_subnet ? 0 : 1
  subnet_id                 = azurerm_subnet.snet[0].id
  network_security_group_id = azurerm_network_security_group.nsg[0].id
}

resource "azurerm_app_service_environment_v3" "asev3" {
  name                = var.ase_name
  resource_group_name = azurerm_resource_group.rg.name
  subnet_id           = var.use_existing_vnet_and_subnet ? azurerm_subnet.snet-delegation[0].id : azurerm_subnet.snet[0].id

  dedicated_host_count         = var.dedicated_host_count >= 2 ? var.dedicated_host_count : null
  zone_redundant               = var.zone_redundant ? var.zone_redundant : null
  internal_load_balancing_mode = var.internal_load_balancing_mode
}

resource "azurerm_private_dns_zone" "zone" {
  count               = (var.create_private_dns && var.internal_load_balancing_mode == "Web, Publishing") ? 1 : 0
  name                = azurerm_app_service_environment_v3.asev3.dns_suffix
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_private_dns_a_record" "a-wildcard" {
  count               = (var.create_private_dns && var.internal_load_balancing_mode == "Web, Publishing") ? 1 : 0
  name                = "*"
  zone_name           = azurerm_private_dns_zone.zone[0].name
  resource_group_name = azurerm_resource_group.rg.name
  ttl                 = 3600
  records             = azurerm_app_service_environment_v3.asev3.internal_inbound_ip_addresses
}

resource "azurerm_private_dns_a_record" "a-scm" {
  count               = (var.create_private_dns && var.internal_load_balancing_mode == "Web, Publishing") ? 1 : 0
  name                = "*.scm"
  zone_name           = azurerm_private_dns_zone.zone[0].name
  resource_group_name = azurerm_resource_group.rg.name
  ttl                 = 3600
  records             = azurerm_app_service_environment_v3.asev3.internal_inbound_ip_addresses
}

resource "azurerm_private_dns_a_record" "a-at" {
  count               = (var.create_private_dns && var.internal_load_balancing_mode == "Web, Publishing") ? 1 : 0
  name                = "@"
  zone_name           = azurerm_private_dns_zone.zone[0].name
  resource_group_name = azurerm_resource_group.rg.name
  ttl                 = 3600
  records             = azurerm_app_service_environment_v3.asev3.internal_inbound_ip_addresses
}