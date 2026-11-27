change "resource-fix" {
  body = "`azurerm_log_analytics_workspace_table` - fix default handling for `retention_in_days` and `total_retention_in_days` to prevent perpetual diffs and update errors"
}
