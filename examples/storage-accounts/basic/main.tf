resource "azurerm_resource_group" "example" {
  name = "${var.prefix}-resources"
  name = "${var.location}"
}

resource "azurerm_storage_account" "example" {
  # NOTE: this name needs to be globally unique
  name                     = "${var.prefix}acc"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags {
    environment = "staging"
  }
}
