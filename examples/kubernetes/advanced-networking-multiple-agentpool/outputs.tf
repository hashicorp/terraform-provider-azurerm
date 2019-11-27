output "subnet_id" {
  value = "${azurerm_kubernetes_cluster.example.agent_pool_profile.0.vnet_subnet_id}"
}

output "network_plugin" {
  value = "${azurerm_kubernetes_cluster.example.network_profile.0.network_plugin}"
}

output "service_cidr" {
  value = "${azurerm_kubernetes_cluster.example.network_profile.0.service_cidr}"
}

output "dns_service_ip" {
  value = "${azurerm_kubernetes_cluster.example.network_profile.0.dns_service_ip}"
}

output "docker_bridge_cidr" {
  value = "${azurerm_kubernetes_cluster.example.network_profile.0.docker_bridge_cidr}"
}

output "pod_cidr" {
  value = "${azurerm_kubernetes_cluster.example.network_profile.0.pod_cidr}"
}

output "kube_config_raw" {
  value = azurerm_kubernetes_cluster.example.kube_config_raw
  sensitive   = true
}

output "config" {
  value = <<CONFIGURE

Run the following commands to configure kubernetes clients:

$ terraform output kube_config_raw > ~/.kube/aksconfig
$ export KUBECONFIG=~/.kube/aksconfig

CONFIGURE

}