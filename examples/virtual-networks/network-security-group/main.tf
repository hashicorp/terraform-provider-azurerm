provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-network"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "example" {
  name                = "${var.prefix}-nsg"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}

# NOTE: this allows SSH from any network
resource "azurerm_network_security_rule" "ssh" {
  name                        = "ssh"
  resource_group_name         = "${azurerm_resource_group.example.name}"
  network_security_group_name = "${azurerm_network_security_group.example.name}"
  priority                    = 102
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
}
