resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "eastus"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststorageaccount"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

}

resource "azurerm_sql_server_blob_auditing_policies" "test"{
  resource_group_name               = "${azurerm_resource_group.test.name}"
  server_name                       = "${azurerm_sql_server.test.name}"
  state                             = "Enabled"
  storage_endpoint                  = "${azurerm_storage_account.test.primary_blob_endpoint}"
  storage_account_access_key        = "${azurerm_storage_account.test.primary_access_key}"
  retention_days                    = 6
  is_storage_secondary_key_in_use   = true
  audit_actions_and_groups          = "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP,FAILED_DATABASE_AUTHENTICATION_GROUP"
  is_azure_monitor_target_enabled   = true

}