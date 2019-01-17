provider "azurerm" {
		version = "~> 999"

        subscription_id="7060bca0-7a3c-44bd-b54c-4bb1e9facfac"
        client_id="32c7a76b-6d19-4395-9ed2-ac1ee3e1aec0"
        client_secret="4ba02547-5c24-46b1-a805-f243040b56cc"
        tenant_id="72f988bf-86f1-41af-91ab-2d7cd011db47"
}

# Can you also get the id by adding a azurerm_storage_account that references the
# same properties as the existing storage account.
variable "existing_storage_account_id" {
  default = "/subscriptions/7060bca0-7a3c-44bd-b54c-4bb1e9facfac/resourcegroups/seusher_dev/providers/Microsoft.Storage/storageAccounts/seusherams3"
}

resource "azurerm_resource_group" "testrg" {
  name     = "seusheramsrg"
  location = "westus"
}

resource "azurerm_media_services" "ams" {

        name = "seushertestamstf"
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