provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "endpoint" {
  name                 = "endpoint"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_private_dns_zone" "example" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_postgresql_server" "example" {
  name                = "${var.prefix}-postgresql"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "GP_Gen5_4"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_private_endpoint" "example" {
  name                 = "${var.prefix}-pe"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  subnet_id            = azurerm_subnet.endpoint.id

  private_dns_zone_group {
    name                 = "private-dns-zone-group"
    private_dns_zone_ids = [azurerm_private_dns_zone.example.id]
  }

  private_service_connection {
    name                           = "tfex-postgresql-connection"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_postgresql_server.example.id
    subresource_names              = ["postgresqlServer"]
  }
}
