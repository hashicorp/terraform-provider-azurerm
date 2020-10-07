provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                                  = "acctestsnet"
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_network_name                  = azurerm_virtual_network.test.name
  address_prefixes                        = ["10.5.1.0/24"]
  private_link_service_network_policies = "Disabled"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.test.name
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_private_link_service" "test" {
  name                = "acctestpls"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name                         = azurerm_public_ip.test.name
    subnet_id                    = azurerm_subnet.test.id
    private_ip_address           = "10.5.1.17"
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
