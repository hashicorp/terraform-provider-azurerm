---
subcategory: "IoT Operations"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotoperations_dataflow_profile"
description: |-
  Manages an Azure IoT Operations Dataflow Profile.
---

# azurerm_iotoperations_dataflow_profile

Manages an Azure IoT Operations Dataflow Profile.

A Dataflow Profile defines the runtime configuration and scaling parameters for dataflow instances in IoT Operations. It controls how many dataflow instances run, their logging levels, and monitoring capabilities, providing a template for dataflow execution environments.

## Example Usage

### Basic Dataflow Profile

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

resource "azurerm_iotoperations_dataflow_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }
}
```

### High-Scale Dataflow Profile

```hcl
resource "azurerm_iotoperations_dataflow_profile" "high_scale" {
  name                = "high-scale-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 10

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "info"
    }

    metrics {
      prometheus_port = 9090
    }
  }
}
```

### Development Dataflow Profile with Debug Logging

```hcl
resource "azurerm_iotoperations_dataflow_profile" "development" {
  name                = "dev-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 2

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "debug"
    }

    metrics {
      prometheus_port = 9091
    }
  }
}
```

### Production Dataflow Profile with Minimal Logging

```hcl
resource "azurerm_iotoperations_dataflow_profile" "production" {
  name                = "prod-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 50

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "error"
    }

    metrics {
      prometheus_port = 9092
    }
  }
}
```

### Edge Computing Dataflow Profile

```hcl
resource "azurerm_iotoperations_dataflow_profile" "edge_computing" {
  name                = "edge-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 3

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "warn"
    }

    metrics {
      prometheus_port = 9093
    }
  }
}
```

### Monitoring-Focused Profile

```hcl
resource "azurerm_iotoperations_dataflow_profile" "monitoring" {
  name                = "monitoring-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 5

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "trace"
    }

    metrics {
      prometheus_port = 9094
    }
  }
}
```

### Multi-Environment Setup

```hcl
# Development environment
resource "azurerm_iotoperations_dataflow_profile" "dev_environment" {
  name                = "dev-env-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 1

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "debug"
    }

    metrics {
      prometheus_port = 9100
    }
  }
}

# Staging environment
resource "azurerm_iotoperations_dataflow_profile" "staging_environment" {
  name                = "staging-env-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 5

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "info"
    }

    metrics {
      prometheus_port = 9101
    }
  }
}

# Production environment
resource "azurerm_iotoperations_dataflow_profile" "production_environment" {
  name                = "prod-env-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 25

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "warn"
    }

    metrics {
      prometheus_port = 9102
    }
  }
}
```

### High-Throughput Profile for Analytics Workloads

```hcl
resource "azurerm_iotoperations_dataflow_profile" "analytics_workload" {
  name                = "analytics-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 100  # Maximum scaling for high throughput

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "error"  # Minimal logging for performance
    }

    metrics {
      prometheus_port = 9200  # Dedicated metrics port
    }
  }
}
```

### Resource-Constrained Profile for Edge Devices

```hcl
resource "azurerm_iotoperations_dataflow_profile" "edge_constrained" {
  name                = "edge-constrained-profile"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  instance_count     = 1  # Minimal instances for resource constraints

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  diagnostics {
    logs {
      level = "error"  # Only critical errors to save resources
    }

    # No metrics configuration to save resources
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Dataflow Profile. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the IoT Operations Dataflow Profile should exist. Must be between 1-90 characters. Changing this forces a new resource to be created.

* `instance_name` - (Required) The name of the IoT Operations Instance. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `extended_location` - (Required) An `extended_location` block as defined below. Changing this forces a new resource to be created.

* `instance_count` - (Optional) The number of dataflow instances to run. Must be between 1-1000. This determines the parallel processing capacity and scaling of dataflow operations.

* `diagnostics` - (Optional) A `diagnostics` block as defined below for configuring logging and monitoring.

---

An `extended_location` block supports the following:

* `name` - (Required) The extended location name where the IoT Operations Dataflow Profile should be deployed.

* `type` - (Required) The extended location type. Must be `CustomLocation`.

---

A `diagnostics` block supports the following:

* `logs` - (Optional) A `logs` block as defined below for logging configuration.

* `metrics` - (Optional) A `metrics` block as defined below for metrics collection configuration.

---

A `logs` block supports the following:

* `level` - (Optional) The logging level for dataflow operations. Possible values are `trace`, `debug`, `info`, `warn`, and `error`. 
  - `trace`: Most detailed logging, includes all operations and data flow tracing
  - `debug`: Detailed debugging information for troubleshooting
  - `info`: General operational information and status updates
  - `warn`: Warning messages about potential issues
  - `error`: Only error conditions and failures

---

A `metrics` block supports the following:

* `prometheus_port` - (Optional) The port number for Prometheus metrics endpoint. Must be between 1-65535. This enables monitoring and observability of dataflow performance and health.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ID of the IoT Operations Dataflow Profile.

* `provisioning_state` - The provisioning state of the Dataflow Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Operations Dataflow Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Operations Dataflow Profile.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Operations Dataflow Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Operations Dataflow Profile.

## Import

An IoT Operations Dataflow Profile can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotoperations_dataflow_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.IoTOperations/instances/instance1/dataflowProfiles/profile1
```
