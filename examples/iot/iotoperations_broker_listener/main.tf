terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
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


resource "azurerm_iotoperations_broker" "example" {
  name                = var.broker_name
  resource_group_name = azurerm_resource_group.example.name
  instance_name       = var.instance_name

  extended_location {
    name = var.custom_location_id
    type = "CustomLocation"
  }

  # Add other supported properties/blocks here as needed
}

# Simple Broker Listener - Basic port configuration only
resource "azurerm_iotoperations_broker_listener" "simple" {

  name                = var.simple_listener_name
  resource_group_name = azurerm_resource_group.example.name
  instance_name       = var.instance_name
  broker_name         = azurerm_iotoperations_broker.example.name

  extended_location_name = var.custom_location_id

  service_name = var.simple_service_name
  service_type = "ClusterIp"

  ports {
    port = var.simple_listener_port
  }
}

# Complex Broker Listener - With TLS, authentication, and authorization
resource "azurerm_iotoperations_broker_listener" "complex" {

  name                = var.complex_listener_name
  resource_group_name = azurerm_resource_group.example.name
  instance_name       = var.instance_name
  broker_name         = azurerm_iotoperations_broker.example.name

  extended_location_name = var.custom_location_id

  service_name = var.complex_service_name
  service_type = var.complex_service_type

  ports {
    port             = var.complex_listener_port
    node_port        = var.complex_node_port
  protocol         = "MQTT"
  authentication_ref = "example-auth-ref"
  authorization_ref  = "example-authz-ref"

    tls {
      mode = var.complex_tls_mode

      cert_manager_certificate_spec {
        duration    = var.complex_tls_cert_duration
        secret_name = var.complex_tls_cert_secret_name
        renew_before = var.complex_tls_cert_renew_before

        issuer_ref {
          name  = var.complex_tls_issuer_name
          kind  = var.complex_tls_issuer_kind
          group = var.complex_tls_issuer_group
        }

        private_key {
          algorithm       = "Rsa2048"
          rotation_policy = var.complex_tls_private_key_rotation_policy
        }

        san {
          dns = var.complex_tls_san_dns
          ip  = var.complex_tls_san_ip
        }
      }
    }
  }
}

# Full/Advanced Broker Listener - Multiple ports with all configuration options
resource "azurerm_iotoperations_broker_listener" "full" {
  count = var.enable_full_listener ? 1 : 0


  name                = var.full_listener_name
  resource_group_name = azurerm_resource_group.example.name
  instance_name       = var.instance_name
  broker_name         = azurerm_iotoperations_broker.example.name

  extended_location_name = var.custom_location_id

  service_name = var.full_service_name
  service_type = var.full_service_type

  # MQTT port with TLS
  ports {
    port             = 8883
    node_port        = 30883
  protocol         = "MQTT"
  authentication_ref = "example-auth-ref"
  authorization_ref  = "example-authz-ref"

    tls {
      mode = "Automatic"

      cert_manager_certificate_spec {
        duration    = "8760h"
        secret_name = "mqtt-tls-secret"
        renew_before = "720h"

        issuer_ref {
          name  = "cluster-issuer"
          kind  = "ClusterIssuer"
          group = "cert-manager.io"
        }

        private_key {
    algorithm       = "Rsa2048"
          rotation_policy = "Always"
        }

        san {
          dns = ["mqtt.example.com", "broker.local"]
          ip  = ["10.0.0.1"]
        }
      }
    }
  }

  # WebSocket port
  ports {
    port             = 8080
    node_port        = 30080
    protocol         = "WebSockets"
    authentication_ref = var.full_authentication_ref
    authorization_ref  = var.full_authorization_ref
  }

  # HTTP port for management
  ports {
    port      = 8081
    node_port = 30081
    protocol  = "Http"
  }
}