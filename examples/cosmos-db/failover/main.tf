provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_cosmosdb_account" "example" {
  name                      = "${var.prefix}-cosmosdb"
  location                  = "${azurerm_resource_group.example.location}"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  offer_type                = "Standard"
  kind                      = "GlobalDocumentDB"
  enable_automatic_failover = true

  //set ip_range_filter to allow azure services (0.0.0.0) and azure portal.
  ip_range_filter = "0.0.0.0,104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45,52.187.184.26"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    prefix            = "${var.prefix}-customid"
    location          = "${azurerm_resource_group.example.location}"
    failover_priority = 2
  }

  geo_location {
    location          = "${var.failover_location}"
    failover_priority = 0
  }
}
