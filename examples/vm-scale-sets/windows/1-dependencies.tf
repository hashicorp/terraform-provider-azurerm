locals {
  virtual_machine_name = "${var.prefix}-vmss"
}

resource "azurerm_resource_group" "main" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "main" {
  name                = "${var.prefix}-network"
  location            = "${azurerm_resource_group.main.location}"
  address_space       = ["10.0.0.0/16"]
  resource_group_name = "${azurerm_resource_group.main.name}"
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  address_prefix       = "10.0.0.0/24"
  resource_group_name  = "${azurerm_resource_group.main.name}"
  virtual_network_name = "${azurerm_virtual_network.main.name}"
}

resource "azurerm_public_ip" "main" {
  name                         = "${var.prefix}-pip"
  location                     = "${azurerm_resource_group.main.location}"
  resource_group_name          = "${azurerm_resource_group.main.name}"
  public_ip_address_allocation = "Dynamic"
  domain_name_label            = "${local.virtual_machine_name}"
}
