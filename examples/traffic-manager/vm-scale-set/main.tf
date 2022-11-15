provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-tmresources"
  location = var.location
}

resource "azurerm_traffic_manager_profile" "example" {
  name                   = "${var.prefix}-trafficmgr"
  resource_group_name    = azurerm_resource_group.example.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "${var.prefix}-trafficmgr"
    ttl           = 100
  }

  monitor_config {
    protocol = "HTTP"
    port     = 80
    path     = "/"
  }
}

module "region1" {
  source   = "./modules/region"
  prefix   = "${var.prefix}-region1"
  location = var.location
}

resource "azurerm_traffic_manager_azure_endpoint" "region1" {
  name               = "${var.prefix}-region1"
  profile_id         = azurerm_traffic_manager_profile.example.id
  weight             = 100
  target_resource_id = module.region1.public_ip_address_id
}


module "region2" {
  source   = "./modules/region"
  prefix   = "${var.prefix}-region1"
  location = var.alt_location
}

resource "azurerm_traffic_manager_azure_endpoint" "region2" {
  name               = "${var.prefix}-region2"
  profile_id         = azurerm_traffic_manager_profile.example.id
  weight             = 100
  target_resource_id = module.region2.public_ip_address_id
}
