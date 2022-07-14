provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

data "azuread_service_principal" "cosmosdb" {
  display_name = "Azure Cosmos DB"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                = "${var.prefix}kv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  purge_protection_enabled = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "List",
      "Create",
      "Delete",
      "Get",
      "Update",
    ]

  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.cosmosdb.id

    key_permissions = [
      "Get",
      "UnwrapKey",
      "WrapKey",
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "${var.prefix}key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 3072

  key_opts = [
    "decrypt",
    "encrypt",
    "wrapKey",
    "unwrapKey",
  ]
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "${var.prefix}-cosmosdb"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "MongoDB"
  key_vault_key_id    = azurerm_key_vault_key.example.id

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}
