---
subcategory: "IoT Operations"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotoperations_broker_authorization"
description: |-
  Manages an Azure IoT Operations Broker Authorization.
---

# azurerm_iotoperations_broker_authorization

Manages an Azure IoT Operations Broker Authorization.

A Broker Authorization defines access control policies for the IoT Operations broker, determining what actions authenticated clients can perform. It supports fine-grained permissions for MQTT operations (connect, publish, subscribe) and state store access, with flexible principal matching based on client IDs, usernames, or custom attributes.

## Example Usage

### Basic MQTT Authorization

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

resource "azurerm_iotoperations_broker_authorization" "basic_mqtt" {
  name                = "basic-mqtt-authz"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  authorization_policies {
    cache = "Enabled"

    rules {
      # Allow device connection
      broker_resources {
        method = "Connect"
      }

      # Allow publishing to device topics
      broker_resources {
        method = "Publish"
        topics = [
          "devices/+/telemetry",
          "devices/+/status"
        ]
      }

      # Allow subscribing to command topics
      broker_resources {
        method = "Subscribe"
        topics = [
          "devices/+/commands",
          "devices/+/config"
        ]
      }

      principals {
        clients = [
          "device-*",
          "sensor-*"
        ]
      }
    }
  }
}
```

### Role-Based Authorization

```hcl
resource "azurerm_iotoperations_broker_authorization" "role_based" {
  name                = "role-based-authz"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  authorization_policies {
    cache = "Enabled"

    # Admin users - full access
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Publish"
        topics = ["*"]
      }

      broker_resources {
        method = "Subscribe"
        topics = ["*"]
      }

      principals {
        usernames = [
          "admin",
          "system-operator"
        ]
      }
    }

    # Device operators - device management
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Publish"
        topics = [
          "devices/+/commands",
          "devices/+/config",
          "management/+"
        ]
      }

      broker_resources {
        method = "Subscribe"
        topics = [
          "devices/+/telemetry",
          "devices/+/status",
          "devices/+/alerts"
        ]
      }

      principals {
        usernames = [
          "device-operator",
          "maintenance-tech"
        ]
      }
    }

    # Read-only monitoring users
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Subscribe"
        topics = [
          "devices/+/telemetry",
          "devices/+/status",
          "system/metrics"
        ]
      }

      principals {
        usernames = [
          "monitor-user",
          "dashboard-service"
        ]
      }
    }
  }
}
```

### Attribute-Based Authorization

```hcl
resource "azurerm_iotoperations_broker_authorization" "attribute_based" {
  name                = "attribute-based-authz"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  authorization_policies {
    cache = "Enabled"

    # Manufacturing floor devices
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Publish"
        topics = [
          "factory/floor1/+/data",
          "factory/floor1/+/status"
        ]
      }

      broker_resources {
        method = "Subscribe"
        topics = [
          "factory/floor1/+/commands",
          "factory/broadcast"
        ]
      }

      principals {
        attributes = [
          {
            "location"    = "floor1"
            "device-type" = "sensor"
            "access"      = "production"
          }
        ]
      }
    }

    # Quality control systems
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Subscribe"
        topics = [
          "factory/+/+/data",
          "quality/+/reports"
        ]
      }

      broker_resources {
        method = "Publish"
        topics = [
          "quality/+/alerts",
          "quality/+/reports"
        ]
      }

      principals {
        attributes = [
          {
            "department" = "quality-control"
            "clearance"  = "level-2"
          }
        ]
      }
    }

    # Maintenance team access
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Publish"
        topics = [
          "maintenance/+/schedule",
          "maintenance/+/reports"
        ]
      }

      broker_resources {
        method = "Subscribe"
        topics = [
          "factory/+/+/diagnostics",
          "maintenance/+/alerts"
        ]
      }

      principals {
        attributes = [
          {
            "team"   = "maintenance"
            "shift"  = "day"
            "level"  = "technician"
          },
          {
            "team"   = "maintenance"
            "shift"  = "night"
            "level"  = "supervisor"
          }
        ]
      }
    }
  }
}
```

### State Store Authorization

```hcl
resource "azurerm_iotoperations_broker_authorization" "state_store" {
  name                = "state-store-authz"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  authorization_policies {
    cache = "Disabled"  # Disable caching for dynamic permissions

    # Application state management
    rules {
      broker_resources {
        method = "Connect"
      }

      # State store access for application data
      state_store_resources {
        key_type = "String"
        keys = [
          "app.config.*",
          "app.state.*",
          "app.cache.*"
        ]
        method = "ReadWrite"
      }

      # Binary data access
      state_store_resources {
        key_type = "Binary"
        keys = [
          "blobs/images/*",
          "blobs/documents/*"
        ]
        method = "ReadWrite"
      }

      # Pattern-based access for dynamic keys
      state_store_resources {
        key_type = "Pattern"
        keys = [
          "user\\..*\\.settings",
          "session\\..*\\.data"
        ]
        method = "ReadWrite"
      }

      principals {
        clients = [
          "application-*",
          "service-*"
        ]
      }
    }

    # Read-only analytics access
    rules {
      broker_resources {
        method = "Connect"
      }

      # Read-only access to analytics data
      state_store_resources {
        key_type = "String"
        keys = [
          "analytics.*",
          "reports.*",
          "metrics.*"
        ]
        method = "Read"
      }

      principals {
        usernames = [
          "analytics-service",
          "reporting-engine"
        ]
        attributes = [
          {
            "service-type" = "analytics"
            "access-level" = "read-only"
          }
        ]
      }
    }

    # Device configuration management
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Subscribe"
        topics = [
          "devices/+/config-request"
        ]
      }

      broker_resources {
        method = "Publish"
        topics = [
          "devices/+/config-response"
        ]
      }

      # Device-specific configuration access
      state_store_resources {
        key_type = "Pattern"
        keys = [
          "device\\..+\\.config",
          "device\\..+\\.firmware"
        ]
        method = "Write"
      }

      principals {
        clients = [
          "config-manager"
        ]
        attributes = [
          {
            "service" = "device-management"
            "role"    = "configurator"
          }
        ]
      }
    }
  }
}
```

### Multi-Tenant Authorization

```hcl
resource "azurerm_iotoperations_broker_authorization" "multi_tenant" {
  name                = "multi-tenant-authz"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  broker_name        = azurerm_iotoperations_broker.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  authorization_policies {
    cache = "Enabled"

    # Tenant A - restricted to their namespace
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Publish"
        topics = [
          "tenants/tenant-a/+",
          "tenants/tenant-a/+/+"
        ]
      }

      broker_resources {
        method = "Subscribe"
        topics = [
          "tenants/tenant-a/+",
          "tenants/tenant-a/+/+",
          "shared/announcements"
        ]
      }

      state_store_resources {
        key_type = "Pattern"
        keys = [
          "tenant-a\\..*"
        ]
        method = "ReadWrite"
      }

      principals {
        attributes = [
          {
            "tenant-id" = "tenant-a"
            "org"       = "company-alpha"
          }
        ]
        clients = [
          "tenant-a-*"
        ]
      }
    }

    # Tenant B - restricted to their namespace
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Publish"
        topics = [
          "tenants/tenant-b/+",
          "tenants/tenant-b/+/+"
        ]
      }

      broker_resources {
        method = "Subscribe"
        topics = [
          "tenants/tenant-b/+",
          "tenants/tenant-b/+/+",
          "shared/announcements"
        ]
      }

      state_store_resources {
        key_type = "Pattern"
        keys = [
          "tenant-b\\..*"
        ]
        method = "ReadWrite"
      }

      principals {
        attributes = [
          {
            "tenant-id" = "tenant-b"
            "org"       = "company-beta"
          }
        ]
        clients = [
          "tenant-b-*"
        ]
      }
    }

    # Platform administrators - cross-tenant access
    rules {
      broker_resources {
        method = "Connect"
      }

      broker_resources {
        method = "Publish"
        topics = [
          "*"
        ]
      }

      broker_resources {
        method = "Subscribe"
        topics = [
          "*"
        ]
      }

      state_store_resources {
        key_type = "String"
        keys = [
          "*"
        ]
        method = "ReadWrite"
      }

      principals {
        usernames = [
          "platform-admin",
          "system-monitor"
        ]
        attributes = [
          {
            "role"  = "platform-admin"
            "level" = "global"
          }
        ]
      }
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Broker Authorization. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the IoT Operations Broker Authorization should exist. Must be between 1-90 characters. Changing this forces a new resource to be created.

* `instance_name` - (Required) The name of the IoT Operations Instance. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `broker_name` - (Required) The name of the IoT Operations Broker. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `extended_location` - (Required) An `extended_location` block as defined below. Changing this forces a new resource to be created.

* `authorization_policies` - (Required) An `authorization_policies` block as defined below for configuring access control policies.

---

An `extended_location` block supports the following:

* `name` - (Required) The extended location name where the IoT Operations Broker Authorization should be deployed.

* `type` - (Optional) The extended location type. Defaults to `CustomLocation`.

---

An `authorization_policies` block supports the following:

* `cache` - (Optional) Whether to cache authorization decisions for performance. Possible values are `Enabled` and `Disabled`. Defaults to `Enabled`.

* `rules` - (Required) A list of `rules` blocks as defined below. At least one rule must be configured.

---

A `rules` block supports the following:

* `broker_resources` - (Required) A list of `broker_resources` blocks as defined below for MQTT operation permissions. At least one broker resource rule must be configured.

* `principals` - (Required) A `principals` block as defined below for identifying which clients this rule applies to.

* `state_store_resources` - (Optional) A list of `state_store_resources` blocks as defined below for state store access permissions.

---

A `broker_resources` block supports the following:

* `method` - (Required) The MQTT operation method. Possible values are `Connect`, `Publish`, and `Subscribe`.

* `clients` - (Optional) List of client ID patterns that this rule applies to. Supports wildcards with `+` (single level) and `*` (multi-level).

* `topics` - (Optional) List of MQTT topic patterns that this rule applies to. Required for `Publish` and `Subscribe` methods. Supports wildcards with `+` (single level) and `*` (multi-level).

---

A `principals` block supports the following:

* `clients` - (Optional) List of client ID patterns to match against connecting clients.

* `usernames` - (Optional) List of usernames to match against authenticated clients.

* `attributes` - (Optional) List of attribute maps to match against client authentication attributes. Each map contains key-value pairs that must all match the client's attributes.

---

A `state_store_resources` block supports the following:

* `key_type` - (Required) The type of key matching to perform. Possible values are `Binary`, `Pattern`, and `String`.

* `keys` - (Required) List of keys or key patterns to match. At least one key must be specified. For `Pattern` type, supports regular expressions.

* `method` - (Required) The access method for the state store resource. Possible values are `Read`, `ReadWrite`, and `Write`.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ID of the IoT Operations Broker Authorization.

* `provisioning_state` - The provisioning state of the Broker Authorization.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Operations Broker Authorization.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Operations Broker Authorization.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Operations Broker Authorization.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Operations Broker Authorization.

## Import

An IoT Operations Broker Authorization can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotoperations_broker_authorization.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.IoTOperations/instances/instance1/brokers/broker1/authorizations/authorization1
```
