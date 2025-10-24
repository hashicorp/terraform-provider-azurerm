output "resource_group_id" {
  description = "The ID of the resource group"
  value       = azurerm_resource_group.example.id
}

output "instance_id" {
  description = "The ID of the IoT Operations instance"
  value       = azurerm_iotoperations_instance.example.id
}

output "broker_id" {
  description = "The ID of the IoT Operations broker"
  value       = azurerm_iotoperations_broker.example.id
}

output "broker_listener_id" {
  description = "The ID of the IoT Operations broker listener"
  value       = azurerm_iotoperations_broker_listener.example.id
}

output "broker_listener_name" {
  description = "The name of the IoT Operations broker listener"
  value       = azurerm_iotoperations_broker_listener.example.name
}

output "broker_listener_port" {
  description = "The port of the IoT Operations broker listener"
  value       = azurerm_iotoperations_broker_listener.example.port
}

output "broker_listener_service_name" {
  description = "The service name of the IoT Operations broker listener"
  value       = azurerm_iotoperations_broker_listener.example.service_name
}

output "broker_listener_service_type" {
  description = "The service type of the IoT Operations broker listener"
  value       = azurerm_iotoperations_broker_listener.example.service_type
}