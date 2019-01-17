# Can you also get the id by adding a azurerm_storage_account that references the
# same properties as the existing storage account.
variable "existing_storage_account_id" {
  default = "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/testrg/providers/Microsoft.Storage/storageAccounts/teststorage"
}

resource "azurerm_resource_group" "testrg" {
  name     = "amstestrg"
  location = "westus"
}

resource "azurerm_media_services" "ams" {

        name = "amstest"
        location = "${azurerm_resource_group.testrg.location}"
        resource_group_name = "${azurerm_resource_group.testrg.name}"
		
        storage_account {
				  id = "${var.existing_storage_account_id}"
				  is_primary = true
		    }
}

output "rendered" {
  value = "${azurerm_media_services.ams.id}"
}