resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.app_location}"
}

resource "azurerm_app_service_plan" "example" {
  name                = "${var.prefix}-asp"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "${var.prefix}-as"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  app_service_plan_id = "${azurerm_app_service_plan.example.id}"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "GatewaySubnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "example" {
  name                = "${var.prefix}-publicIP"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  allocation_method = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "example" {
  name                = "${var.prefix}-vnetGateway"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = "${azurerm_public_ip.example.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = "${azurerm_subnet.example.id}"
  }

  vpn_client_configuration {
    address_space = ["10.2.0.0/24"]
    vpn_client_protocols = ["SSTP"]
  }

  lifecycle {
    ignore_changes = [
        vpn_client_configuration.0.root_certificate,
    ]
  }
}

resource "azurerm_app_service_virtual_network_connection_gateway" "example" {
  app_service_name      = "${azurerm_app_service.example.name}"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  virtual_network_id    = "${azurerm_virtual_network.example.id}"
  virtual_network_gateway_id = "${azurerm_virtual_network_gateway.example.id}"
}