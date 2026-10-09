change "resource-fix" {
  body = "`azurerm_policy_definition` - the `parameters` property now forces recreation when parameter names change as the API does not support renaming parameters"
}
