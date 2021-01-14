provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

module "virtual-network" {
  source              = "./modules/virtual-network"
  resource_group_name = "${azurerm_resource_group.example.name}"
  prefix              = "${var.prefix}"
}

module "first-virtual-machine" {
  source              = "./modules/virtual-machine"
  resource_group_name = "${azurerm_resource_group.example.name}"
  prefix              = "${var.prefix}1"
  subnet_id           = "${module.virtual-network.subnet_id}"
}

module "second-virtual-machine" {
  source              = "./modules/virtual-machine"
  resource_group_name = "${azurerm_resource_group.example.name}"
  prefix              = "${var.prefix}2"
  subnet_id           = "${module.virtual-network.subnet_id}"
}

resource "azurerm_traffic_manager_profile" "example" {
  name                   = "${var.prefix}-tmprofile"
  resource_group_name    = "${azurerm_resource_group.example.name}"
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "${azurerm_resource_group.example.name}"
    ttl           = 30
  }

  monitor_config {
    protocol = "http"
    port     = 80
    path     = "/"
  }
}

resource "azurerm_traffic_manager_endpoint" "first-vm" {
  name                = "${var.prefix}-first-vm"
  resource_group_name = "${azurerm_resource_group.example.name}"
  profile_name        = "${azurerm_traffic_manager_profile.example.name}"
  target_resource_id  = "${module.first-virtual-machine.network_interface_id}"
  type                = "azureEndpoints"
  weight              = 1
}

resource "azurerm_traffic_manager_endpoint" "second-vm" {
  name                = "${var.prefix}-second-vm"
  resource_group_name = "${azurerm_resource_group.example.name}"
  profile_name        = "${azurerm_traffic_manager_profile.example.name}"
  target_resource_id  = "${module.second-virtual-machine.network_interface_id}"
  type                = "azureEndpoints"
  weight              = 1
}
