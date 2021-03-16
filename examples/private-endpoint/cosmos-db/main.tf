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

resource "azurerm_cosmosdb_account" "example" {
  name                = "${var.prefix}-cosmosdb-example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  enable_automatic_failover         = false
  is_virtual_network_filter_enabled = true
  ip_range_filter                   = "104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45,52.187.184.26"

  capabilities {
    name = "EnableMongo"
  }

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 310
    max_staleness_prefix    = 101000
  }

  geo_location {
    prefix            = "${var.prefix}-cosmos-db-customid"
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_private_endpoint" "example" {
  name                 = "${var.prefix}-pe"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  subnet_id            = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "tfex-cosmosdb-connection"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_cosmosdb_account.example.id
    subresource_names              = ["MongoDB"]
  }
}
