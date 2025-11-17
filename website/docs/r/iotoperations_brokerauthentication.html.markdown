---
subcategory: "IoT Operations"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotoperations_broker_authentication"
description: |-
  Manages an Azure IoT Operations Broker Authentication.
---

# azurerm_iotoperations_broker_authentication

Manages an Azure IoT Operations Broker Authentication.

A Broker Authentication defines how clients authenticate with the IoT Operations broker. It supports multiple authentication methods including X.509 certificates, service account tokens, and custom authentication providers, providing flexible security options for different client types and deployment scenarios.

## Example Usage

### X.509 Certificate Authentication

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iotoperations_instance" "example" {
  name                   = "example-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  extended_location_name = "microsoftiotoperations"
  extended_location_type = "CustomLocation"
}

resource "azurerm_iotoperations_broker" "example" {
  name                = "example-broker"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = azurerm_iotoperations_instance.example.extended_location_type
  }
}

resource "azurerm_iotoperations_broker_authentication" "x509_auth" {
  name                = "x509-auth"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  authentication_methods {
    method = "X509"

    x509_settings {
      trusted_client_ca_cert = "ca-certificates"

      authorization_attributes {
        name    = "device-id"
        subject = "CN"
        attributes = {
          "device-type" = "sensor"
          "location"    = "factory-floor"
        }
      }

      authorization_attributes {
        name    = "organization"
        subject = "O"
        attributes = {
          "department" = "manufacturing"
          "access"     = "read-write"
        }
      }
    }
  }
}
```

### Service Account Token Authentication

```hcl
resource "azurerm_iotoperations_broker_authentication" "sat_auth" {
  name                = "sat-auth"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  authentication_methods {
    method = "ServiceAccountToken"

    service_account_token_settings {
      audiences = [
        "iotoperations.azure.com",
        "mqtt.broker.local",
        "https://management.azure.com/"
      ]
    }
  }
}
```

### Custom Authentication Provider

```hcl
resource "azurerm_iotoperations_broker_authentication" "custom_auth" {
  name                = "custom-auth"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  authentication_methods {
    method = "Custom"

    custom_settings {
      endpoint = "https://auth.example.com/validate"

      auth {
        x509 {
          secret_ref = "auth-server-cert"
        }
      }

      ca_cert_config_map = "custom-ca-certs"

      headers = {
        "Content-Type"    = "application/json"
        "X-API-Version"   = "v1"
        "X-Client-ID"     = "iot-broker"
      }
    }
  }
}
```

### Multiple Authentication Methods

```hcl
resource "azurerm_iotoperations_broker_authentication" "multi_auth" {
  name                = "multi-auth"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  # X.509 certificate authentication for devices
  authentication_methods {
    method = "X509"

    x509_settings {
      trusted_client_ca_cert = "device-ca-certs"

      authorization_attributes {
        name    = "device-type"
        subject = "CN"
        attributes = {
          "category" = "iot-device"
          "protocol" = "mqtt"
        }
      }
    }
  }

  # Service account token for applications
  authentication_methods {
    method = "ServiceAccountToken"

    service_account_token_settings {
      audiences = [
        "iotoperations.azure.com",
        "application.local"
      ]
    }
  }

  # Custom authentication for special clients
  authentication_methods {
    method = "Custom"

    custom_settings {
      endpoint           = "https://auth.enterprise.com/oauth/validate"
      ca_cert_config_map = "enterprise-ca"

      auth {
        x509 {
          secret_ref = "oauth-client-cert"
        }
      }

      headers = {
        "Authorization" = "Bearer ${var.api_token}"
        "Accept"       = "application/json"
      }
    }
  }
}
```

### Enterprise X.509 Authentication with Multiple Attributes

```hcl
resource "azurerm_iotoperations_broker_authentication" "enterprise_x509" {
  name                = "enterprise-x509"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  authentication_methods {
    method = "X509"

    x509_settings {
      trusted_client_ca_cert = "enterprise-root-ca"

      # Device identification
      authorization_attributes {
        name    = "device-identity"
        subject = "CN"
        attributes = {
          "device-id"     = "unique-identifier"
          "manufacturer"  = "company-name"
          "model"        = "device-model"
          "firmware"     = "version-info"
        }
      }

      # Location-based access
      authorization_attributes {
        name    = "location-access"
        subject = "L"
        attributes = {
          "facility"    = "building-name"
          "floor"       = "level-number"
          "zone"        = "security-zone"
          "access-level" = "restricted"
        }
      }

      # Organizational unit
      authorization_attributes {
        name    = "department"
        subject = "OU"
        attributes = {
          "division"    = "business-unit"
          "team"        = "operational-group"
          "role"        = "function-type"
          "clearance"   = "security-level"
        }
      }

      # Organization details
      authorization_attributes {
        name    = "organization"
        subject = "O"
        attributes = {
          "company"     = "corporation-name"
          "subsidiary"  = "regional-office"
          "contract"    = "agreement-type"
          "tier"        = "service-level"
        }
      }
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Broker Authentication. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the IoT Operations Broker Authentication should exist. Must be between 1-90 characters. Changing this forces a new resource to be created.

* `instance_name` - (Required) The name of the IoT Operations Instance. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `broker_name` - (Required) The name of the IoT Operations Broker. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `extended_location` - (Required) An `extended_location` block as defined below. Changing this forces a new resource to be created.

* `authentication_methods` - (Required) A list of `authentication_methods` blocks as defined below. At least one authentication method must be configured.

---

An `extended_location` block supports the following:

* `name` - (Required) The extended location name where the IoT Operations Broker Authentication should be deployed.

* `type` - (Optional) The extended location type. Defaults to `CustomLocation`.

---

An `authentication_methods` block supports the following:

* `method` - (Required) The authentication method type. Possible values are `Custom`, `ServiceAccountToken`, and `X509`.

* `custom_settings` - (Optional) A `custom_settings` block as defined below. Required when `method` is set to `Custom`.

* `service_account_token_settings` - (Optional) A `service_account_token_settings` block as defined below. Required when `method` is set to `ServiceAccountToken`.

* `x509_settings` - (Optional) An `x509_settings` block as defined below. Required when `method` is set to `X509`.

---

A `custom_settings` block supports the following:

* `endpoint` - (Required) The URL of the custom authentication endpoint.

* `auth` - (Optional) An `auth` block as defined below for authenticating with the custom endpoint.

* `ca_cert_config_map` - (Optional) Name of the ConfigMap containing CA certificates for verifying the custom authentication endpoint.

* `headers` - (Optional) A map of HTTP headers to include in requests to the custom authentication endpoint.

---

An `auth` block supports the following:

* `x509` - (Required) An `x509` block as defined below for X.509 certificate authentication.

---

An `x509` block supports the following:

* `secret_ref` - (Required) Reference to the Kubernetes secret containing the X.509 certificate for authenticating with the custom endpoint.

---

A `service_account_token_settings` block supports the following:

* `audiences` - (Required) List of acceptable audiences for the service account token. At least one audience must be specified.

---

An `x509_settings` block supports the following:

* `trusted_client_ca_cert` - (Optional) Name of the ConfigMap containing trusted client CA certificates for X.509 authentication.

* `authorization_attributes` - (Optional) A list of `authorization_attributes` blocks as defined below for mapping certificate attributes to authorization properties.

---

An `authorization_attributes` block supports the following:

* `name` - (Required) The name of the authorization attribute mapping.

* `subject` - (Required) The X.509 certificate subject field to extract (e.g., "CN", "O", "OU", "L", "C").

* `attributes` - (Required) A map of attribute key-value pairs to associate with clients matching this certificate subject.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ID of the IoT Operations Broker Authentication.

* `provisioning_state` - The provisioning state of the Broker Authentication.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Operations Broker Authentication.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Operations Broker Authentication.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Operations Broker Authentication.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Operations Broker Authentication.

## Import

An IoT Operations Broker Authentication can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotoperations_broker_authentication.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.IoTOperations/instances/instance1/brokers/broker1/authentications/authentication1
```
