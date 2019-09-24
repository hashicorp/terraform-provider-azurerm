resource "azurerm_resource_group" "test" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctestnetappaccount"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  active_directories {
    username            = "aduser"
    password            = "aduser"
    smb_server_name     = "SMBSERVER"
    dns                 = "1.2.3.4"
    domain              = "westcentralus.com"
    organizational_unit = "OU=FirstLevel"
  }

  tags = {
    env = "test"
  }
}