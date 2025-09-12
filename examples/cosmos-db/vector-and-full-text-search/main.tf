# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "${var.prefix}-cosmosdb"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  capabilities {
    name = "EnableNoSQLVectorSearch"
  }

  capabilities {
    name = "EnableNoSQLFullTextSearch"
  }

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "${var.prefix}-db"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
}

resource "azurerm_cosmosdb_sql_container" "example" {
  name                = "${var.prefix}-container"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_sql_database.example.name
  partition_key_paths = ["/id"]

  vector_embedding_policy {
    vector_embedding {
      path              = "/vector1"
      data_type         = "float32"
      distance_function = "cosine"
      dimensions        = 1536
    }
    vector_embedding {
      path              = "/vector2"
      data_type         = "uint8"
      distance_function = "euclidean"
      dimensions        = 256
    }
  }

  full_text_policy {
    default_language = "en-US"
    full_text_path {
      path = "/content"
    }
    full_text_path {
      path     = "/title"
      language = "en-US"
    }
  }

  indexing_policy {
    indexing_mode = "consistent"

    included_path {
      path = "/*"
    }

    vector_index {
      path = "/vector1"
      type = "flat"
    }

    vector_index {
      path = "/vector2"
      type = "quantizedFlat"
    }
  }
}
