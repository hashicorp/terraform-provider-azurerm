resource "azurerm_network_security_group" "network_security_group" {
  name                = "kls-fabric-security-group"
  location            = "${var.location}"
  resource_group_name = "${module.fabric_resource_group.name}"
}

resource "azurerm_network_security_rule" "inbound_1" {
  name                        = "KPMG1"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19000"
  source_address_prefix       = "158.180.192.10/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}

resource "azurerm_network_security_rule" "inbound_2" {
  name                        = "KPMG2"
  priority                    = 110
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19080"
  source_address_prefix       = "158.180.192.10/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}

resource "azurerm_network_security_rule" "inbound_3" {
  name                        = "KPMG3"
  priority                    = 120
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19000"
  source_address_prefix       = "154.58.104.234/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}

resource "azurerm_network_security_rule" "inbound_4" {
  name                        = "KPMG4"
  priority                    = 130
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19080"
  source_address_prefix       = "154.58.104.234/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}

resource "azurerm_network_security_rule" "inbound_5" {
  name                        = "KPMG5"
  priority                    = 140
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19000"
  source_address_prefix       = "213.123.41.146/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}

resource "azurerm_network_security_rule" "inbound_6" {
  name                        = "KPMG6"
  priority                    = 150
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19080"
  source_address_prefix       = "213.123.41.146/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}

resource "azurerm_network_security_rule" "inbound_7" {
  name                        = "KPMG7"
  priority                    = 160
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19000"
  source_address_prefix       = "213.123.41.194/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}

resource "azurerm_network_security_rule" "inbound_8" {
  name                        = "KPMG8"
  priority                    = 170
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19080"
  source_address_prefix       = "213.123.41.194/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}

resource "azurerm_network_security_rule" "inbound_9" {
  name                        = "KPMG9"
  priority                    = 180
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19000"
  source_address_prefix       = "212.183.140.18/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}

resource "azurerm_network_security_rule" "inbound_10" {
  name                        = "KPMG10"
  priority                    = 190
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "19080"
  source_address_prefix       = "212.183.140.18/32"
  destination_address_prefix  = "*"
  resource_group_name         = "${module.fabric_resource_group.name}"
  network_security_group_name = "${azurerm_network_security_group.network_security_group.name}"
}
