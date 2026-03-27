terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
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

  schema_registry_ref = var.schema_registry_ref
}

# MQTT Endpoint for IoT data ingestion
resource "azurerm_iotoperations_dataflow_endpoint" "mqtt" {
  name                = "mqtt-endpoint"
  resource_group_name = var.resource_group_name
  instance_name       = var.instance_name

  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }

  endpoint_type = "Mqtt"

  mqtt_settings {
    host = var.mqtt_host

    authentication {
      method = "SystemAssignedManagedIdentity"
    }
  }
}

# Azure Data Explorer Endpoint for processed data
resource "azurerm_iotoperations_dataflow_endpoint" "adx" {
  name                = "adx-endpoint"
  resource_group_name = var.resource_group_name
  instance_name       = var.instance_name

  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }

  endpoint_type = "DataExplorer"

  data_explorer_settings {
    host     = var.adx_host
    database = var.adx_database

    authentication {
      method = "SystemAssignedManagedIdentity"
    }
  }
}

# Azure Blob Storage Endpoint for data archival
resource "azurerm_iotoperations_dataflow_endpoint" "storage" {
  name                = "storage-endpoint"
  resource_group_name = var.resource_group_name
  instance_name       = var.instance_name

  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }

  endpoint_type = "DataLakeStorage"

  data_lake_storage_settings {
    host = var.storage_host

    authentication {
      method = "SystemAssignedManagedIdentity"
    }
  }
}

# Local Storage Endpoint for temporary data
resource "azurerm_iotoperations_dataflow_endpoint" "local" {
  name                = "local-endpoint"
  resource_group_name = var.resource_group_name
  instance_name       = var.instance_name

  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }

  endpoint_type = "LocalStorage"

  local_storage_settings {
    path = "/mnt/data"
  }
}

# Fabric OneLake Endpoint for analytics
resource "azurerm_iotoperations_dataflow_endpoint" "fabric" {
  name                = "fabric-endpoint"
  resource_group_name = var.resource_group_name
  instance_name       = var.instance_name

  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }

  endpoint_type = "FabricOneLake"

  fabric_one_lake_settings {
    host              = var.fabric_host
    names             = [var.fabric_lakehouse_name, var.fabric_workspace_name]
    one_lake_path_type = var.fabric_one_lake_path_type
    workspace         = var.fabric_workspace_name

    authentication {
      method = "SystemAssignedManagedIdentity"
    }
  }
}