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

  # If you need to bypass AzAPI's input schema validation while wiring IDs, you can temporarily set:
  # schema_validation_enabled = false
}

# ---------- Outputs ----------
output "iotoperations_instance_id" {
  value       = azapi_resource.iotops_instance.id
  description = "ARM ID of the created IoT Operations instance."
}