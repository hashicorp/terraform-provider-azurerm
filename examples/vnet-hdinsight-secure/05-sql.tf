resource "random_id" "sql_dbserver_name_unique" {
  byte_length = 8
}

resource "azurerm_sql_server" "dbserver" {
  count                        = "${length(var.azuresqldb_databases) > 0 ? 1 : 0}"
  name                         = "${random_id.sql_dbserver_name_unique.hex}"
  resource_group_name          = "${azurerm_resource_group.rg.name}"
  location                     = "${azurerm_resource_group.rg.location}"
  version                      = "12.0"
  administrator_login          = "${var.sql_server_admin_user}"
  administrator_login_password = "${var.sql_server_admin_password}"
  tags                         = "${var.tags}"
}

# Enables the "Allow Access to Azure services" box as described in the API docs 
# https://docs.microsoft.com/en-us/rest/api/sql/firewallrules/createorupdate

resource "azurerm_sql_firewall_rule" "sqlfw" {
  count               = "${length(var.azuresqldb_databases) > 0 ? 1 : 0}"
  name                = "allow-azure-services"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  server_name         = "${azurerm_sql_server.dbserver.name}"
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}

resource "azurerm_sql_virtual_network_rule" "sqlvnet" {
  count               = "${length(var.azuresqldb_databases) > 0 ? 1 : 0}"
  name                = "allow-vnet"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  server_name         = "${azurerm_sql_server.dbserver.name}"
  subnet_id           = "${azurerm_subnet.subnet.id}"
}

resource "azurerm_sql_database" "sqldatabase" {
  count                            = "${length(var.azuresqldb_databases)}"
  name                             = "${var.azuresqldb_databases[count.index]}"
  resource_group_name              = "${azurerm_resource_group.rg.name}"
  location                         = "${azurerm_resource_group.rg.location}"
  edition                          = "Basic"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  create_mode                      = "Default"
  requested_service_objective_name = "Basic"
  server_name                      = "${azurerm_sql_server.dbserver.name}"
}
