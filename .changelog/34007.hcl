change "resource-fix" {
  body = "`azurerm_kubernetes_cluster` - default `load_balancer_profile.idle_timeout_in_minutes` to 30 during read when omitted by the API to prevent invalid imported state"
}
