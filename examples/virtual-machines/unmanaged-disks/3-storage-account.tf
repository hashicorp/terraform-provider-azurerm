resource "azurerm_storage_account" "main" {
  name                     = "${var.prefix}stor"
  resource_group_name      = "${azurerm_resource_group.main.name}"
  location                 = "${azurerm_resource_group.main.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  tags                     = "${var.tags}"
}

resource "azurerm_storage_container" "disks" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.main.name}"
  storage_account_name  = "${azurerm_storage_account.main.name}"
  container_access_type = "private"
}
