provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "xiaxintestRG-lb-cr"
  location = "eastus"
}

resource "azurerm_public_ip" "test1" {
  name                = "xiaxintest-ip1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "test1" {
  name                = "xiaxintestvn"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}
resource "azurerm_lb" "test1" {
  name                = "xiaxintest-loadbalancer-R1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "one"
    public_ip_address_id = azurerm_public_ip.test1.id
  }
}

resource "azurerm_lb_backend_address_pool" "test1" {
  loadbalancer_id = azurerm_lb.test1.id
  name            = "be"
}

resource "azurerm_lb_backend_address_pool_address" "test1"{
  name = "address1"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test1.id
  virtual_network_id      = azurerm_virtual_network.test1.id
  ip_address              = "191.168.0.1"
}

resource "azurerm_public_ip" "test2" {
  name                = "xiaxintest-ip2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test2" {
  name                = "xiaxintest-loadbalancer-R2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "one"
    public_ip_address_id = azurerm_public_ip.test2.id
  }
}

resource "azurerm_lb_backend_address_pool" "test2" {
  loadbalancer_id = azurerm_lb.test2.id
  name            = "be"
}

resource "azurerm_public_ip" "testip-cr" {
  name                = "xiaxintest-ip-cr"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "testcr" {
  name                = "xiaxin-test-loadbalancercr"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  sku_tier            = "Global"

  frontend_ip_configuration {
    name                 = "one"
    public_ip_address_id = azurerm_public_ip.testip-cr.id
  }
}

#resource "azurerm_lb_backend_address_pool" "testcrbe1" {
#  loadbalancer_id = azurerm_lb.test1.id
#  name            = "be"
#}
#
#resource "azurerm_lb_backend_address_pool" "testcrbe2" {
#  loadbalancer_id = azurerm_lb.test1.id
#  name            = "be"
#}


resource "azurerm_lb_rule" "test" {
  name                = "myHTTPRule-cr"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.testcr.id}"

  protocol      = "Tcp"
  frontend_port = 80
  backend_port  = 80

  //backend_address_pool_id = azurerm_lb_backend_address_pool.testcrbe1.id
  idle_timeout_in_minutes = 15

  frontend_ip_configuration_name = azurerm_lb.testcr.frontend_ip_configuration[0].name
  enable_tcp_reset               = true
  enable_floating_ip             = false
}