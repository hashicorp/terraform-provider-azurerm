resource "azurerm_resource_group" "network" {
  name     = "${var.resource_group_name}-${var.environment_name}"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "main" {
  name                = "${var.resource_group_name}-${var.environment_name}-net"
  address_space       = ["${var.address_space}"]
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}-${var.environment_name}"
  dns_servers         = ["${var.dns_servers}"]

  depends_on = ["azurerm_resource_group.network"]
}

resource "azurerm_subnet" "dc-subnet" {
  name                 = "${var.resource_group_name}-${var.dcsubnet_name}-${var.environment_name}"
  resource_group_name  = "${var.resource_group_name}-${var.environment_name}"
  virtual_network_name = "${azurerm_virtual_network.main.name}"
  address_prefix       = "${var.dcsubnet_prefix}"
}

resource "azurerm_subnet" "waf-subnet" {
  name                 = "${var.resource_group_name}-${var.wafsubnet_name}-${var.environment_name}"
  resource_group_name  = "${var.resource_group_name}-${var.environment_name}"
  virtual_network_name = "${azurerm_virtual_network.main.name}"
  address_prefix       = "${var.wafsubnet_prefix}"
}

resource "azurerm_subnet" "rp-subnet" {
  name                 = "${var.resource_group_name}-${var.rpsubnet_name}-${var.environment_name}"
  resource_group_name  = "${var.resource_group_name}-${var.environment_name}"
  virtual_network_name = "${azurerm_virtual_network.main.name}"
  address_prefix       = "${var.rpsubnet_prefix}"
}

resource "azurerm_subnet" "is-subnet" {
  name                 = "${var.resource_group_name}-${var.issubnet_name}-${var.environment_name}"
  resource_group_name  = "${var.resource_group_name}-${var.environment_name}"
  virtual_network_name = "${azurerm_virtual_network.main.name}"
  address_prefix       = "${var.issubnet_prefix}"
}

resource "azurerm_subnet" "db-subnet" {
  name                 = "${var.resource_group_name}-${var.dbsubnet_name}-${var.environment_name}"
  resource_group_name  = "${var.resource_group_name}-${var.environment_name}"
  virtual_network_name = "${azurerm_virtual_network.main.name}"
  address_prefix       = "${var.dbsubnet_prefix}"
}
