terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

# IoT Operations Instance (prerequisite for broker)
resource "azurerm_iotoperations_instance" "example" {
  name                = "example-iot-instance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  # Extended location for IoT Operations
  extended_location {
    name = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.example.name}/providers/Microsoft.ExtendedLocation/customLocations/example-location"
    type = "CustomLocation"
  }
}

# IoT Operations Broker
resource "azurerm_iotoperations_broker" "example" {
  name                = "example-broker"
  resource_group_name = azurerm_resource_group.example.name
  instance_name       = azurerm_iotoperations_instance.example.name

  properties {
    memory_profile = "Medium"
  }

  extended_location {
    name = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.example.name}/providers/Microsoft.ExtendedLocation/customLocations/example-location"
    type = "CustomLocation"
  }
}

data "azurerm_client_config" "current" {}

output "broker_id" {
  description = "The ID of the IoT Operations Broker"
  value       = azurerm_iotoperations_broker.example.id
}

output "broker_name" {
  description = "The name of the IoT Operations Broker"
  value       = azurerm_iotoperations_broker.example.name
}

output "instance_name" {
  description = "The name of the IoT Operations Instance"
  value       = azurerm_iotoperations_instance.example.name
}
