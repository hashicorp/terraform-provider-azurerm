provider "azurerm" {
		version = "~> 999"

        subscription_id="7060bca0-7a3c-44bd-b54c-4bb1e9facfac"
        client_id="32c7a76b-6d19-4395-9ed2-ac1ee3e1aec0"
        client_secret="4ba02547-5c24-46b1-a805-f243040b56cc"
        tenant_id="72f988bf-86f1-41af-91ab-2d7cd011db47"
}

resource "azurerm_resource_group" "testrg" {
  name     = "seusheramsrg"
  location = "westus"
}

resource "azurerm_storage_account" "testsa1" {
  name                     = "seusheramsstore1"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags {
    environment = "staging"
  }
}

#
# Remove comments here and in the azure_media_services resource after the initial run
# to see a storage account get added to an existing AMS account.
#
/*
resource "azurerm_storage_account" "testsa2" {
  name                     = "seusheramsstore2"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags {
    environment = "staging"
  }
}
*/

resource "azurerm_media_services" "ams" {

        name = "seushertestamstf"
        location = "${azurerm_resource_group.testrg.location}"
        resource_group_name = "${azurerm_resource_group.testrg.name}"
		
        storage_account {
				id = "${azurerm_storage_account.testsa1.id}"
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