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

  extended_location_name = var.custom_location_id
  extended_location_type = "CustomLocation"

  identity {
    type = "SystemAssigned"
  }

  tags = var.tags
}

resource "azurerm_iotoperations_dataflow_profile" "example" {
  name                = var.dataflow_profile_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  instance_count = var.dataflow_profile_instance_count

  diagnostics {
    logs {
      level = var.dataflow_profile_log_level
    }
    metrics {
      prometheus_port = var.dataflow_profile_prometheus_port
    }
  }
}

resource "azurerm_iotoperations_dataflow" "example" {
  name                      = var.dataflow_name
  dataflow_profile_id      = azurerm_iotoperations_dataflow_profile.example.id

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  mode = var.dataflow_mode

  dynamic "source" {
    for_each = var.dataflow_sources
    content {
      name           = source.value.name
      endpoint_ref   = source.value.endpoint_ref
      asset_ref      = source.value.asset_ref
      schema_ref     = source.value.schema_ref
      
      dynamic "serialization" {
        for_each = source.value.serialization != null ? [source.value.serialization] : []
        content {
          format = serialization.value.format
        }
      }
    }
  }

  dynamic "destination" {
    for_each = var.dataflow_destinations
    content {
      name         = destination.value.name
      endpoint_ref = destination.value.endpoint_ref
      schema_ref   = destination.value.schema_ref

      dynamic "serialization" {
        for_each = destination.value.serialization != null ? [destination.value.serialization] : []
        content {
          format = serialization.value.format
        }
      }
    }
  }

  dynamic "transformation" {
    for_each = var.dataflow_transformations
    content {
      type = transformation.value.type
      
      dynamic "filter" {
        for_each = transformation.value.filter != null ? [transformation.value.filter] : []
        content {
          expression = filter.value.expression
          type      = filter.value.type
        }
      }

      dynamic "map" {
        for_each = transformation.value.map != null ? [transformation.value.map] : []
        content {
          expression = map.value.expression
          type      = map.value.type
        }
      }
    }
  }

  dynamic "operation" {
    for_each = var.dataflow_operations
    content {
      name                = operation.value.name
      operation_type      = operation.value.operation_type
      source_name         = operation.value.source_name
      destination_name    = operation.value.destination_name
      
      dynamic "built_in_transformation" {
        for_each = operation.value.built_in_transformation != null ? [operation.value.built_in_transformation] : []
        content {
          serialize_type = built_in_transformation.value.serialize_type
          schema_ref     = built_in_transformation.value.schema_ref
        }
      }
    }
  }
}