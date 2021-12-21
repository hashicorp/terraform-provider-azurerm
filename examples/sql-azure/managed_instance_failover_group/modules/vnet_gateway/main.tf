resource "azurerm_subnet" "gateway_snet_1" {
  name                 = "GatewaySubnet"
  resource_group_name  = var.resource_group_name_1
  virtual_network_name = var.vnet_name_1
  address_prefixes     = [var.gateway_subnet_range_1]
}

resource "azurerm_public_ip" "pip_1" {
  name                = "${var.prefix}_pip_1"
  location            = var.location_1
  resource_group_name = var.resource_group_name_1
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "vnet_gw_1" {
  name                = "${var.prefix}_vnet_gw_1"
  location            = var.location_1
  resource_group_name = var.resource_group_name_1

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.pip_1.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.gateway_snet_1.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "gw_connection_1" {
  name                = "${var.prefix}_gw_connection_1"
  location            = var.location_1
  resource_group_name = var.resource_group_name_1

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.vnet_gw_1.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.vnet_gw_2.id

  shared_key = var.shared_key
}

resource "azurerm_subnet" "gateway_snet_2" {
  name                 = "GatewaySubnet"
  resource_group_name  = var.resource_group_name_2
  virtual_network_name = var.vnet_name_2
  address_prefixes     = [var.gateway_subnet_range_2]
}

resource "azurerm_public_ip" "pip_2" {
  name                = "${var.prefix}_pip_2"
  location            = var.location_2
  resource_group_name = var.resource_group_name_2
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "vnet_gw_2" {
  name                = "${var.prefix}_vnet_gw_2"
  location            = var.location_2
  resource_group_name = var.resource_group_name_2

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.pip_2.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.gateway_snet_2.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "gw_connection_2" {
  name                = "${var.prefix}_gw_connection_2"
  location            = var.location_2
  resource_group_name = var.resource_group_name_2

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.vnet_gw_2.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.vnet_gw_1.id

  shared_key = var.shared_key
}
