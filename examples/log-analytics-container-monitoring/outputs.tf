// These outputs can be used to deploy the monitoring Daemonset into your k8s cluster
// https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-monitor
output "workspace_id" {
  value = "${azurerm_log_analytics_workspace.test.workspace_id}"
}

output "workspace_key" {
  value = "${azurerm_log_analytics_workspace.test.primary_shared_key}"
}
