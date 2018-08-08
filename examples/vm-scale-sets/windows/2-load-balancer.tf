resource "azurerm_lb" "main" {
  name                = "${var.prefix}-lb"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"

  frontend_ip_configuration {
    name                 = "LBFrontEnd"
    public_ip_address_id = "${azurerm_public_ip.main.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "main" {
  name                = "BackEndAddressPool"
  resource_group_name = "${azurerm_resource_group.main.name}"
  loadbalancer_id     = "${azurerm_lb.main.id}"
}

resource "azurerm_lb_nat_pool" "main" {
  resource_group_name            = "${azurerm_resource_group.main.name}"
  loadbalancer_id                = "${azurerm_lb.main.id}"
  name                           = "NATPool"
  protocol                       = "Tcp"
  frontend_port_start            = 50000
  frontend_port_end              = 50119
  backend_port                   = 22
  frontend_ip_configuration_name = "LBFrontEnd"
}
