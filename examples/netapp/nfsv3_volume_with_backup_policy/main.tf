provider "azurerm" {
  features {
    netapp {
      prevent_volume_destruction = false
      delete_backups_on_backup_vault_destroy = true
	  }
  }
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-virtualnetwork"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_backup_vault" "example" {
  name                = "${var.prefix}-netappbackupvault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  account_name        = azurerm_netapp_account.example.name
}

resource "azurerm_netapp_backup_policy" "example" {
  name                    = "${var.prefix}-netappbackuppolicy"
  resource_group_name     = azurerm_resource_group.example.name
  location                = azurerm_resource_group.example.location
  account_name            = azurerm_netapp_account.example.name
  daily_backups_to_keep   = 2
  weekly_backups_to_keep  = 2
  monthly_backups_to_keep = 2
  enabled = true
}

resource "azurerm_netapp_pool" "example" {
  name                = "${var.prefix}-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "example" {
  lifecycle {
    prevent_destroy = true
  }

  name                = "${var.prefix}-netappvolume"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "${var.prefix}-netappvolume"
  service_level       = "Standard"
  protocols           = ["NFSv3"]
  network_features    = "Basic"
  subnet_id           = azurerm_subnet.example.id
  storage_quota_in_gb = 100

  data_protection_backup_policy {
    backup_vault_id  = azurerm_netapp_backup_vault.example.id
    backup_policy_id = azurerm_netapp_backup_policy.example.id
    policy_enforced  = true
  }

  export_policy_rule {
    rule_index = 1
    allowed_clients = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_write = true
  }
}
