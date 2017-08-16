# provider "azurerm" {
#   subscription_id = "REPLACE-WITH-YOUR-SUBSCRIPTION-ID"
#   client_id       = "REPLACE-WITH-YOUR-CLIENT-ID"
#   client_secret   = "REPLACE-WITH-YOUR-CLIENT-SECRET"
#   tenant_id       = "REPLACE-WITH-YOUR-TENANT-ID"
# }

resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group}"
  location = "${var.location}"
}

resource "azurerm_sql_database" "db" {
  name                             = "${var.db_name}"
  resource_group_name              = "${azurerm_resource_group.rg.name}"
  location                         = "${var.location}"
  edition                          = "${var.db_edition}"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  server_name                      = "${azurerm_sql_server.server.name}"
  create_mode                      = "Default"
  requested_service_objective_name = "Basic"
}

resource "azurerm_sql_server" "server" {
  name                         = "${var.resource_group}-sqlsvr"
  resource_group_name          = "${azurerm_resource_group.rg.name}"
  location                     = "${var.location}"
  version                      = "12.0"
  administrator_login          = "${var.sql_admin_username}"
  administrator_login_password = "${var.sql_password}"
}

resource "azurerm_sql_firewall_rule" "fw" {
  name                = "firewallrules"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  server_name         = "${azurerm_sql_server.server.name}"
  start_ip_address    = "${var.start_ip_address}"
  end_ip_address      = "${var.end_ip_address}"
}
