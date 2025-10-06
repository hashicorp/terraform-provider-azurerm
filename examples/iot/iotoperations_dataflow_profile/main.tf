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

# High-performance dataflow profile for real-time processing
resource "azurerm_iotoperations_dataflow_profile" "high_performance" {
  name                = var.high_performance_profile_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  instance_count = var.high_performance_instance_count

  diagnostics {
    logs {
      level = var.high_performance_log_level
    }
    metrics {
      prometheus_port = var.high_performance_prometheus_port
    }
    self_check {
      mode                = var.high_performance_self_check_mode
      interval_seconds    = var.high_performance_self_check_interval
      timeout_seconds     = var.high_performance_self_check_timeout
    }
  }
}

# Standard dataflow profile for batch processing
resource "azurerm_iotoperations_dataflow_profile" "standard" {
  name                = var.standard_profile_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  instance_count = var.standard_instance_count

  diagnostics {
    logs {
      level = var.standard_log_level
    }
    metrics {
      prometheus_port = var.standard_prometheus_port
    }
    self_check {
      mode                = var.standard_self_check_mode
      interval_seconds    = var.standard_self_check_interval
      timeout_seconds     = var.standard_self_check_timeout
    }
  }
}

# Low-resource dataflow profile for edge scenarios
resource "azurerm_iotoperations_dataflow_profile" "edge" {
  name                = var.edge_profile_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  instance_count = var.edge_instance_count

  diagnostics {
    logs {
      level = var.edge_log_level
    }
    metrics {
      prometheus_port = var.edge_prometheus_port
    }
    self_check {
      mode                = var.edge_self_check_mode
      interval_seconds    = var.edge_self_check_interval
      timeout_seconds     = var.edge_self_check_timeout
    }
  }
}

# Development/testing dataflow profile
resource "azurerm_iotoperations_dataflow_profile" "development" {
  count = var.create_development_profile ? 1 : 0

  name                = var.development_profile_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  instance_count = var.development_instance_count

  diagnostics {
    logs {
      level = var.development_log_level
    }
    metrics {
      prometheus_port = var.development_prometheus_port
    }
    self_check {
      mode                = var.development_self_check_mode
      interval_seconds    = var.development_self_check_interval
      timeout_seconds     = var.development_self_check_timeout
    }
  }
}