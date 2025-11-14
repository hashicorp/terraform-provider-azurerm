output "resource_group_id" {
  description = "The ID of the resource group"
  value       = azurerm_resource_group.example.id
}

output "instance_id" {
  description = "The ID of the IoT Operations instance"
  value       = azurerm_iotoperations_instance.example.id
}

# High Performance Profile Outputs
output "high_performance_profile_id" {
  description = "The ID of the high-performance dataflow profile"
  value       = azurerm_iotoperations_dataflow_profile.high_performance.id
}

output "high_performance_profile_name" {
  description = "The name of the high-performance dataflow profile"
  value       = azurerm_iotoperations_dataflow_profile.high_performance.name
}

# Standard Profile Outputs
output "standard_profile_id" {
  description = "The ID of the standard dataflow profile"
  value       = azurerm_iotoperations_dataflow_profile.standard.id
}

output "standard_profile_name" {
  description = "The name of the standard dataflow profile"
  value       = azurerm_iotoperations_dataflow_profile.standard.name
}

# Edge Profile Outputs
output "edge_profile_id" {
  description = "The ID of the edge dataflow profile"
  value       = azurerm_iotoperations_dataflow_profile.edge.id
}

output "edge_profile_name" {
  description = "The name of the edge dataflow profile"
  value       = azurerm_iotoperations_dataflow_profile.edge.name
}

# Development Profile Outputs
output "development_profile_id" {
  description = "The ID of the development dataflow profile"
  value       = var.create_development_profile ? azurerm_iotoperations_dataflow_profile.development[0].id : null
}

output "development_profile_name" {
  description = "The name of the development dataflow profile"
  value       = var.create_development_profile ? azurerm_iotoperations_dataflow_profile.development[0].name : null
}

# Profile Configuration Summary
output "profiles_summary" {
  description = "Summary of all configured dataflow profiles"
  value = {
    high_performance = {
      id             = azurerm_iotoperations_dataflow_profile.high_performance.id
      name           = azurerm_iotoperations_dataflow_profile.high_performance.name
      instance_count = azurerm_iotoperations_dataflow_profile.high_performance.instance_count
      log_level      = var.high_performance_log_level
      prometheus_port = var.high_performance_prometheus_port
    }
    standard = {
      id             = azurerm_iotoperations_dataflow_profile.standard.id
      name           = azurerm_iotoperations_dataflow_profile.standard.name
      instance_count = azurerm_iotoperations_dataflow_profile.standard.instance_count
      log_level      = var.standard_log_level
      prometheus_port = var.standard_prometheus_port
    }
    edge = {
      id             = azurerm_iotoperations_dataflow_profile.edge.id
      name           = azurerm_iotoperations_dataflow_profile.edge.name
      instance_count = azurerm_iotoperations_dataflow_profile.edge.instance_count
      log_level      = var.edge_log_level
      prometheus_port = var.edge_prometheus_port
    }
    development = var.create_development_profile ? {
      id             = azurerm_iotoperations_dataflow_profile.development[0].id
      name           = azurerm_iotoperations_dataflow_profile.development[0].name
      instance_count = azurerm_iotoperations_dataflow_profile.development[0].instance_count
      log_level      = var.development_log_level
      prometheus_port = var.development_prometheus_port
    } : null
  }
}

# Metrics Endpoints Summary
output "metrics_endpoints" {
  description = "Prometheus metrics endpoints for all profiles"
  value = {
    high_performance = "http://localhost:${var.high_performance_prometheus_port}/metrics"
    standard        = "http://localhost:${var.standard_prometheus_port}/metrics"
    edge            = "http://localhost:${var.edge_prometheus_port}/metrics"
    development     = var.create_development_profile ? "http://localhost:${var.development_prometheus_port}/metrics" : null
  }
}