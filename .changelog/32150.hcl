change "resource-fix" {
  body = "`azurerm_kubernetes_cluster` - normalize the casing of `private_dns_zone_id` returned from the API to prevent unwanted diffs"
}
