terraform {
  required_version = ">= 1.6"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

# ---------- Providers ----------
provider "azurerm" {
  features {}
}

# variables needed
variable "resource_group_name" {
  type        = string
  default     = "terraformiotest"
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
  tags     = var.tags
}

# IoT Operations Instance via AzureRM (ARM template: Microsoft.IoTOperations/instances@2024-11-01)
resource "azurerm_resource_group_template_deployment" "iotops_instance" {
  name                = "${var.instance_name}-deploy"
  resource_group_name = azurerm_resource_group.example.name
  deployment_mode     = "Incremental"

  template_content = <<JSON
{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "instanceName":     { "type": "string" },
    "location":         { "type": "string" },
    "customLocationId": { "type": "string" },
    "schemaRegistryId": { "type": "string" },
    "tags":             { "type": "object" }
  },
  "resources": [
    {
      "type": "Microsoft.IoTOperations/instances",
      "apiVersion": "2024-11-01",
      "name": "[parameters('instanceName')]",
      "location": "[parameters('location')]",
      "extendedLocation": {
        "type": "CustomLocation",
        "name": "[parameters('customLocationId')]"
      },
      "tags": "[parameters('tags')]",
      "properties": {
        "description": "IoT Operations instance created via AzureRM template deployment",
        "schemaRegistryRef": {
          "resourceId": "[parameters('schemaRegistryId')]"
        }
      }
    }
  ],
  "outputs": {
    "instanceId": {
      "type": "string",
      "value": "[resourceId('Microsoft.IoTOperations/instances', parameters('instanceName'))]"
    }
  }
}
JSON

  parameters_content = jsonencode({
    instanceName     = { value = var.instance_name }
    location         = { value = var.location }
    customLocationId = { value = var.custom_location_id }
    schemaRegistryId = { value = var.schema_registry_id }
    tags             = { value = var.tags }
  })
}

# ---------- Outputs ----------
output "iotoperations_instance_id" {
  value       = jsondecode(azurerm_resource_group_template_deployment.iotops_instance.output_content).instanceId.value
  description = "ARM ID of the created IoT Operations instance."
}
