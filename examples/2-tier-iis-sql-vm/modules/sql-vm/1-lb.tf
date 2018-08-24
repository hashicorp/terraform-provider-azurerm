resource "azurerm_lb" "sql-loadbalancer" {
  name                = "${var.prefix}-sql-loadbalancer"
  resource_group_name = "${var.resource_group_name}"
  location            = "${var.location}"
  sku                 = "Standard"
  frontend_ip_configuration {
    name                 = "LoadBalancerFrontEnd"
    subnet_id                     = "${var.subnet_id}"
    private_ip_address_allocation = "static"
    private_ip_address            = "${var.lbprivate_ip_address}"
  }
}
resource "azurerm_lb_backend_address_pool" "loadbalancer_backend" {
  name                = "loadbalancer_backend"
  resource_group_name = "${var.resource_group_name}"
  loadbalancer_id     = "${azurerm_lb.sql-loadbalancer.id}"
}
resource "azurerm_lb_probe" "loadbalancer_probe" {
  resource_group_name = "${var.resource_group_name}"
  loadbalancer_id     = "${azurerm_lb.sql-loadbalancer.id}"
  name                = "SQLAlwaysOnEndPointProbe"
  protocol            = "tcp"
  port                = 59999
  interval_in_seconds = 5
  number_of_probes    = 2
}

resource "azurerm_lb_rule" "SQLAlwaysOnEndPointListener" {
  resource_group_name            = "${var.resource_group_name}"
  loadbalancer_id                = "${azurerm_lb.sql-loadbalancer.id}"
  name                           = "SQLAlwaysOnEndPointListener"
  protocol                       = "Tcp"
  frontend_port                  = 1433
  backend_port                   = 1433
  frontend_ip_configuration_name = "LoadBalancerFrontEnd"
  backend_address_pool_id        = "${azurerm_lb_backend_address_pool.loadbalancer_backend.id}"
  probe_id                       = "${azurerm_lb_probe.loadbalancer_probe.id}"
}
