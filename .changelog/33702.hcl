change "resource-enhancement" {
  body = "`azurerm_kubernetes_cluster` - support configuring the maximum default node pool name length up to 39 characters with the `kubernetes_cluster.node_pool_name_max_length` provider feature"
}
change "resource-enhancement" {
  body = "`azurerm_kubernetes_cluster_node_pool` - support configuring the maximum node pool name length up to 39 characters with the `kubernetes_cluster.node_pool_name_max_length` provider feature"
}
change "resource-fix" {
  body = "Data Source: `azurerm_kubernetes_cluster_node_pool` - allow reading node pools with names up to 39 characters"
}
