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

# IoT Operations broker authentication
resource "azurerm_iotoperations_broker_authentication" "example" {
  name                = "${var.prefix}-broker-auth"
  resource_group_name = data.azurerm_resource_group.example.name
  instance_name       = var.instance_name
  broker_name         = var.broker_name
  
  authentication_methods {
    method = "ServiceAccountToken"
    
    custom_settings {
      auth {
        audience = var.audience
      }
    }
  }
}