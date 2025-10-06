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
  
  authorization_policies {
    cache = "Enabled"
    rules {
      brokers_resources = ["*"]
      method            = "Connect"
      clients           = ["*"]
      state_store_key_value_pairs {
        key   = "group"
        value = "admin"
      }
    }
  }
}