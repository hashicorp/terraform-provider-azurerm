data "azurerm_resource_group" "example" {
  name = "${var.resource_group_name}"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-network"
  resource_group_name = "${data.azurerm_resource_group.example.name}"
  location            = "${data.azurerm_resource_group.example.location}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  resource_group_name  = "${data.azurerm_resource_group.example.name}"
  address_prefixes     = ["10.0.1.0/24"]
}
