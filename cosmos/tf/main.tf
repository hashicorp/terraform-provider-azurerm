resource "azurerm_resource_group" "test" {
  name     = "kt-cosmos-201903-3"
  location = "westus"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "kt-cosmos-201903-3"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "BoundedStaleness"

  }

  geo_location {
    location          = "${azurerm_resource_group.test.location}"
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_database" "test" {
  name         = "letstrythisagain1133333"
  account_name = "${azurerm_cosmosdb_account.test.name}"
  account_key  = "${azurerm_cosmosdb_account.test.primary_master_key}"
}

resource "azurerm_cosmosdb_database" "testmore" {
  name         = "another_data-_base_for_us11"
  account_name = "${azurerm_cosmosdb_account.test.name}"
  account_key  = "${azurerm_cosmosdb_account.test.primary_master_key}"
}