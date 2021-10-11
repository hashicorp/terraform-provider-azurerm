output "openshift_version" {
  value = azurerm_redhatopenshift_cluster.example.version
}

output "console_url" {
  value = azurerm_redhatopenshift_cluster.example.console_url
}
