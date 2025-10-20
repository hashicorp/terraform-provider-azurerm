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

# IoT Operations broker authorization
resource "azurerm_iotoperations_broker_authorization" "example" {
  name                = "${var.prefix}-broker-authz"
  resource_group_name = data.azurerm_resource_group.example.name
  instance_name       = var.instance_name
  broker_name         = var.broker_name
  
  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }
  
  authorization_policies {
    cache = "Enabled"
    
    rules {
      principals {
        clients = ["my-client-id"]
        attributes = [
          {
            "floor" = "floor1"
            "site"  = "site1"
          }
        ]
      }
      
      broker_resources {
        method = "Connect"
      }
      
      broker_resources {
        method = "Subscribe"
        topics = ["topic", "topic/with/wildcard/#"]
      }
      
      state_store_resources {
        method   = "ReadWrite"
        key_type = "Pattern"
        keys     = ["*"]
      }
    }
  }
}