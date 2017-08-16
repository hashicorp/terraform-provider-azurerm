resource "azurerm_virtual_network" "adha_vnet" {
  name                = "${var.config["vnet_name"]}"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  location            = "${azurerm_resource_group.quickstartad.location}"
  address_space       = ["${var.config["vnet_address_range"]}"]
}

resource "azurerm_subnet" "ad_subnet" {
  name                = "${var.config["subnet_name"]}"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  address_prefix       = "${var.config["subnet_address_range"]}"
  virtual_network_name = "${azurerm_virtual_network.adha_vnet.name}"
}