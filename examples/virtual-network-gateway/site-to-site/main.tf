resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-network"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "GatewaySubnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_local_network_gateway" "example" {
  name                = "onpremise"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  gateway_address     = "168.62.225.23"
  address_space       = ["10.1.1.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                         = "${var.prefix}-pip"
  location                     = "${azurerm_resource_group.example.location}"
  resource_group_name          = "${azurerm_resource_group.example.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "example" {
  name                = "${var.prefix}-vng"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  type                = "Vpn"
  vpn_type            = "RouteBased"
  active_active       = false
  enable_bgp          = false
  sku                 = "Basic"

  ip_configuration {
    public_ip_address_id          = "${azurerm_public_ip.example.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = "${azurerm_subnet.example.id}"
  }
}

resource "azurerm_virtual_network_gateway_connection" "example" {
  name                       = "onpremise"
  location                   = "${azurerm_resource_group.example.location}"
  resource_group_name        = "${azurerm_resource_group.example.name}"
  type                       = "IPsec"
  virtual_network_gateway_id = "${azurerm_virtual_network_gateway.example.id}"
  local_network_gateway_id   = "${azurerm_local_network_gateway.example.id}"
  shared_key                 = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
