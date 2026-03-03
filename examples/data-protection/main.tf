# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-dataprotection-rg"
  location = var.location
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "${var.prefix}-dataprotection-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"
}

resource "azurerm_data_protection_backup_policy_blob_storage" "example" {
  name                                  = "${var.prefix}-blob-policy"
  vault_id                              = azurerm_data_protection_backup_vault.example.id
  operational_default_retention_duration = "P30D"
}

data "azurerm_data_protection_backup_policy_blob_storage" "example" {
  name     = azurerm_data_protection_backup_policy_blob_storage.example.name
  vault_id = azurerm_data_protection_backup_vault.example.id
}

resource "azurerm_data_protection_backup_policy_disk" "example" {
  name                            = "${var.prefix}-disk-policy"
  vault_id                        = azurerm_data_protection_backup_vault.example.id
  backup_repeating_time_intervals = ["R/2025-01-01T00:00:00+00:00/PT4H"]
  default_retention_duration      = "P7D"
}

data "azurerm_data_protection_backup_policy_disk" "example" {
  name     = azurerm_data_protection_backup_policy_disk.example.name
  vault_id = azurerm_data_protection_backup_vault.example.id
}

resource "azurerm_data_protection_backup_policy_kubernetes_cluster" "example" {
  name                = "${var.prefix}-k8s-policy"
  resource_group_name = azurerm_resource_group.example.name
  vault_name          = azurerm_data_protection_backup_vault.example.name

  backup_repeating_time_intervals = ["R/2025-01-01T06:00:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P7D"
      data_store_type = "OperationalStore"
    }
  }
}

data "azurerm_data_protection_backup_policy_kubernetes_cluster" "example" {
  name     = azurerm_data_protection_backup_policy_kubernetes_cluster.example.name
  vault_id = azurerm_data_protection_backup_vault.example.id
}

resource "azurerm_data_protection_backup_policy_mysql_flexible_server" "example" {
  name                            = "${var.prefix}-mysql-policy"
  vault_id                        = azurerm_data_protection_backup_vault.example.id
  backup_repeating_time_intervals = ["R/2025-01-01T06:00:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }
}

data "azurerm_data_protection_backup_policy_mysql_flexible_server" "example" {
  name     = azurerm_data_protection_backup_policy_mysql_flexible_server.example.name
  vault_id = azurerm_data_protection_backup_vault.example.id
}

resource "azurerm_data_protection_backup_policy_postgresql_flexible_server" "example" {
  name                            = "${var.prefix}-psql-policy"
  vault_id                        = azurerm_data_protection_backup_vault.example.id
  backup_repeating_time_intervals = ["R/2025-01-01T06:00:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }
}

data "azurerm_data_protection_backup_policy_postgresql_flexible_server" "example" {
  name     = azurerm_data_protection_backup_policy_postgresql_flexible_server.example.name
  vault_id = azurerm_data_protection_backup_vault.example.id
}
