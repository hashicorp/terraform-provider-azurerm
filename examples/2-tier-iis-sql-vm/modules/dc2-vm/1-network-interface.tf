resource "azurerm_public_ip" "dc2-external" {
  name                         = "${var.prefix}-dc2-ext"
  location                     = "${var.location}"
  resource_group_name          = "${var.resource_group_name}"
  public_ip_address_allocation = "Static"
  idle_timeout_in_minutes      = 30
}

resource "azurerm_network_interface" "dc2primary" {
  name                    = "${var.prefix}-dc2-primary"
  location                = "${var.location}"
  resource_group_name     = "${var.resource_group_name}"
  internal_dns_name_label = "${local.dc2virtual_machine_name}"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = "${var.subnet_id}"
    private_ip_address_allocation = "static"
    private_ip_address            = "${var.dc2private_ip_address}"
    public_ip_address_id          = "${azurerm_public_ip.dc2-external.id}"
  }
}
