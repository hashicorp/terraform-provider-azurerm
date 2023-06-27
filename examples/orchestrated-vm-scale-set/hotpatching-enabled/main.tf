# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

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

resource "azurerm_network_interface" "main" {
  name                = "${var.prefix}-nic"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.internal.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "main" {
  name                = "${var.prefix}-OVMSS"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name

  sku_name  = "Standard_F2s_v2"
  instances = 1

  platform_fault_domain_count = 2

  os_profile {
    windows_configuration {
      computer_name_prefix = var.prefix
      admin_username       = "adminuser"
      admin_password       = "P@$$w0rd1234!"

      winrm_listener {
        protocol = "Http"
      }
    }
  }

  network_interface {
    name    = "${var.prefix}-NetworkProfile"
    primary = true

    ip_configuration {
      name      = "PrimaryIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.internal.id

      public_ip_address {
        name                    = "${var.prefix}-PublicIpConfiguration"
        domain_name_label       = "${var.prefix}-domain-label"
        idle_timeout_in_minutes = 4
      }
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition-core"
    version   = "latest"
  }

  extension {
    name                               = "${var.prefix}-HealthExtension"
    publisher                          = "Microsoft.ManagedServices"
    type                               = "ApplicationHealthWindows"
    type_handler_version               = "1.0"
    auto_upgrade_minor_version_enabled = true

    settings = jsonencode({
      "protocol"    = "http"
      "port"        = "80"
      "requestPath" = "/healthEndpoint"
    })
  }
}
