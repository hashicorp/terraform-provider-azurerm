terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "rg-netapp-cool-access"
  location = "East US 2"
}

resource "azurerm_virtual_network" "example" {
  name                = "vnet-netapp-cool-access"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "snet-netapp-cool-access"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
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
  name                = "netappaccount-cool-access"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_pool" "example" {
  name                = "netapppool-cool-access"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "example" {
  name                         = "netappvolume-cool-access"
  location                     = azurerm_resource_group.example.location
  resource_group_name          = azurerm_resource_group.example.name
  account_name                 = azurerm_netapp_account.example.name
  pool_name                    = azurerm_netapp_pool.example.name
  volume_path                  = "cool-access-volume"
  service_level                = "Standard"
  subnet_id                    = azurerm_subnet.example.id
  storage_quota_in_gb          = 1000
  protocols                    = ["NFSv3"]

  cool_access_enabled          = true
  coolness_period              = 30
  cool_access_retrieval_policy = "Default"
  cool_access_tiering_policy   = "Auto"

  tags = {
    Environment = "example"
    Feature     = "cool-access"
  }
}
