provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-tmresources"
  location = "${var.location}"
}

resource "azurerm_traffic_manager_profile" "example" {
  name                   = "${var.prefix}-trafficmgr"
  resource_group_name    = "${azurerm_resource_group.example.name}"
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "${var.prefix}-trafficmgr"
    ttl           = 100
  }

  monitor_config {
    protocol = "http"
    port     = 80
    path     = "/"
  }
}

module "region1" {
  source   = "./modules/region"
  prefix   = "${var.prefix}-region1"
  location = "${var.location}"
}

resource "azurerm_traffic_manager_endpoint" "region1" {
  name                = "${var.prefix}-region1"
  resource_group_name = "${azurerm_resource_group.example.name}"
  profile_name        = "${azurerm_traffic_manager_profile.example.name}"
  target_resource_id  = "${module.region1.public_ip_address_id}"
  type                = "azureEndpoints"
  weight              = 100
}

module "region2" {
  source   = "./modules/region"
  prefix   = "${var.prefix}-region1"
  location = "${var.alt_location}"
}

resource "azurerm_traffic_manager_endpoint" "region2" {
  name                = "${var.prefix}-region2"
  resource_group_name = "${azurerm_resource_group.example.name}"
  profile_name        = "${azurerm_traffic_manager_profile.example.name}"
  target_resource_id  = "${module.region2.public_ip_address_id}"
  type                = "azureEndpoints"
  weight              = 100
}
