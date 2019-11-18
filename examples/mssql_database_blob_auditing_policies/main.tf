resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_sql_server" "example" {
  name                         = "${var.prefix}-sqlsvr"
  resource_group_name          = "${azurerm_resource_group.example.name}"
  location                     = "${azurerm_resource_group.example.location}"
  version                      = "12.0"
  administrator_login          = "${var.login}"
  administrator_login_password = "${var.login_pwd}"
}

resource "azurerm_sql_database" "example" {
  name                             = "${var.prefix}-sqldb"
  resource_group_name              = "${azurerm_resource_group.example.name}"
  server_name                      = "${azurerm_sql_server.example.name}"
  location                         = "${azurerm_resource_group.example.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}-storageAccount"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

}

resource "azurerm_mssql_database_blob_auditing_policies" "example"{
  resource_group_name               = "${azurerm_resource_group.example.name}"
  server_name                       = "${azurerm_sql_server.example.name}"
  database_name                     = "${azurerm_sql_database.example.name}"
  state                             = "Enabled"
  storage_endpoint                  = "${azurerm_storage_account.example.primary_blob_endpoint}"
  storage_account_access_key        = "${azurerm_storage_account.example.primary_access_key}"
  retention_days                    = 6
  is_storage_secondary_key_in_use   = true
  audit_actions_and_groups          = "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP,FAILED_DATABASE_AUTHENTICATION_GROUP"
  is_azure_monitor_target_enabled   = true

}