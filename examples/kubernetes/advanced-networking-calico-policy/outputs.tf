output "subnet_id" {
  value = "${azurerm_kubernetes_cluster.test.agent_pool_profile.0.vnet_subnet_id}"
}

output "network_plugin" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.network_plugin}"
}

output "network_policy" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.network_policy}"
}

output "service_cidr" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.service_cidr}"
}

output "dns_service_ip" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.dns_service_ip}"
}

output "docker_bridge_cidr" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.docker_bridge_cidr}"
}

output "pod_cidr" {
  value = "${azurerm_kubernetes_cluster.test.network_profile.0.pod_cidr}"
}
