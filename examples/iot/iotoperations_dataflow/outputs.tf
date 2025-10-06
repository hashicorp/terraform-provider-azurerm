output "resource_group_id" {
  description = "The ID of the resource group"
  value       = azurerm_resource_group.example.id
}

output "instance_id" {
  description = "The ID of the IoT Operations instance"
  value       = azurerm_iotoperations_instance.example.id
}

output "dataflow_profile_id" {
  description = "The ID of the IoT Operations dataflow profile"
  value       = azurerm_iotoperations_dataflow_profile.example.id
}

output "dataflow_profile_name" {
  description = "The name of the IoT Operations dataflow profile"
  value       = azurerm_iotoperations_dataflow_profile.example.name
}

output "dataflow_id" {
  description = "The ID of the IoT Operations dataflow"
  value       = azurerm_iotoperations_dataflow.example.id
}

output "dataflow_name" {
  description = "The name of the IoT Operations dataflow"
  value       = azurerm_iotoperations_dataflow.example.name
}

output "dataflow_mode" {
  description = "The mode of the IoT Operations dataflow"
  value       = azurerm_iotoperations_dataflow.example.mode
}

output "dataflow_sources" {
  description = "The sources configured for the dataflow"
  value       = azurerm_iotoperations_dataflow.example.source
}

output "dataflow_destinations" {
  description = "The destinations configured for the dataflow"
  value       = azurerm_iotoperations_dataflow.example.destination
}