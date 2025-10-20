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
  
  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }
  
  authentication_methods {
    method = "Custom"
    
    custom_settings {
      endpoint           = "https://www.example.com"
      ca_cert_config_map = "pdecudefqyolvncbus"
      
      headers = {
        "key8518" = "bwityjy"
      }
      
      auth {
        x509 {
          secret_ref = "secret-name"
        }
      }
    }
    
    service_account_token_settings {
      audiences = ["jqyhyqatuydg"]
    }
    
    x509_settings {
      trusted_client_ca_cert = "vlctsqddl"
      
      authorization_attributes = {
        "key3384" = {
          subject = "jpgwctfeixitptfgfnqhua"
          attributes = {
            "key186" = "ucpajramsz"
          }
        }
      }
    }
  }
}