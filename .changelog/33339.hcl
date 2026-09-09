
change "resource-fix" {
  body = "`azurerm_resource_group` - the `managed_by` property now forces recreation when changed as the API does not support changing this value"
}
