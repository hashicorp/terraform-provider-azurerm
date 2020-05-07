resource "azurerm_public_ip" "static" {
  name                         = "${var.prefix}-client-ppip"
  location                     = "${var.location}"
  resource_group_name          = "${var.resource_group_name}"
  allocation_method = "Static"
}

resource "azurerm_network_interface" "primary" {
  name                    = "${var.prefix}-client-nic"
  location                = "${var.location}"
  resource_group_name     = "${var.resource_group_name}"
  internal_dns_name_label = "${local.virtual_machine_name}"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = "${var.subnet_id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.static.id}"
  }
}
