provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netappaccount"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  active_directory {
      username            = "aduser"
      password            = "aduserpwd"
      smb_server_name     = "SMBSERVER"
      dns_servers         = ["1.2.3.4"]
      domain              = "westcentralus.com"
      organizational_unit = "OU=FirstLevel"
   }
}
