output "resource_group_id" {
  description = "The ID of the resource group"
  value       = azurerm_resource_group.example.id
}

output "instance_id" {
  description = "The ID of the IoT Operations instance"
  value       = azurerm_iotoperations_instance.example.id
}

# MQTT Endpoint Outputs
output "mqtt_endpoint_id" {
  description = "The ID of the MQTT dataflow endpoint"
  value       = azurerm_iotoperations_dataflow_endpoint.mqtt.id
}

output "mqtt_endpoint_name" {
  description = "The name of the MQTT dataflow endpoint"
  value       = azurerm_iotoperations_dataflow_endpoint.mqtt.name
}

# Azure Data Explorer Endpoint Outputs
output "adx_endpoint_id" {
  description = "The ID of the Azure Data Explorer dataflow endpoint"
  value       = azurerm_iotoperations_dataflow_endpoint.adx.id
}

output "adx_endpoint_name" {
  description = "The name of the Azure Data Explorer dataflow endpoint"
  value       = azurerm_iotoperations_dataflow_endpoint.adx.name
}

# Azure Storage Endpoint Outputs
output "storage_endpoint_id" {
  description = "The ID of the Azure Storage dataflow endpoint"
  value       = azurerm_iotoperations_dataflow_endpoint.storage.id
}

output "storage_endpoint_name" {
  description = "The name of the Azure Storage dataflow endpoint"
  value       = azurerm_iotoperations_dataflow_endpoint.storage.name
}

# Local Storage Endpoint Outputs
output "local_endpoint_id" {
  description = "The ID of the local storage dataflow endpoint"
  value       = azurerm_iotoperations_dataflow_endpoint.local.id
}

output "local_endpoint_name" {
  description = "The name of the local storage dataflow endpoint"
  value       = azurerm_iotoperations_dataflow_endpoint.local.name
}

# Fabric OneLake Endpoint Outputs
output "fabric_endpoint_id" {
  description = "The ID of the Fabric OneLake dataflow endpoint"
  value       = var.enable_fabric_endpoint ? azurerm_iotoperations_dataflow_endpoint.fabric[0].id : null
}

output "fabric_endpoint_name" {
  description = "The name of the Fabric OneLake dataflow endpoint"
  value       = var.enable_fabric_endpoint ? azurerm_iotoperations_dataflow_endpoint.fabric[0].name : null
}

# Endpoint Configuration Summary
output "endpoint_summary" {
  description = "Summary of all configured dataflow endpoints"
  value = {
    mqtt_endpoint = {
      id   = azurerm_iotoperations_dataflow_endpoint.mqtt.id
      name = azurerm_iotoperations_dataflow_endpoint.mqtt.name
      type = azurerm_iotoperations_dataflow_endpoint.mqtt.endpoint_type
      host = var.mqtt_host
      port = var.mqtt_port
    }
    adx_endpoint = {
      id   = azurerm_iotoperations_dataflow_endpoint.adx.id
      name = azurerm_iotoperations_dataflow_endpoint.adx.name
      type = azurerm_iotoperations_dataflow_endpoint.adx.endpoint_type
      host = var.adx_cluster_uri
      database = var.adx_database_name
    }
    storage_endpoint = {
      id   = azurerm_iotoperations_dataflow_endpoint.storage.id
      name = azurerm_iotoperations_dataflow_endpoint.storage.name
      type = azurerm_iotoperations_dataflow_endpoint.storage.endpoint_type
      host = var.storage_account_host
      container = var.storage_container_name
    }
    local_endpoint = {
      id   = azurerm_iotoperations_dataflow_endpoint.local.id
      name = azurerm_iotoperations_dataflow_endpoint.local.name
      type = azurerm_iotoperations_dataflow_endpoint.local.endpoint_type
      pvc  = var.local_storage_pvc_name
    }
  }
}