resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_sql_server" "example" {
  name                         = "${var.prefix}-sqlsvr"
  resource_group_name          = "${azurerm_resource_group.example.name}"
  location                     = "${azurerm_resource_group.example.location}"
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_database" "example" {
  name                             = "${var.prefix}-db"
  resource_group_name              = "${azurerm_resource_group.example.name}"
  location                         = "${azurerm_resource_group.example.location}"
  server_name                      = "${azurerm_sql_server.example.name}"
  edition                          = "Basic"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  create_mode                      = "Default"
  requested_service_objective_name = "Basic"
}

# Enables the "Allow Access to Azure services" box as described in the API docs
# https://docs.microsoft.com/en-us/rest/api/sql/firewallrules/createorupdate
resource "azurerm_sql_firewall_rule" "example" {
  name                = "allow-azure-services"
  resource_group_name = "${azurerm_resource_group.example.name}"
  server_name         = "${azurerm_sql_server.example.name}"
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}
