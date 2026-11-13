change "breaking" {
  body = "`azurerm_mssql_managed_instance_start_stop_schedule` - the `schedule` block will change from a list to a set in version 6.0 of the AzureRM Provider"
}

change "resource-fix" {
  body = "`azurerm_mssql_managed_instance_start_stop_schedule` - prevent persistent differences when the API reorders otherwise unchanged `schedule` blocks"
}
