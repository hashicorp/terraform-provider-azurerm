change "resource-fix" {
  body = "`azurerm_data_protection_backup_instance_kubernetes_cluster` - wait for backup protection to be configured before completing creation"
}

change "resource-fix" {
  body = "`azurerm_data_protection_backup_instance_kubernetes_cluster` - send default datasource parameters when `backup_datasource_parameters` is omitted to satisfy Azure API requirements without making the block required"
}
