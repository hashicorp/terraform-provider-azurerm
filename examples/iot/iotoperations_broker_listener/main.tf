terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_iotoperations_instance" "example" {
  name                = var.instance_name
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }

  tags = var.tags
}

resource "azurerm_iotoperations_broker" "example" {
  name                  = var.broker_name
  iot_operations_instance_id = azurerm_iotoperations_instance.example.id

  cardinality {
    backend = {
      partitions   = 2
      redundancy_factor = 1
      workers      = 1
    }
    frontend = {
      replicas = 2
      workers  = 1
    }
  }

  memory_profile = "Medium"
}

resource "azurerm_iotoperations_broker_listener" "example" {
  name                = var.listener_name
  broker_id          = azurerm_iotoperations_broker.example.id

  port = var.listener_port

  service_name = var.service_name
  service_type = var.service_type

  dynamic "tls" {
    for_each = var.enable_tls ? [1] : []
    content {
      mode = var.tls_mode
      cert_manager_certificate_spec {
        duration           = var.tls_cert_duration
        issue_ref {
          name  = var.tls_issuer_name
          kind  = var.tls_issuer_kind
          group = var.tls_issuer_group
        }
        renew_before       = var.tls_cert_renew_before
        secret_name        = var.tls_cert_secret_name
        subject {
          organization       = [var.tls_subject_organization]
          organizational_unit = [var.tls_subject_organizational_unit]
        }
      }
    }
  }

  dynamic "authentication_ref" {
    for_each = var.authentication_ref_name != "" ? [1] : []
    content {
      name = var.authentication_ref_name
    }
  }

  dynamic "authorization_ref" {
    for_each = var.authorization_ref_name != "" ? [1] : []
    content {
      name = var.authorization_ref_name
    }
  }
}