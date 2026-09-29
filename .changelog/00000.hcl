change "resource-fix" {
  body = "`azurerm_network_watcher_flow_log` - remove client-side validation restricting `version` to `1` or `2`, allowing the service to validate supported values"
}
