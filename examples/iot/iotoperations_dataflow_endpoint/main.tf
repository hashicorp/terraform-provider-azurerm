terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_iotoperations_instance" "example" {
  name                = var.instance_name
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }

  tags = var.tags
}

# MQTT Endpoint for IoT data ingestion
resource "azurerm_iotoperations_dataflow_endpoint" "mqtt" {
  name                = var.mqtt_endpoint_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  endpoint_type = "Mqtt"

  mqtt_settings {
    host = var.mqtt_host
    port = var.mqtt_port
    
    dynamic "authentication" {
      for_each = var.mqtt_authentication_enabled ? [1] : []
      content {
        method = var.mqtt_auth_method
        username = var.mqtt_username
        password_secret_name = var.mqtt_password_secret
      }
    }

    dynamic "tls" {
      for_each = var.mqtt_tls_enabled ? [1] : []
      content {
        mode = var.mqtt_tls_mode
        trusted_ca_certificate_config_map = var.mqtt_tls_ca_cert_config_map
      }
    }

    keep_alive_seconds = var.mqtt_keep_alive_seconds
    retain             = var.mqtt_retain
    session_expiry_seconds = var.mqtt_session_expiry_seconds
    quality_of_service = var.mqtt_qos
  }
}

# Azure Data Explorer Endpoint for processed data
resource "azurerm_iotoperations_dataflow_endpoint" "adx" {
  name                = var.adx_endpoint_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  endpoint_type = "DataExplorer"

  data_explorer_settings {
    host     = var.adx_cluster_uri
    database = var.adx_database_name

    authentication {
      method                = "SystemAssignedManagedIdentity"
      system_assigned_managed_identity_audience = var.adx_audience
    }

    batching {
      latency_seconds = var.adx_batching_latency
      max_messages    = var.adx_batching_max_messages
    }
  }
}

# Azure Blob Storage Endpoint for data archival
resource "azurerm_iotoperations_dataflow_endpoint" "storage" {
  name                = var.storage_endpoint_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  endpoint_type = "DataLakeStorage"

  data_lake_storage_settings {
    host           = var.storage_account_host
    container_name = var.storage_container_name

    authentication {
      method                = "SystemAssignedManagedIdentity"
      system_assigned_managed_identity_audience = var.storage_audience
    }

    batching {
      latency_seconds = var.storage_batching_latency
      max_messages    = var.storage_batching_max_messages
    }
  }
}

# Local Storage Endpoint for temporary data
resource "azurerm_iotoperations_dataflow_endpoint" "local" {
  name                = var.local_endpoint_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  endpoint_type = "LocalStorage"

  local_storage_settings {
    persistent_volume_claim_ref = var.local_storage_pvc_name
  }
}

# Fabric OneLake Endpoint for analytics
resource "azurerm_iotoperations_dataflow_endpoint" "fabric" {
  count = var.enable_fabric_endpoint ? 1 : 0

  name                = var.fabric_endpoint_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  endpoint_type = "FabricOneLake"

  fabric_one_lake_settings {
    host             = var.fabric_host
    workspace_id     = var.fabric_workspace_id
    lakehouse_name   = var.fabric_lakehouse_name

    authentication {
      method                = "SystemAssignedManagedIdentity"
      system_assigned_managed_identity_audience = var.fabric_audience
    }

    batching {
      latency_seconds = var.fabric_batching_latency
      max_messages    = var.fabric_batching_max_messages
    }
  }
}