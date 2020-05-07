provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}stor"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "example" {
  name                 = "aci-test-share"
  storage_account_name = "${azurerm_storage_account.example.name}"
  quota                = 50
}

resource "azurerm_container_group" "example" {
  name                = "${var.prefix}-continst"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  ip_address_type     = "public"
  dns_name_label      = "${var.prefix}-continst"
  os_type             = "linux"

  container {
    name     = "webserver"
    image    = "seanmckenna/aci-hellofiles"
    cpu      = "1"
    memory   = "1.5"
    port     = "80"
    protocol = "tcp"

    volume {
      name       = "logs"
      mount_path = "/aci/logs"
      read_only  = false
      share_name = "${azurerm_storage_share.example.name}"

      storage_account_name = "${azurerm_storage_account.example.name}"
      storage_account_key  = "${azurerm_storage_account.example.primary_access_key}"
    }
  }

  tags = {
    environment = "testing"
  }
}
