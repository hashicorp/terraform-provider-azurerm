resource "azurerm_public_ip" "example" {
  name                = "${var.prefix}-publicip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

locals {
  frontend_ip_configuration_name = "internal"
}

resource "azurerm_lb" "example" {
  name                = "${var.prefix}-lb"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_lb_backend_address_pool" "example" {
  name                = "backend"
  loadbalancer_id     = azurerm_lb.example.id
}

resource "azurerm_lb_probe" "example" {
  name                = "probe"
  loadbalancer_id     = azurerm_lb.example.id
  protocol            = "Tcp"
  port                = 80
}

resource "azurerm_lb_rule" "example" {
  name                           = "http-lb-rule"
  loadbalancer_id                = azurerm_lb.example.id
  probe_id                       = azurerm_lb_probe.example.id
  backend_address_pool_ids = [azurerm_lb_backend_address_pool.example.id]
  frontend_ip_configuration_name = local.frontend_ip_configuration_name
  protocol                       = "Tcp"
  frontend_port                  = "80"
  backend_port                   = "80"
}
