change "resource-fix" {
  body = "`azurerm_kubernetes_cluster_node_pool` - add additional locking on the pod and node virtual network IDs to prevent conflicts when creating multiple node pools on the same virtual network"
}
