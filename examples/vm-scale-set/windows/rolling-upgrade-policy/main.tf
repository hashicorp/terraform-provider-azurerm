provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "main" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_virtual_network" "main" {
  name                = "${var.prefix}-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "main" {
  name                = "${var.prefix}-publicip"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "main" {
  name                = "${var.prefix}-lb"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name

  frontend_ip_configuration {
    name                 = "internal"
    public_ip_address_id = azurerm_public_ip.main.id
  }
}

resource "azurerm_lb_backend_address_pool" "main" {
  name                = "backend-pool"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  loadbalancer_id     = azurerm_lb.main.id
}

resource "azurerm_lb_nat_pool" "main" {
  name                           = "nat-pool"
  resource_group_name            = azurerm_resource_group.main.name
  loadbalancer_id                = azurerm_lb.main.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port_start            = 80
  frontend_port_end              = 90
  backend_port                   = 8080
}

resource "azurerm_lb_probe" "main" {
  name                = "lb-probe"
  resource_group_name = azurerm_resource_group.main.name
  loadbalancer_id     = azurerm_lb.main.id
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "main" {
  name                           = "lb-rule"
  resource_group_name            = azurerm_resource_group.main.name
  loadbalancer_id                = azurerm_lb.main.id
  probe_id                       = azurerm_lb_probe.main.id
  backend_address_pool_id        = azurerm_lb_backend_address_pool.main.id
  frontend_ip_configuration_name = "internal"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
}

resource "azurerm_windows_virtual_machine_scale_set" "main" {
  name                = "${var.prefix}-vmss"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  sku                 = "Standard_F2"
  instances           = 3
  admin_username      = "adminuser"
  admin_password      = "P@ssw0rd1234!"
  health_probe_id     = azurerm_lb_probe.main.id
  upgrade_mode        = "Rolling"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name                                   = "internal"
      primary                                = true
      subnet_id                              = azurerm_subnet.internal.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.main.id]
      load_balancer_inbound_nat_rules_ids    = [azurerm_lb_nat_pool.main.id]
    }
  }

  rolling_upgrade_policy {
    max_batch_instance_percent              = 21
    max_unhealthy_instance_percent          = 22
    max_unhealthy_upgraded_instance_percent = 23
    pause_time_between_batches              = "PT30S"
  }

  depends_on = ["azurerm_lb_rule.main"]
}
