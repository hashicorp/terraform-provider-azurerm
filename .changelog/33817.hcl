change "resource-fix" {
  body = "`azurerm_api_management` - set `tenant_access.enabled` to `false` on read, and no longer fail on create or update when `tenant_access.enabled` is `false`, when the direct management API is disabled for the service"
}

change "resource-fix" {
  body = "Data Source: `azurerm_api_management` - set `tenant_access.enabled` to `false` when the direct management API is disabled for the service"
}

change "resource-enhancement" {
  body = "`azurerm_api_management` - the `tenant_access` block has been deprecated and will be removed in v6.0 of the provider"
}

change "resource-enhancement" {
  body = "Data Source: `azurerm_api_management` - the `tenant_access` block has been deprecated and will be removed in v6.0 of the provider"
}
