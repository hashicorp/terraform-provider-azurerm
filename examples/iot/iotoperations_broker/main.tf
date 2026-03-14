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

# IoT Operations broker
resource "azurerm_iotoperations_broker" "example" {
  name                = var.broker_name
  resource_group_name = data.azurerm_resource_group.example.name
  instance_name       = var.instance_name
  
  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }
  
  properties {
    memory_profile = "Medium"
    
    cardinality {
      backend_chain {
        partitions        = 2
        redundancy_factor = 1
        workers          = 1
      }
      
      frontend {
        replicas = 2
        workers  = 1
      }
    }
    
    advanced {
      encrypt_internal_traffic = "Enabled"
      
      clients {
        max_session_expiry_seconds = 3600
        max_message_expiry_seconds = 3600
        max_packet_size_bytes      = 1048576
        max_receive_maximum        = 100
        max_keep_alive_seconds     = 3600
        
        subscriber_queue_limit {
          length   = 1000
          strategy = "DropOldest"
        }
      }
      
      internal_certs {
        duration     = "8760h"
        renew_before = "720h"
        
        private_key {
          algorithm       = "RSA"
          rotation_policy = "Always"
        }
      }
    }
    
    diagnostics {
      logs {
        level = "info"
      }
      
      metrics {
        prometheus_port = 9090
      }
      
      self_check {
        mode             = "Enabled"
        interval_seconds = 30
        timeout_seconds  = 15
      }
      
      traces {
        mode                 = "Enabled"
        cache_size_megabytes = 16
        span_channel_capacity = 1000
        
        self_tracing {
          mode             = "Enabled"
          interval_seconds = 30
        }
      }
    }
    
    disk_backed_message_buffer {
      max_size = "1Gi"
      
      ephemeral_volume_claim_spec {
        access_modes = ["ReadWriteOnce"]
        
        resources {
          requests = {
            "storage" = "1Gi"
          }
        }
      }
    }
    
    generate_resource_limits {
      cpu = "Enabled"
    }
  }
}