# Azure load balancer module
resource "azurerm_resource_group" "azlb" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
  tags     = "${var.tags}"
}

resource "azurerm_public_ip" "azlb" {
  name                         = "${var.prefix}-publicIP"
  location                     = "${var.location}"
  resource_group_name          = "${azurerm_resource_group.azlb.name}"
  public_ip_address_allocation = "${var.public_ip_address_allocation}"
}

resource "azurerm_lb" "azlb" {
  name                      = "${var.prefix}-lb"  
  resource_group_name       = "${azurerm_resource_group.azlb.name}"
  location                  = "${var.location}"
  frontend_ip_configuration {
    name                 = "${var.frontend_name}"
    public_ip_address_id = "${azurerm_public_ip.azlb.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "azlb" {
  resource_group_name = "${azurerm_resource_group.azlb.name}"
  loadbalancer_id     = "${azurerm_lb.azlb.id}"
  name                = "BackEndAddressPool"
}

resource "azurerm_lb_nat_rule" "azlb" {
  count                          = "${length(var.remote_port)}"
  resource_group_name            = "${azurerm_resource_group.azlb.name}"
  loadbalancer_id                = "${azurerm_lb.azlb.id}"
  name                           = "VM-${count.index}"
  protocol                       = "tcp"
  frontend_port                  = "5000${count.index + 1}"
  backend_port                   = "${element(var.remote_port["${element(keys(var.remote_port), count.index)}"], 1)}"
  frontend_ip_configuration_name = "${var.frontend_name}"
  count                          = "${var.number_of_endpoints}"
}

resource "azurerm_lb_probe" "azlb" {
  count               = "${length(var.lb_port)}"
  resource_group_name = "${azurerm_resource_group.azlb.name}"
  loadbalancer_id     = "${azurerm_lb.azlb.id}"
  name                = "${element(keys(var.lb_port), count.index)}"
  protocol            = "${element(var.lb_port["${element(keys(var.lb_port), count.index)}"], 1)}"
  port                = "${element(var.lb_port["${element(keys(var.lb_port), count.index)}"], 2)}"
  interval_in_seconds = "${var.lb_probe_interval}"
  number_of_probes    = "${var.lb_probe_unhealthy_threshold}"
}

resource "azurerm_lb_rule" "azlb" {
  count                          = "${length(var.lb_port)}"
  resource_group_name            = "${azurerm_resource_group.azlb.name}"
  loadbalancer_id                = "${azurerm_lb.azlb.id}"
  name                           = "${element(keys(var.lb_port), count.index)}"
  protocol                       = "${element(var.lb_port["${element(keys(var.lb_port), count.index)}"], 1)}"
  frontend_port                  = "${element(var.lb_port["${element(keys(var.lb_port), count.index)}"], 0)}"
  backend_port                   = "${element(var.lb_port["${element(keys(var.lb_port), count.index)}"], 2)}"
  frontend_ip_configuration_name = "${var.frontend_name}"
  enable_floating_ip             = false
  backend_address_pool_id        = "${azurerm_lb_backend_address_pool.azlb.id}"
  idle_timeout_in_minutes        = 5
  probe_id                       = "${element(azurerm_lb_probe.azlb.*.id,count.index)}"
  depends_on                     = ["azurerm_lb_probe.azlb"]
}