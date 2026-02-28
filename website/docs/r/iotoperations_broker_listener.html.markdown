---
subcategory: "IoT Operations"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotoperations_broker_listener"
description: |-
  Manages an Azure IoT Operations Broker Listener.
---

# azurerm_iotoperations_broker_listener

Manages an Azure IoT Operations Broker Listener.

A Broker Listener defines how clients can connect to the IoT Operations broker, including network configuration, authentication, and TLS settings. It provides network endpoints for MQTT and WebSocket protocols with configurable security options.

## Example Usage

### Basic Broker Listener

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iotoperations_instance" "example" {
  name                     = "example-instance"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  extended_location_name   = "microsoftiotoperations"
  extended_location_type   = "CustomLocation"
}

resource "azurerm_iotoperations_broker" "example" {
  name                     = "example-broker"
  resource_group_name      = azurerm_resource_group.example.name
  instance_name           = azurerm_iotoperations_instance.example.name
  extended_location_name  = azurerm_iotoperations_instance.example.extended_location_name
  extended_location_type  = azurerm_iotoperations_instance.example.extended_location_type
}

resource "azurerm_iotoperations_broker_listener" "example" {
  name                    = "example-listener"
  resource_group_name     = azurerm_resource_group.example.name
  instance_name          = azurerm_iotoperations_instance.example.name
  broker_name            = azurerm_iotoperations_broker.example.name
  extended_location_name = azurerm_iotoperations_instance.example.extended_location_name
  service_name           = "example-service"
  service_type           = "ClusterIp"

  ports {
    port     = 1883
    protocol = "MQTT"
  }

  ports {
    port     = 8883
    protocol = "MQTT"
    
    tls {
      mode = "Manual"
      
      manual {
        secret_ref = "example-tls-secret"
      }
    }
  }
}
```

### Broker Listener with Automatic TLS using Cert Manager

```hcl
resource "azurerm_iotoperations_broker_listener" "example_auto_tls" {
  name                    = "example-listener-auto-tls"
  resource_group_name     = azurerm_resource_group.example.name
  instance_name          = azurerm_iotoperations_instance.example.name
  broker_name            = azurerm_iotoperations_broker.example.name
  extended_location_name = azurerm_iotoperations_instance.example.extended_location_name
  service_name           = "example-service-tls"
  service_type           = "LoadBalancer"

  ports {
    port      = 8883
    node_port = 30883
    protocol  = "MQTT"
    
    tls {
      mode = "Automatic"
      
      cert_manager_certificate_spec {
        duration     = "8760h"  # 1 year
        secret_name  = "broker-tls-secret"
        renew_before = "720h"   # 30 days
        
        issuer_ref {
          group = "cert-manager.io"
          kind  = "ClusterIssuer"
          name  = "letsencrypt-prod"
        }
        
        private_key {
          algorithm       = "Rsa2048"
          rotation_policy = "Always"
        }
        
        san {
          dns = ["broker.example.com", "*.broker.example.com"]
          ip  = ["192.168.1.100", "10.0.0.100"]
        }
      }
    }
  }

  ports {
    port     = 9001
    protocol = "WebSockets"
    
    tls {
      mode = "Automatic"
      
      cert_manager_certificate_spec {
        issuer_ref {
          group = "cert-manager.io"
          kind  = "ClusterIssuer"
          name  = "letsencrypt-prod"
        }
        
        san {
          dns = ["websocket.example.com"]
          ip  = ["192.168.1.101"]
        }
      }
    }
  }
}
```

### Broker Listener with Authentication and Authorization

```hcl
resource "azurerm_iotoperations_broker_listener" "example_auth" {
  name                    = "example-listener-auth"
  resource_group_name     = azurerm_resource_group.example.name
  instance_name          = azurerm_iotoperations_instance.example.name
  broker_name            = azurerm_iotoperations_broker.example.name
  extended_location_name = azurerm_iotoperations_instance.example.extended_location_name
  service_name           = "auth-service"
  service_type           = "NodePort"

  ports {
    port               = 1883
    node_port          = 31883
    protocol           = "MQTT"
    authentication_ref = "example-auth-method"
    authorization_ref  = "example-authz-policy"
  }

  ports {
    port               = 8883
    node_port          = 31884
    protocol           = "MQTT"
    authentication_ref = "example-auth-method"
    authorization_ref  = "example-authz-policy"
    
    tls {
      mode = "Manual"
      
      manual {
        secret_ref = "mqtt-tls-secret"
      }
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Broker Listener. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the IoT Operations Broker Listener should exist. Changing this forces a new resource to be created.

* `instance_name` - (Required) The name of the IoT Operations Instance. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `broker_name` - (Required) The name of the IoT Operations Broker. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `extended_location_name` - (Required) The extended location name where the IoT Operations Broker Listener should be deployed. Changing this forces a new resource to be created.

* `ports` - (Required) A list of `ports` blocks as defined below. At least one port must be configured.

* `service_name` - (Optional) The name of the Kubernetes service. Must be between 1-63 characters.

* `service_type` - (Optional) The type of Kubernetes service. Possible values are `LoadBalancer`, `NodePort`, and `ClusterIp`. Defaults to `ClusterIp`.

---

A `ports` block supports the following:

* `port` - (Required) The port number for the broker listener. Must be between 1-65535.

* `node_port` - (Optional) The node port number when `service_type` is set to `NodePort` or `LoadBalancer`. Must be between 30000-32767.

* `protocol` - (Optional) The protocol for the port. Possible values are `MQTT` and `WebSockets`. Defaults to `MQTT`.

* `authentication_ref` - (Optional) Reference to the authentication method to use for this port. Must be between 1-253 characters.

* `authorization_ref` - (Optional) Reference to the authorization policy to use for this port. Must be between 1-253 characters.

* `tls` - (Optional) A `tls` block as defined below for configuring TLS settings.

---

A `tls` block supports the following:

* `mode` - (Required) The TLS mode for the port. Possible values are `Automatic` and `Manual`.

* `cert_manager_certificate_spec` - (Optional) A `cert_manager_certificate_spec` block as defined below. This is required when `mode` is set to `Automatic`.

* `manual` - (Optional) A `manual` block as defined below. This is required when `mode` is set to `Manual`.

---

A `cert_manager_certificate_spec` block supports the following:

* `issuer_ref` - (Required) An `issuer_ref` block as defined below specifying the certificate issuer.

* `duration` - (Optional) The duration of the certificate validity period. Must be between 1-50 characters (e.g., "8760h" for 1 year).

* `secret_name` - (Optional) The name of the Kubernetes secret where the certificate will be stored. Must be between 1-253 characters.

* `renew_before` - (Optional) The time before expiry when the certificate should be renewed. Must be between 1-50 characters (e.g., "720h" for 30 days).

* `private_key` - (Optional) A `private_key` block as defined below for private key configuration.

* `san` - (Optional) A `san` block as defined below for Subject Alternative Names.

---

An `issuer_ref` block supports the following:

* `group` - (Required) The API group of the certificate issuer. Must be between 1-253 characters (typically "cert-manager.io").

* `kind` - (Required) The kind of certificate issuer. Possible values are `ClusterIssuer` and `Issuer`.

* `name` - (Required) The name of the certificate issuer. Must be between 1-253 characters.

---

A `private_key` block supports the following:

* `algorithm` - (Required) The algorithm for the private key. Possible values are `Rsa2048`, `Rsa4096`, `Rsa8192`, `Ec256`, `Ec384`, `Ec521`, and `Ed25519`.

* `rotation_policy` - (Required) The rotation policy for the private key. Possible values are `Always` and `Never`.

---

A `san` block supports the following:

* `dns` - (Required) A list of DNS names to include in the certificate's Subject Alternative Names. Each DNS name must be between 1-253 characters.

* `ip` - (Required) A list of IP addresses to include in the certificate's Subject Alternative Names. Each IP must be a valid IP address.

---

A `manual` block supports the following:

* `secret_ref` - (Required) Reference to the Kubernetes secret containing the TLS certificate. Must be between 1-253 characters.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ID of the IoT Operations Broker Listener.

* `provisioning_state` - The provisioning state of the Broker Listener.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Operations Broker Listener.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Operations Broker Listener.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Operations Broker Listener.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Operations Broker Listener.

## Import

An IoT Operations Broker Listener can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotoperations_broker_listener.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.IoTOperations/instances/instance1/brokers/broker1/listeners/listener1
```
