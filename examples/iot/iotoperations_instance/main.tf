# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_version = ">= 1.6"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {}
}

# Use existing resource group
data "azurerm_resource_group" "example" {
  name = var.resource_group_name
}

# IoT Operations instance
resource "azurerm_iotoperations_instance" "example" {
  name                = "${var.prefix}-iotoperations"
  resource_group_name = data.azurerm_resource_group.example.name
  location            = data.azurerm_resource_group.example.location
  
  # Extended location (Custom Location for Arc-enabled Kubernetes)
  extended_location_name = var.custom_location_id
  extended_location_type = "CustomLocation"
  
  # Required schema registry reference
  schema_registry_ref = var.schema_registry_ref
  
  # Optional properties
  description = "IoT Operations instance created via Terraform"
}
