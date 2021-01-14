provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-virtualnetwork"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "netapp"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netappaccount"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_netapp_pool" "example" {
  name                = "${var.prefix}-netapppool"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_netapp_account.example.name}"
  service_level       = "Premium"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "example" {
  name                = "${var.prefix}-netappvolume"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_netapp_account.example.name}"
  pool_name           = "${azurerm_netapp_pool.example.name}"
  volume_path         = "my-unique-file-path"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.example.id}"
  storage_quota_in_gb = 100
}

resource "azurerm_netapp_snapshot" "example" {
  name                = "${var.prefix}-netappsnapshot"
  account_name        = "${azurerm_netapp_account.example.name}"
  pool_name           = "${azurerm_netapp_pool.example.name}"
  volume_name         = "${azurerm_netapp_volume.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}
