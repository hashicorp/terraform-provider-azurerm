provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

module "virtual-network" {
  source              = "./modules/virtual-network"
  resource_group_name = azurerm_resource_group.example.name
  prefix              = var.prefix
}

module "first-virtual-machine" {
  source              = "./modules/virtual-machine"
  resource_group_name = azurerm_resource_group.example.name
  prefix              = "${var.prefix}1"
  subnet_id           = module.virtual-network.subnet_id
}

module "second-virtual-machine" {
  source              = "./modules/virtual-machine"
  resource_group_name = azurerm_resource_group.example.name
  prefix              = "${var.prefix}2"
  subnet_id           = module.virtual-network.subnet_id
}

resource "azurerm_traffic_manager_profile" "example" {
  name                   = "${var.prefix}-tmprofile"
  resource_group_name    = azurerm_resource_group.example.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = azurerm_resource_group.example.name
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTP"
    port     = 80
    path     = "/"
  }
}

resource "azurerm_traffic_manager_azure_endpoint" "first-vm" {
  name               = "${var.prefix}-first-vm"
  profile_id         = azurerm_traffic_manager_profile.example.id
  weight             = 1
  target_resource_id = module.first-virtual-machine.network_interface_id
}

resource "azurerm_traffic_manager_azure_endpoint" "second-vm" {
  name               = "${var.prefix}-second-vm"
  profile_id         = azurerm_traffic_manager_profile.example.id
  weight             = 1
  target_resource_id = module.second-virtual-machine.network_interface_id
}