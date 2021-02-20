provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_mssql_server" "example" {
  name                         = "${var.prefix}-server-primary"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = var.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_server" "secondary" {
  name                         = "${var.prefix}-server-secondary"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = var.location_alt
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "example" {
  name         = "${var.prefix}-db-primary"
  server_id    = azurerm_mssql_server.example.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "BasePrice"
  sku_name     = "GP_Gen5_2"
}

resource "azurerm_mssql_database" "secondary" {
  name                        = "${var.prefix}-db-secondary"
  server_id                   = azurerm_mssql_server.secondary.id
  create_mode                 = "Secondary"
  creation_source_database_id = azurerm_mssql_database.example.id
}

resource "azurerm_sql_failover_group" "example" {
  name                = "${var.prefix}-failover-group"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mssql_server.example.name
  databases           = [azurerm_mssql_database.example.id]

  partner_servers {
    id = azurerm_mssql_server.secondary.id
  }

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }

  depends_on = [azurerm_mssql_database.secondary]
}