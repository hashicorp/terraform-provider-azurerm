
resource "azurerm_lb" "adha_lb" {
  name                = "adha-lb"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  location            = "${azurerm_resource_group.quickstartad.location}"

  frontend_ip_configuration {
    name                          = "${var.config["ad_loadbalancer_frontend_name"]}"
    public_ip_address_id          = "${azurerm_public_ip.ad_loadbalancer_publicip.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_public_ip" "ad_loadbalancer_publicip" {
  name                         = "${var.config["network_public_ipaddress_name"]}"
  resource_group_name          = "${azurerm_resource_group.quickstartad.name}"
  location                     = "${azurerm_resource_group.quickstartad.location}"
  public_ip_address_allocation = "${var.config["network_public_ipaddress_type"]}"

}

resource "azurerm_lb_nat_rule" "pdc_rdp" {
  resource_group_name            = "${azurerm_resource_group.quickstartad.name}"
  loadbalancer_id                = "${azurerm_lb.adha_lb.id}"
  name                           = "${var.config["ad_pdc_rdp_nat_name"]}"
  protocol                       = "tcp"
  frontend_port                  = "${var.config["pdc_rdp_port"]}"
  backend_port                   = 3389
  frontend_ip_configuration_name = "${var.config["ad_loadbalancer_frontend_name"]}"
}

resource "azurerm_lb_nat_rule" "bdc_rdp" {
  resource_group_name            = "${azurerm_resource_group.quickstartad.name}"
  loadbalancer_id                = "${azurerm_lb.adha_lb.id}"
  name                           = "${var.config["ad_bdc_rdp_nat_name"]}"
  protocol                       = "tcp"
  frontend_port                  = "${var.config["bdc_rdp_port"]}"
  backend_port                   = 3389
  frontend_ip_configuration_name = "${var.config["ad_loadbalancer_frontend_name"]}"
}

resource "azurerm_lb_backend_address_pool" "adha-lb" {
  name                = "${var.config["ad_loadbalancer_backend_name"]}"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  loadbalancer_id     = "${azurerm_lb.adha_lb.id}"
}
