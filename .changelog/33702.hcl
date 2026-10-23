change "resource-fix" {
  body = "`azurerm_kubernetes_cluster` - relax default node pool name validation to allow up to 39 characters, subject to AKS service limits"
}
change "resource-fix" {
  body = "`azurerm_kubernetes_cluster_node_pool` - relax node pool name validation to allow up to 39 characters, subject to AKS service limits"
}
change "resource-fix" {
  body = "Data Source: `azurerm_kubernetes_cluster_node_pool` - allow reading node pools with names up to 39 characters"
}
