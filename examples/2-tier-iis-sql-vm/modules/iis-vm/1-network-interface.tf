resource "azurerm_public_ip" "static" {
  name                         = "${var.prefix}-iis${1 + count.index}-ext"
  location                     = "${var.location}"
  resource_group_name          = "${var.resource_group_name}"
  public_ip_address_allocation = "static"
  count                        = "${var.vmcount}"
}

resource "azurerm_network_interface" "primary" {
  name                    = "${var.prefix}-iis${1 + count.index}-int"
  location                = "${var.location}"
  resource_group_name     = "${var.resource_group_name}"
  internal_dns_name_label = "${var.prefix}-iis${1 + count.index}"
  count                   = "${var.vmcount}"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = "${var.subnet_id}"
    private_ip_address_allocation = "static"
    private_ip_address            = "10.100.30.${10 + count.index}"
    public_ip_address_id          = "${azurerm_public_ip.static.*.id[count.index]}"
  }
}
