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

# IoT Operations Instance (prerequisite)
resource "azurerm_iotoperations_instance" "example" {
  name                = "example-iot-instance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  extended_location {
    name = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.example.name}/providers/Microsoft.ExtendedLocation/customLocations/example-location"
    type = "CustomLocation"
  }
}

# IoT Operations Broker (prerequisite)
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

# IoT Operations Broker Authentication
resource "azurerm_iotoperations_broker_authentication" "example" {
  name                = "example-auth"
  resource_group_name = azurerm_resource_group.example.name
  instance_name       = azurerm_iotoperations_instance.example.name
  broker_name         = azurerm_iotoperations_broker.example.name

  authentication_methods {
    method = "ServiceAccountToken"
    
    service_account_token_settings {
      audiences = ["audience1", "audience2"]
    }
  }

  authentication_methods {
    method = "X509"
    
    x509_settings {
      trusted_client_ca_cert = "example-ca-cert"
      authorization_attributes = {
        "building" = "17"
        "floor"    = "1"
      }
    }
  }

  authentication_methods {
    method = "Custom"
    
    custom_settings {
      auth {
        x509 {
          secret_ref = "example-secret"
        }
      }
      
      ca_cert_config_map = "example-ca-configmap"
      endpoint           = "https://example-auth-endpoint.com"
      
      headers = {
        "X-Custom-Header" = "example-value"
        "Authorization"   = "Bearer token"
      }
    }
  }
}

data "azurerm_client_config" "current" {}

output "broker_authentication_id" {
  description = "The ID of the IoT Operations Broker Authentication"
  value       = azurerm_iotoperations_broker_authentication.example.id
}

output "broker_authentication_name" {
  description = "The name of the IoT Operations Broker Authentication"
  value       = azurerm_iotoperations_broker_authentication.example.name
}

output "instance_name" {
  description = "The name of the IoT Operations Instance"
  value       = azurerm_iotoperations_instance.example.name
}

output "broker_name" {
  description = "The name of the IoT Operations Broker"
  value       = azurerm_iotoperations_broker.example.name
}
