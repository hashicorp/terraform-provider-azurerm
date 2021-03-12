provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  active_directory {
    username            = "pmcadmin"
    password            = var.password
    smb_server_name     = "SMBSERVER"
    dns_servers         = ["10.2.0.4"]
    domain              = "anf.local"
  }
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
  protocols           = ["CIFS","NFSv3"]
  subnet_id           = var.subnet_id
  storage_quota_in_gb = 100

  export_policy_rule {
    rule_index = 1
    allowed_clients = ["0.0.0.0/0"]
    cifs_enabled = true
    nfsv3_enabled = true
    nfsv4_enabled = false
    unix_read_write = true
  }
}
