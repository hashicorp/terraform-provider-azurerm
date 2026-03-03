# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

output "blob_policy_id" {
  value = data.azurerm_data_protection_backup_policy_blob_storage.example.id
}

output "disk_policy_id" {
  value = data.azurerm_data_protection_backup_policy_disk.example.id
}

output "kubernetes_policy_id" {
  value = data.azurerm_data_protection_backup_policy_kubernetes_cluster.example.id
}

output "mysql_flex_policy_id" {
  value = data.azurerm_data_protection_backup_policy_mysql_flexible_server.example.id
}

output "psql_flex_policy_id" {
  value = data.azurerm_data_protection_backup_policy_postgresql_flexible_server.example.id
}
