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
  subscription_id = "d4ccd08b-0809-446d-a8b7-7af8a90109cd"

}

resource "azurerm_resource_group" "example" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_iotoperations_instance" "example" {
  name                = var.instance_name
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  extended_location_name = var.custom_location_id
  extended_location_type = "CustomLocation"

  schema_registry_ref = var.schema_registry_ref

  tags = var.tags
}

# High-performance dataflow profile for real-time processing
resource "azurerm_iotoperations_dataflow_profile" "high_performance" {
  name                = var.high_performance_profile_name
  resource_group_name = var.resource_group_name
  instance_name       = var.instance_name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  instance_count = var.high_performance_instance_count

  diagnostics {
    logs {
      level = var.high_performance_log_level
    }
    metrics {
      prometheus_port = var.high_performance_prometheus_port
    }
  }
}

# Standard dataflow profile for batch processing
resource "azurerm_iotoperations_dataflow_profile" "standard" {
  name                = var.standard_profile_name
  resource_group_name = var.resource_group_name
  instance_name       = var.instance_name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  instance_count = var.standard_instance_count

  diagnostics {
    logs {
      level = var.standard_log_level
    }
    metrics {
      prometheus_port = var.standard_prometheus_port
    }
  }
}

# Low-resource dataflow profile for edge scenarios
resource "azurerm_iotoperations_dataflow_profile" "edge" {
  name                = var.edge_profile_name
  resource_group_name = var.resource_group_name
  instance_name       = var.instance_name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  instance_count = var.edge_instance_count

  diagnostics {
    logs {
      level = var.edge_log_level
    }
    metrics {
      prometheus_port = var.edge_prometheus_port
    }
  }
}

# Development/testing dataflow profile
resource "azurerm_iotoperations_dataflow_profile" "development" {
  count = var.create_development_profile ? 1 : 0

  name                = var.development_profile_name
  resource_group_name = var.resource_group_name
  instance_name       = var.instance_name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  instance_count = var.development_instance_count

  diagnostics {
    logs {
      level = var.development_log_level
    }
    metrics {
      prometheus_port = var.development_prometheus_port
    }
  }
}