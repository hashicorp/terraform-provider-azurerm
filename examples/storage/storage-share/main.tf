provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_storage_account" "example" {
    name                     = "${var.prefix}storageacct"
    resource_group_name      = azurerm_resource_group.example.name
    location                 = azurerm_resource_group.example.location
    account_tier             = "Premium"
    account_replication_type = "LRS"
	account_kind             = "FileStorage"
}

resource "azurerm_storage_share" "example" {
  name                 = "${var.prefix}storageshare"
  storage_account_name = azurerm_storage_account.example.name
}
