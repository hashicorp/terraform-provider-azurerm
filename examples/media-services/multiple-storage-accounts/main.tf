provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}stor1"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_storage_account" "example2" {
  name                     = "${var.prefix}stor2"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_media_services" "example " {
  name                = "${var.prefix}-mediasvc"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  storage_account {
    id         = "${azurerm_storage_account.example.id}"
    is_primary = true
  }

  storage_account {
    id         = "${azurerm_storage_account.example2.id}"
    is_primary = false
  }
}

output "rendered" {
  value = "${azurerm_media_services.example.id}"
}
