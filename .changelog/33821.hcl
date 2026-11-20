change "resource-fix" {
  body = "`azurerm_data_protection_backup_instance_kubernetes_cluster` - wait for backup protection to be configured before completing creation"
}

change "resource-fix" {
  body = "`azurerm_data_protection_backup_instance_kubernetes_cluster` - require the `backup_datasource_parameters` block to satisfy Azure API requirements"
}
