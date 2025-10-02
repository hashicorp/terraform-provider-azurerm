terraform {
  required_version = ">= 1.6"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
    azapi = {
      source  = "Azure/azapi"
      version = "~> 2.0"
    }
  }
}

# ---------- Providers ----------
provider "azurerm" {
  features {}
}

provider "azapi" {}

# variables needed
variable "resource_group_name" {
  type        = string
  default     = "example-iotoperations"
  description = "Resource Group name to house the IoT Operations instance."
}

variable "location" {
  type        = string
  default     = "West Europe"
  description = "Azure location/region."
}

variable "instance_name" {
  type        = string
  default     = "terraforminstancecreated"
  description = "Name of the IoT Operations instance resource."
}

variable "broker_name" {
  type        = string
  default     = "terraformbroker"
  description = "Name of the IoT Operations broker resource."
}

variable "authorization_name" {
  type        = string
  default     = "terraformauthorization"
  description = "Name of the IoT Operations broker authorization resource."
}

# FULL ARM ID of your Custom Location (Arc-enabled K8s)
# /subscriptions/<subId>/resourceGroups/<rg>/providers/Microsoft.ExtendedLocation/customLocations/<customLocationName>
variable "custom_location_id" {
  type        = string
  description = "ARM ID of the Custom Location used by AIO."
}

# FULL ARM ID of your Device Registry Schema Registry:
# /subscriptions/<subId>/resourceGroups/<rg>/providers/Microsoft.DeviceRegistry/schemaRegistries/<schemaRegistryName>
variable "schema_registry_id" {
  type        = string
  description = "ARM ID of the Schema Registry referenced by the AIO instance."
}

# Optional: tags
variable "tags" {
  type = map(string)
  default = {
    Environment = "Dev"
    Owner       = "team"
  }
}

# -Data (useful for subscription)
data "azurerm_client_config" "current" {}

# Resources
resource "azurerm_resource_group" "example" {
  name     = var.resource_group_name
  location = var.location
}

# IoT Operations Instance via AzAPI (ARM: Microsoft.IoTOperations/instances@2024-11-01)
resource "azapi_resource" "iotops_instance" {
  type      = "Microsoft.IoTOperations/instances@2024-11-01"
  name      = var.instance_name
  parent_id = azurerm_resource_group.example.id
  location  = azurerm_resource_group.example.location
  tags      = var.tags

  # NOTE: For AzAPI v2, body must be object
  body = {
    # REQUIRED: AIO deploys into a Custom Location (Arc-enabled K8s)
    extendedLocation = {
      type = "CustomLocation"
      name = var.custom_location_id
    }

    properties = {
      # 'version' is READ-ONLY in this API version; do NOT include it.
      description = "Example IoT Operations instance managed by Terraform"

      # REQUIRED: reference an existing Schema Registry by ARM ID
      schemaRegistryRef = {
        resourceId = var.schema_registry_id
      }
    }
  }
}

# IoT Operations Broker via AzAPI (ARM: Microsoft.IoTOperations/instances/brokers@2024-11-01)
resource "azapi_resource" "iotops_broker" {
  type      = "Microsoft.IoTOperations/instances/brokers@2024-11-01"
  name      = var.broker_name
  parent_id = azapi_resource.iotops_instance.id

  body = {
    extendedLocation = {
      type = "CustomLocation"
      name = var.custom_location_id
    }

    properties = {
      memoryProfile = "Medium"
    }
  }

  depends_on = [azapi_resource.iotops_instance]
}

# IoT Operations Broker Authorization
resource "azurerm_broker_authorization" "example" {
  name                = var.authorization_name
  resource_group_name = azurerm_resource_group.example.name
  instance_name       = azapi_resource.iotops_instance.name
  broker_name         = azapi_resource.iotops_broker.name

  authorization_policies {
    type = "RoleBasedAccessControl"
    
    rbac_config {
      principals = [
        {
          usernames = ["user1", "user2"]
          attributes = {
            "building" = "17"
            "floor"    = "1"
          }
        }
      ]

      resources = [
        {
          method = "Connect"
          topics = ["temperature/*", "humidity/*"]
        },
        {
          method = "Publish"
          topics = ["alerts/*"]
        },
        {
          method = "Subscribe"
          topics = ["commands/*"]
        }
      ]
    }
  }

  authorization_policies {
    type = "AttributeBasedAccessControl"
    
    abac_config {
      rules = [
        {
          principals = {
            attributes = {
              "department" = "engineering"
              "clearance"  = "high"
            }
          }
          
          resources = {
            topics = ["sensitive/*"]
            clientIds = ["trusted-client-*"]
          }
          
          operations = ["Connect", "Publish", "Subscribe"]
        }
      ]
    }
  }

  # Extended location for Arc-enabled Kubernetes
  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }

  depends_on = [azapi_resource.iotops_broker]
}

# ---------- Outputs ----------
output "iotoperations_instance_id" {
  value       = azapi_resource.iotops_instance.id
  description = "ARM ID of the created IoT Operations instance."
}

output "iotoperations_broker_id" {
  value       = azapi_resource.iotops_broker.id
  description = "ARM ID of the created IoT Operations broker."
}

output "broker_authorization_id" {
  value       = azurerm_broker_authorization.example.id
  description = "ARM ID of the created IoT Operations broker authorization."
}

output "broker_authorization_name" {
  value       = azurerm_broker_authorization.example.name
  description = "Name of the IoT Operations broker authorization."
}

output "resource_hierarchy" {
  value = {
    subscription_id    = data.azurerm_client_config.current.subscription_id
    resource_group     = azurerm_resource_group.example.name
    instance_name      = azapi_resource.iotops_instance.name
    broker_name        = azapi_resource.iotops_broker.name
    authorization_name = azurerm_broker_authorization.example.name
  }
  description = "Complete resource hierarchy for the IoT Operations broker authorization."
}
