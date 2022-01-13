output "databricks_workspace_id"{
  value = azurerm_databricks_workspace.test.workspace_id
}

output "databricks_workspace_URL" {
  value = azurerm_databricks_workspace.test.workspace_url
}
