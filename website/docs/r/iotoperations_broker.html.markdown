---
subcategory: "IoT Operations"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotoperations_broker"
description: |-
  Manages an Azure IoT Operations Broker.
---

# azurerm_iotoperations_broker

Manages an Azure IoT Operations Broker.

An IoT Operations Broker is the central message routing component that handles MQTT communication between IoT devices and applications. It provides high-performance message routing with configurable scaling, diagnostics, security, and storage options.

## Example Usage

### Basic Broker

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
```

### Broker with Performance Tuning

```hcl
resource "azurerm_iotoperations_broker" "high_performance" {
  name                = "high-perf-broker"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = azurerm_iotoperations_instance.example.extended_location_type
  }

  properties {
    memory_profile = "High"

    cardinality {
      backend_chain {
        partitions        = 8
        redundancy_factor = 3
        workers          = 4
      }

      frontend {
        replicas = 4
        workers  = 2
      }
    }

    generate_resource_limits {
      cpu = "Enabled"
    }
  }
}
```

### Broker with Advanced Configuration

```hcl
resource "azurerm_iotoperations_broker" "advanced" {
  name                = "advanced-broker"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = azurerm_iotoperations_instance.example.extended_location_type
  }

  properties {
    memory_profile = "Medium"

    advanced {
      encrypt_internal_traffic = "Enabled"

      clients {
        max_session_expiry_seconds = 3600
        max_message_expiry_seconds = 1800
        max_packet_size_bytes     = 1048576
        max_receive_maximum       = 100
        max_keep_alive_seconds    = 300

        subscriber_queue_limit {
          length   = 1000
          strategy = "DropOldest"
        }
      }

      internal_certs {
        duration     = "8760h"  # 1 year
        renew_before = "720h"   # 30 days

        private_key {
          algorithm       = "Rsa2048"
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
        timeout_seconds  = 10
      }

      traces {
        mode                  = "Enabled"
        cache_size_megabytes = 16
        span_channel_capacity = 1000

        self_tracing {
          mode             = "Enabled"
          interval_seconds = 30
        }
      }
    }
  }
}
```

### Broker with Persistent Storage

```hcl
resource "azurerm_iotoperations_broker" "with_storage" {
  name                = "storage-broker"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = azurerm_iotoperations_instance.example.extended_location_type
  }

  properties {
    memory_profile = "Medium"

    disk_backed_message_buffer {
      max_size = "1Gi"

      persistent_volume_claim_spec {
        access_modes         = ["ReadWriteOnce"]
        storage_class_name   = "fast-ssd"
        volume_mode         = "Filesystem"

        resources {
          requests = {
            storage = "10Gi"
          }
          limits = {
            storage = "50Gi"
          }
        }

        selector {
          match_labels = {
            tier = "storage"
            type = "fast"
          }

          match_expressions {
            key      = "environment"
            operator = "In"
            values   = ["production", "staging"]
          }
        }
      }
    }
  }
}
```

### Broker with Ephemeral Storage

```hcl
resource "azurerm_iotoperations_broker" "ephemeral_storage" {
  name                = "ephemeral-broker"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = azurerm_iotoperations_instance.example.extended_location_type
  }

  properties {
    disk_backed_message_buffer {
      max_size = "512Mi"

      ephemeral_volume_claim_spec {
        access_modes = ["ReadWriteOnce"]
        volume_mode = "Filesystem"

        resources {
          requests = {
            storage = "1Gi"
          }
        }

        data_source {
          api_group = "snapshot.storage.k8s.io"
          kind      = "VolumeSnapshot"
          name      = "broker-snapshot"
        }
      }
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Broker. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the IoT Operations Broker should exist. Must be between 1-90 characters. Changing this forces a new resource to be created.

* `instance_name` - (Required) The name of the IoT Operations Instance. Changing this forces a new resource to be created.

* `extended_location` - (Required) An `extended_location` block as defined below. Changing this forces a new resource to be created.

* `properties` - (Optional) A `properties` block as defined below for configuring broker behavior.

---

An `extended_location` block supports the following:

* `name` - (Required) The extended location name where the IoT Operations Broker should be deployed.

* `type` - (Required) The extended location type.

---

A `properties` block supports the following:

* `memory_profile` - (Optional) The memory profile for the broker. Possible values are `Tiny`, `Low`, `Medium`, and `High`.

* `advanced` - (Optional) An `advanced` block as defined below for advanced broker configuration.

* `cardinality` - (Optional) A `cardinality` block as defined below for scaling configuration.

* `diagnostics` - (Optional) A `diagnostics` block as defined below for monitoring and logging.

* `disk_backed_message_buffer` - (Optional) A `disk_backed_message_buffer` block as defined below for persistent message storage.

* `generate_resource_limits` - (Optional) A `generate_resource_limits` block as defined below for automatic resource limit generation.

---

An `advanced` block supports the following:

* `encrypt_internal_traffic` - (Optional) Whether to encrypt internal traffic between broker components.

* `clients` - (Optional) A `clients` block as defined below for client connection settings.

* `internal_certs` - (Optional) An `internal_certs` block as defined below for internal certificate configuration.

---

A `clients` block supports the following:

* `max_session_expiry_seconds` - (Optional) Maximum session expiry time in seconds.

* `max_message_expiry_seconds` - (Optional) Maximum message expiry time in seconds.

* `max_packet_size_bytes` - (Optional) Maximum MQTT packet size in bytes.

* `max_receive_maximum` - (Optional) Maximum number of QoS 1 and QoS 2 publications that the client is willing to process concurrently.

* `max_keep_alive_seconds` - (Optional) Maximum keep alive time in seconds.

* `subscriber_queue_limit` - (Optional) A `subscriber_queue_limit` block as defined below.

---

A `subscriber_queue_limit` block supports the following:

* `length` - (Optional) Maximum number of messages in the subscriber queue.

* `strategy` - (Optional) Strategy for handling queue overflow. Possible values include message drop strategies.

---

An `internal_certs` block supports the following:

* `duration` - (Optional) The duration of the certificate validity period (e.g., "8760h" for 1 year).

* `renew_before` - (Optional) The time before expiry when the certificate should be renewed (e.g., "720h" for 30 days).

* `private_key` - (Optional) A `private_key` block as defined below.

---

A `private_key` block supports the following:

* `algorithm` - (Optional) The algorithm for the private key.

* `rotation_policy` - (Optional) The rotation policy for the private key.

---

A `cardinality` block supports the following:

* `backend_chain` - (Optional) A `backend_chain` block as defined below for backend scaling configuration.

* `frontend` - (Optional) A `frontend` block as defined below for frontend scaling configuration.

---

A `backend_chain` block supports the following:

* `partitions` - (Optional) Number of partitions for the backend chain. Must be between 1-16.

* `redundancy_factor` - (Optional) Redundancy factor for the backend chain. Must be between 1-5.

* `workers` - (Optional) Number of worker threads for the backend chain. Must be between 1-16.

---

A `frontend` block supports the following:

* `replicas` - (Optional) Number of frontend replicas. Must be between 1-16.

* `workers` - (Optional) Number of worker threads per frontend replica. Must be between 1-16.

---

A `diagnostics` block supports the following:

* `logs` - (Optional) A `logs` block as defined below for logging configuration.

* `metrics` - (Optional) A `metrics` block as defined below for metrics collection.

* `self_check` - (Optional) A `self_check` block as defined below for health checking.

* `traces` - (Optional) A `traces` block as defined below for distributed tracing.

---

A `logs` block supports the following:

* `level` - (Optional) The logging level (e.g., "info", "debug", "warn", "error").

---

A `metrics` block supports the following:

* `prometheus_port` - (Optional) Port number for Prometheus metrics endpoint. Must be between 0-65535.

---

A `self_check` block supports the following:

* `mode` - (Optional) Whether self-checking is enabled.

* `interval_seconds` - (Optional) Interval between self-checks in seconds.

* `timeout_seconds` - (Optional) Timeout for self-check operations in seconds.

---

A `traces` block supports the following:

* `mode` - (Optional) Whether distributed tracing is enabled.

* `cache_size_megabytes` - (Optional) Size of the trace cache in megabytes.

* `span_channel_capacity` - (Optional) Capacity of the span channel.

* `self_tracing` - (Optional) A `self_tracing` block as defined below.

---

A `self_tracing` block supports the following:

* `mode` - (Optional) Whether self-tracing is enabled.

* `interval_seconds` - (Optional) Interval between self-trace operations in seconds.

---

A `disk_backed_message_buffer` block supports the following:

* `max_size` - (Optional) Maximum size of the message buffer (e.g., "1Gi", "512Mi").

* `ephemeral_volume_claim_spec` - (Optional) A `volume_claim_spec` block as defined below for ephemeral storage.

* `persistent_volume_claim_spec` - (Optional) A `volume_claim_spec` block as defined below for persistent storage.

---

A `volume_claim_spec` block supports the following:

* `volume_name` - (Optional) Name of the volume.

* `volume_mode` - (Optional) Volume mode (e.g., "Filesystem", "Block").

* `storage_class_name` - (Optional) Storage class name for the volume.

* `access_modes` - (Optional) List of access modes for the volume (e.g., ["ReadWriteOnce", "ReadOnlyMany"]).

* `data_source` - (Optional) A `data_source` block as defined below.

* `data_source_ref` - (Optional) A `data_source_ref` block as defined below.

* `resources` - (Optional) A `resources` block as defined below for resource requirements.

* `selector` - (Optional) A `selector` block as defined below for volume selection.

---

A `data_source` block supports the following:

* `api_group` - (Optional) API group of the data source.

* `kind` - (Optional) Kind of the data source.

* `name` - (Optional) Name of the data source.

---

A `data_source_ref` block supports the following:

* `api_group` - (Optional) API group of the data source reference.

* `kind` - (Optional) Kind of the data source reference.

* `name` - (Optional) Name of the data source reference.

* `namespace` - (Optional) Namespace of the data source reference.

---

A `resources` block supports the following:

* `limits` - (Optional) Map of resource limits (e.g., `{"storage" = "10Gi", "cpu" = "500m"}`).

* `requests` - (Optional) Map of resource requests (e.g., `{"storage" = "5Gi", "memory" = "1Gi"}`).

---

A `selector` block supports the following:

* `match_labels` - (Optional) Map of labels to match.

* `match_expressions` - (Optional) List of `match_expression` blocks as defined below.

---

A `match_expression` block supports the following:

* `key` - (Optional) The label key.

* `operator` - (Optional) The operator for matching (e.g., "In", "NotIn", "Exists", "DoesNotExist").

* `values` - (Optional) List of values to match against.

---

A `generate_resource_limits` block supports the following:

* `cpu` - (Optional) Whether to automatically generate CPU resource limits.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ID of the IoT Operations Broker.

* `provisioning_state` - The provisioning state of the Broker.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Operations Broker.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Operations Broker.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Operations Broker.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Operations Broker.

## Import

An IoT Operations Broker can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotoperations_broker.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.IoTOperations/instances/instance1/brokers/broker1
```
