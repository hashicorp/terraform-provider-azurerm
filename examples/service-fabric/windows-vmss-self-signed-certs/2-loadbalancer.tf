resource "azurerm_public_ip" "example" {
  name                = "${var.prefix}PIPForLB"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
  domain_name_label   = "exampleservicefabric"
}

resource "azurerm_lb" "example" {
  name                = "${var.prefix}LoadBalancer"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_lb_backend_address_pool" "example" {
  resource_group_name = azurerm_resource_group.example.name
  loadbalancer_id     = azurerm_lb.example.id
  name                = "${var.prefix}BEAPool"
}

resource "azurerm_lb_nat_pool" "example" {
  resource_group_name            = azurerm_resource_group.example.name
  loadbalancer_id                = azurerm_lb.example.id
  name                           = "${var.prefix}SFApplicationPool"
  protocol                       = "Tcp"
  frontend_port_start            = 3389
  frontend_port_end              = 4500
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.example.frontend_ip_configuration[0].name
}

resource "azurerm_lb_rule" "example_tcp" {
  resource_group_name            = azurerm_resource_group.example.name
  loadbalancer_id                = azurerm_lb.example.id
  name                           = "${var.prefix}LBRuleTcp"
  protocol                       = "Tcp"
  frontend_port                  = 19000
  backend_port                   = 19000
  frontend_ip_configuration_name = azurerm_lb.example.frontend_ip_configuration[0].name
  backend_address_pool_id        = azurerm_lb_backend_address_pool.example.id
}

resource "azurerm_lb_rule" "example_http" {
  resource_group_name            = azurerm_resource_group.example.name
  loadbalancer_id                = azurerm_lb.example.id
  name                           = "${var.prefix}LBRuleHttp"
  protocol                       = "Tcp"
  frontend_port                  = 19080
  backend_port                   = 19080
  frontend_ip_configuration_name = azurerm_lb.example.frontend_ip_configuration[0].name
  backend_address_pool_id        = azurerm_lb_backend_address_pool.example.id
}

resource "azurerm_lb_probe" "example_tcp" {
  resource_group_name = azurerm_resource_group.example.name
  loadbalancer_id     = azurerm_lb.example.id
  name                = "${var.prefix}SFTcpGatewayProbe"
  port                = 19000
}

resource "azurerm_lb_probe" "example_http" {
  resource_group_name = azurerm_resource_group.example.name
  loadbalancer_id     = azurerm_lb.example.id
  name                = "${var.prefix}SFHttpGatewayProbe"
  port                = 19080
}

