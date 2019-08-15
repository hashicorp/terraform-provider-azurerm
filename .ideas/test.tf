resource "azurerm_resource_group" "test" {
  name     = "example-resources-2"
  location = "West Europe"
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  address_space       = ["192.168.1.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "AzureBastionSubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "192.168.1.224/27"
}

resource "azurerm_public_ip" "test" {
  name                = "testpip"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_bastion_host" "test" {
  name                = "testbastion"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = "${azurerm_subnet.test.id}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  tags = {
    environment = "dev"
    location    = "${azurerm_resource_group.test.location}"
  }
}