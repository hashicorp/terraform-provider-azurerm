output "resource_group_id" {
  description = "The ID of the resource group"
  value       = azurerm_resource_group.example.id
}


output "broker_id" {
  description = "The ID of the IoT Operations broker"
  value       = azurerm_iotoperations_broker.example.id
}

