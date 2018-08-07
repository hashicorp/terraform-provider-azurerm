resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "vnet" {
  name                = "${var.resource_group_name}vnet"
  location            = "${azurerm_resource_group.rg.location}"
  address_space       = ["10.0.0.0/16"]
  resource_group_name = "${azurerm_resource_group.rg.name}"
}

resource "azurerm_subnet" "subnet" {
  name                 = "subnet"
  address_prefix       = "10.0.0.0/24"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet.name}"
}

resource "azurerm_public_ip" "pip" {
  name                         = "${var.hostname}-pip"
  location                     = "${azurerm_resource_group.rg.location}"
  resource_group_name          = "${azurerm_resource_group.rg.name}"
  public_ip_address_allocation = "Dynamic"
  domain_name_label            = "${var.hostname}"
}

resource "azurerm_lb" "lb" {
  name                = "LoadBalancer"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  depends_on          = ["azurerm_public_ip.pip"]

  frontend_ip_configuration {
    name                 = "LBFrontEnd"
    public_ip_address_id = "${azurerm_public_ip.pip.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "backlb" {
  name                = "BackEndAddressPool"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  loadbalancer_id     = "${azurerm_lb.lb.id}"
}

resource "azurerm_lb_probe" "prob" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.lb.id}"
  name                = "ssh-running-probe"
  port                = 22
  protocol            = "Tcp"
}

resource "azurerm_lb_rule" "lbr" {
  resource_group_name            = "${azurerm_resource_group.rg.name}"
  loadbalancer_id                = "${azurerm_lb.lb.id}"
  name                           = "LBRule"
  protocol                       = "Tcp"
  frontend_port                  = 22
  backend_port                   = 22
  frontend_ip_configuration_name = "LBFrontEnd"
  probe_id                       = "${azurerm_lb_probe.prob.id}"
  backend_address_pool_id        = "${azurerm_lb_backend_address_pool.backlb.id}"
}

resource "azurerm_storage_account" "stor" {
  name                     = "${var.resource_group_name}stor"
  location                 = "${azurerm_resource_group.rg.location}"
  resource_group_name      = "${azurerm_resource_group.rg.name}"
  account_tier             = "${var.storage_account_tier}"
  account_replication_type = "${var.storage_replication_type}"
}

resource "azurerm_storage_container" "vhds" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.rg.name}"
  storage_account_name  = "${azurerm_storage_account.stor.name}"
  container_access_type = "blob"
}

resource "azurerm_virtual_machine_scale_set" "scaleset" {
  name                = "autoscalewad"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  upgrade_policy_mode = "Manual"
  overprovision       = true
  depends_on          = ["azurerm_lb.lb", "azurerm_virtual_network.vnet"]

  upgrade_policy_mode = "Rolling"

  automatic_os_upgrade = true

  rolling_upgrade_policy {
    max_batch_instance_percent              = 20
    max_unhealthy_instance_percent          = 20
    max_unhealthy_upgraded_instance_percent = 20
    pause_time_between_batches              = "PT0S"
  }

  health_probe_id = "${azurerm_lb_probe.prob.id}"
  depends_on      = ["azurerm_lb_rule.lbr"]

  sku {
    name     = "${var.vm_sku}"
    tier     = "Standard"
    capacity = "${var.instance_count}"
  }

  os_profile {
    computer_name_prefix = "${var.vmss_name_prefix}"
    admin_username       = "${var.admin_username}"
    admin_password       = "${var.admin_password}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  network_profile {
    name    = "${var.hostname}-nic"
    primary = true

    ip_configuration {
      primary                                = true
      name                                   = "${var.hostname}ipconfig"
      subnet_id                              = "${azurerm_subnet.subnet.id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.backlb.id}"]
    }
  }

  storage_profile_os_disk {
    name           = "${var.hostname}"
    caching        = "ReadWrite"
    create_option  = "FromImage"
    vhd_containers = ["${azurerm_storage_account.stor.primary_blob_endpoint}${azurerm_storage_container.vhds.name}"]
  }

  storage_profile_image_reference {
    publisher = "${var.image_publisher}"
    offer     = "${var.image_offer}"
    sku       = "${var.ubuntu_os_version}"
    version   = "latest"
  }
}
