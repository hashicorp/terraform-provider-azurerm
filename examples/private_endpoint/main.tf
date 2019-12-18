resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "example" {
  name     = "${var.prefix}-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet" "example" {
  name     = "${var.prefix}-subnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.1.0/24"
  private_link_service_network_policies = "Disabled"
  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_public_ip" "example" {
  name     = "${var.prefix}-pip"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name     = "${var.prefix}-lb"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  frontend_ip_configuration {
    name                 = "${azurerm_public_ip.example.name}"
    public_ip_address_id = "${azurerm_public_ip.example.id}"
  }
}

resource "azurerm_private_link_service" "example" {
  name     = "${var.prefix}-pls"
  location = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  nat_ip_configuration {
    name = "${azurerm_public_ip.example.name}"
    subnet_id = "${azurerm_subnet.example.id}"
  }
  load_balancer_frontend_ip_configuration_ids = ["${azurerm_lb.test.frontend_ip_configuration.0.id}"]
}

resource "azurerm_private_link_endpoint" "example" {
  name     = "${var.prefix}-pe"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  subnet_id           = "${azurerm_subnet.example.id}"
  tags = {
    env = "test"
    version = "2"
  }

  private_link_service_connections {
    name = "testplsconnection"
    private_link_service_id = "${azurerm_private_link_service.example.id}"
    request_message         = "Please approve my connection"
  }
}
