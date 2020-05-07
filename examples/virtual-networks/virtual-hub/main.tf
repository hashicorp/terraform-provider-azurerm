provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_virtual_wan" "example" {
  name                = "${var.prefix}-virtualwan"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}

resource "azurerm_virtual_hub" "example" {
  name                = "${var.prefix}-virtualhub"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = "${azurerm_virtual_wan.example.id}"
}
