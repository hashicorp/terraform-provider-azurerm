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

resource "azurerm_iotoperations_dataflow_profile" "example" {
  name                = var.dataflow_profile_name
  resource_group_name = azurerm_resource_group.example.name
  instance_name       = azurerm_iotoperations_instance.example.name

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
  name                = var.dataflow_name
  resource_group_name = azurerm_resource_group.example.name
  instance_name       = azurerm_iotoperations_instance.example.name
  dataflow_profile_name = azurerm_iotoperations_dataflow_profile.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  mode = var.dataflow_mode

  operations {
    name           = "operation-transform"
    operation_type = "BuiltInTransformation"

    source_settings {
      data_sources         = ["source-mqtt"]
      asset_ref            = "temperature-asset"
      endpoint_ref         = "mqtt-endpoint"
      schema_ref           = "temperature-schema"
      serialization_format = "Json"
    }

    destination_settings {
      data_destination = "destination-adx"
      endpoint_ref     = "adx-endpoint"
    }

    built_in_transformation_settings {
      datasets {
        key    = "dataset1"
        inputs = ["input1"]
      }
      filter {
        expression = "temperature > 20"
        inputs     = ["input1"]
      }
      map {
        output = "output1"
        inputs = ["input1"]
      }
    }
  }
}