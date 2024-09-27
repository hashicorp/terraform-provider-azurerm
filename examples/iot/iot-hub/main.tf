# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

module "naming" {
  source  = "Azure/naming/azurerm"
  version = "0.4.0"
  prefix  = [var.prefix]
}

resource "azurerm_resource_group" "example" {
  name     = module.naming.resource_group.name
  location = "westeurope"
}

resource "azurerm_iothub" "example" {
  name                = module.naming.iothub.name
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku {
    name     = "F1"
    capacity = 1
  }
  local_authentication_enabled  = true
  event_hub_partition_count     = 2
  event_hub_retention_in_days   = 1
  public_network_access_enabled = false

  endpoint {
    name                       = "example_storage"
    type                       = "AzureIotHub.StorageContainer"
    resource_group_name        = azurerm_resource_group.example.name
    authentication_type        = "keyBased"
    connection_string          = var.storage_connection_string
    container_name             = "example_container"
    file_name_format           = "{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}"
    batch_frequency_in_seconds = 60
    encoding                   = "Avro"
    max_chunk_size_in_bytes    = 10485760
  }

  route {
    name           = "example_route"
    source         = "DeviceMessages"
    condition      = "$body.type = \"example\""
    enabled        = true
    endpoint_names = ["example_storage"]
  }

  enrichment {
    key            = "example_enrichment"
    value          = "$iothubname"
    endpoint_names = ["example_storage"]
  }

  fallback_route {
    source         = "DeviceMessages"
    enabled        = true
    endpoint_names = ["events"]
  }

  file_upload {
    authentication_type = "keyBased"
    container_name      = "example_container"
    connection_string   = var.storage_connection_string
    sas_ttl             = "PT1H"
    notifications       = false
    lock_duration       = "PT1M"
    default_ttl         = "PT1H"
    max_delivery_count  = 10
  }

  identity {
    type = "SystemAssigned"
  }

  cloud_to_device {
    default_ttl        = "PT1H"
    max_delivery_count = 10
    feedback {
      time_to_live       = "PT1H"
      lock_duration      = "PT60S"
      max_delivery_count = 10
    }
  }

  network_rule_set {
    default_action                     = "Deny"
    apply_to_builtin_eventhub_endpoint = true
    ip_rule {
      name    = "example_ip_rule"
      ip_mask = "10.0.0.0/24"
      action  = "Allow"
    }
  }

  tags = {
    environment = "Development"
    region      = azurerm_resource_group.example.location
  }
}
