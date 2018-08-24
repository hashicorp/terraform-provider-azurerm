resource "azurerm_network_security_group" "allow-rdp" {
  name                = "allow-rdp"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"
}

resource "azurerm_network_security_rule" "allow-rdp" {
  name                        = "allow-rdp"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "3389"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.resource_group_name}"
  network_security_group_name = "${azurerm_network_security_group.allow-rdp.name}"
}

resource "azurerm_public_ip" "static" {
  name                         = "${var.prefix}-sql${1 + count.index}-ext"
  location                     = "${var.location}"
  resource_group_name          = "${var.resource_group_name}"
  public_ip_address_allocation = "static"
  count                        = "${var.sqlvmcount}"
  sku                          = "Standard"
}

resource "azurerm_network_interface" "primary" {
  name                    = "${var.prefix}-sql${1 + count.index}-int"
  location                = "${var.location}"
  resource_group_name     = "${var.resource_group_name}"
  internal_dns_name_label = "${var.prefix}-sql${1 + count.index}"
  network_security_group_id = "${azurerm_network_security_group.allow-rdp.id}"
  count                   = "${var.sqlvmcount}"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = "${var.subnet_id}"
    private_ip_address_allocation = "static"
    private_ip_address            = "10.100.50.${10 + count.index}"
    public_ip_address_id          = "${azurerm_public_ip.static.*.id[count.index]}"
    load_balancer_backend_address_pools_ids = ["${azurerm_lb_backend_address_pool.loadbalancer_backend.id}"]
  }
}