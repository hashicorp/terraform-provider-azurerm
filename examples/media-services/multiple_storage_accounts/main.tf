resource "azurerm_resource_group" "testrg" {
  name     = "amstestrg"
  location = "westus"
}

resource "azurerm_storage_account" "testsa" {
  name                     = "amstestsa"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

#
# Remove comments here and in the azure_media_services resource after the initial run
# to see a storage account get added to an existing AMS account.
#
/*
resource "azurerm_storage_account" "testsa2" {
  name                     = "amstestsa2"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}
*/

resource "azurerm_media_services" "ams" {
  name                = "amstest"
  location            = "${azurerm_resource_group.testrg.location}"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  storage_account {
    id         = "${azurerm_storage_account.testsa.id}"
    is_primary = true
  }

  /*
  storage_account {
    id = "${azurerm_storage_account.testsa2.id}"
    is_primary = false
  }
  */
}

output "rendered" {
  value = "${azurerm_media_services.ams.id}"
}
