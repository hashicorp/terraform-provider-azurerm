# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_version = ">= 1.6"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
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

# IoT Operations broker
resource "azurerm_iotoperations_broker" "example" {
  name                = "${var.prefix}-broker"
  resource_group_name = data.azurerm_resource_group.example.name
  instance_name       = var.instance_name
  location            = data.azurerm_resource_group.example.location
  
  # Extended location (Custom Location for Arc-enabled Kubernetes)
  extended_location_name = var.custom_location_id
  extended_location_type = "CustomLocation"
  
  # Optional properties
  properties {
    description = "IoT Operations broker created via Terraform"
  }
}