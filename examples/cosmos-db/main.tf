
resource "azurerm_resource_group" "rg" {
    name     = "${var.resource_group_name}"
    location = "${var.resource_group_location}"
}

resource "random_integer" "ri" {
    min = 10000
    max = 99999
}

resource "azurerm_cosmosdb_account" "db" {
    name                = "tfex-cosmos-db-${random_integer.ri.result}"
    location            = "${azurerm_resource_group.rg.location}"
    resource_group_name = "${azurerm_resource_group.rg.name}"
    offer_type          = "Standard"
    kind                = "GlobalDocumentDB"

    enable_automatic_failover = true

    consistency_policy {
        consistency_level       = "BoundedStaleness"
        max_interval_in_seconds = 10
        max_staleness_prefix    = 200
    }

    geo_location {
      prefix              = "tfex-cosmos-db-${random_integer.ri.result}-customid"
      location          = "${azurerm_resource_group.rg.location}"
      failover_priority = 2
    }

    geo_location {
        location          = "${var.failover_location}"
        failover_priority = 0
    }
}