---
subcategory: "IoT Operations"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotoperations_dataflow_endpoint"
description: |-
  Manages an Azure IoT Operations Dataflow Endpoint.
---

# azurerm_iotoperations_dataflow_endpoint

Manages an Azure IoT Operations Dataflow Endpoint.

A Dataflow Endpoint defines connectivity and authentication settings for data sources and destinations used in IoT Operations dataflows. It supports multiple endpoint types including MQTT brokers, Kafka clusters, Azure Data Explorer, Data Lake Storage, Fabric OneLake, and local storage with various authentication methods and configuration options.

## Example Usage

### MQTT Endpoint with TLS

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

resource "azurerm_iotoperations_dataflow_endpoint" "mqtt_endpoint" {
  name                = "mqtt-broker-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  endpoint_type      = "Mqtt"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  mqtt_settings {
    host                   = "broker.example.com:8883"
    protocol              = "Mqtt"
    keep_alive_seconds    = 60
    session_expiry_seconds = 3600
    max_inflight_messages = 100
    qos                   = 1
    retain                = "Keep"
    client_id_prefix      = "iot-dataflow"

    tls_settings {
      mode                                      = "Enabled"
      trusted_ca_certificate_config_map_ref = "mqtt-ca-certs"
    }

    authentication {
      method = "X509Certificate"

      x509_certificate_settings {
        secret_ref = "mqtt-client-cert"
      }
    }
  }
}
```

### Kafka Endpoint with SASL Authentication

```hcl
resource "azurerm_iotoperations_dataflow_endpoint" "kafka_endpoint" {
  name                = "kafka-cluster-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  endpoint_type      = "Kafka"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  kafka_settings {
    host              = "kafka-cluster.example.com:9092"
    consumer_group_id = "iot-consumer-group"
    compression       = "Gzip"
    partition_strategy = "Static"
    kafka_acks        = "All"
    copy_mqtt_properties = "Enabled"

    batching {
      mode         = "Enabled"
      latency_ms   = 1000
      max_bytes    = 1048576  # 1MB
      max_messages = 1000
    }

    tls_settings {
      mode                                      = "Enabled"
      trusted_ca_certificate_config_map_ref = "kafka-ca-certs"
    }

    authentication {
      method = "Sasl"

      sasl_settings {
        sasl_type  = "ScramSha256"
        secret_ref = "kafka-credentials"
      }
    }
  }
}
```

### Azure Data Explorer Endpoint with Managed Identity

```hcl
resource "azurerm_iotoperations_dataflow_endpoint" "adx_endpoint" {
  name                = "data-explorer-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  endpoint_type      = "DataExplorer"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  data_explorer_settings {
    host     = "https://mycluster.westeurope.kusto.windows.net"
    database = "IotTelemetry"

    batching {
      latency_seconds = 60
      max_messages    = 10000
    }

    authentication {
      method = "SystemAssignedManagedIdentity"

      system_assigned_managed_identity_settings {
        audience = "https://management.azure.com/"
      }
    }
  }
}
```

### Data Lake Storage Endpoint with User-Assigned Managed Identity

```hcl
resource "azurerm_user_assigned_identity" "example" {
  name                = "dataflow-identity"
  resource_group_name = azurerm_resource_group.example.name
  location           = azurerm_resource_group.example.location
}

resource "azurerm_iotoperations_dataflow_endpoint" "adls_endpoint" {
  name                = "data-lake-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  endpoint_type      = "DataLakeStorage"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  data_lake_storage_settings {
    host = "https://mystorage.dfs.core.windows.net"

    batching {
      latency_seconds = 300    # 5 minutes
      max_messages    = 50000
    }

    authentication {
      method = "UserAssignedManagedIdentity"

      user_assigned_managed_identity_settings {
        client_id = azurerm_user_assigned_identity.example.client_id
        audience  = "https://storage.azure.com/"
      }
    }
  }
}
```

### Fabric OneLake Endpoint

```hcl
resource "azurerm_iotoperations_dataflow_endpoint" "fabric_endpoint" {
  name                = "fabric-onelake-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  endpoint_type      = "FabricOneLake"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  fabric_one_lake_settings {
    host               = "https://onelake.dfs.fabric.microsoft.com"
    one_lake_path_type = "Tables"
    workspace          = "manufacturing-workspace"
    names              = ["sensor-data-lakehouse", "analytics-lakehouse"]

    batching {
      latency_seconds = 120
      max_messages    = 25000
    }

    authentication {
      method = "SystemAssignedManagedIdentity"

      system_assigned_managed_identity_settings {
        audience = "https://storage.azure.com/"
      }
    }
  }
}
```

### Local Storage Endpoint

```hcl
resource "azurerm_iotoperations_dataflow_endpoint" "local_storage_endpoint" {
  name                = "local-storage-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  endpoint_type      = "LocalStorage"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  local_storage_settings {
    path = "/mnt/data/iot-operations"
  }
}
```

### Multi-Protocol MQTT Endpoint with Service Account Token

```hcl
resource "azurerm_iotoperations_dataflow_endpoint" "websocket_mqtt_endpoint" {
  name                = "websocket-mqtt-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  endpoint_type      = "Mqtt"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  mqtt_settings {
    host                     = "wss://mqtt.example.com:443"
    protocol                = "WebSockets"
    keep_alive_seconds      = 30
    session_expiry_seconds  = 7200
    max_inflight_messages   = 50
    qos                     = 2
    retain                  = "Never"
    client_id_prefix        = "websocket-client"
    cloud_event_attributes  = "Propagate"

    authentication {
      method = "ServiceAccountToken"

      service_account_token_settings {
        audience = "iotoperations.azure.com"
      }
    }
  }
}
```

### Enterprise Kafka with Advanced Configuration

```hcl
resource "azurerm_iotoperations_dataflow_endpoint" "enterprise_kafka" {
  name                = "enterprise-kafka-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  endpoint_type      = "Kafka"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  kafka_settings {
    host                   = "kafka.enterprise.com:9093"
    consumer_group_id      = "enterprise-analytics"
    compression           = "Lz4"
    partition_strategy    = "Topic"
    kafka_acks            = "Leader"
    copy_mqtt_properties  = "Disabled"
    cloud_event_attributes = "CreateOrRemap"

    batching {
      mode         = "Enabled"
      latency_ms   = 500
      max_bytes    = 2097152  # 2MB
      max_messages = 5000
    }

    tls_settings {
      mode                                      = "Enabled"
      trusted_ca_certificate_config_map_ref = "enterprise-ca-bundle"
    }

    authentication {
      method = "UserAssignedManagedIdentity"

      user_assigned_managed_identity_settings {
        client_id = azurerm_user_assigned_identity.example.client_id
        audience  = "https://eventhubs.azure.net/"
      }
    }
  }
}
```

### Data Lake with Access Token Authentication

```hcl
resource "azurerm_iotoperations_dataflow_endpoint" "adls_token_endpoint" {
  name                = "adls-token-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  instance_name      = azurerm_iotoperations_instance.example.name
  endpoint_type      = "DataLakeStorage"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  data_lake_storage_settings {
    host = "https://analytics.dfs.core.windows.net"

    batching {
      latency_seconds = 180
      max_messages    = 75000
    }

    authentication {
      method = "AccessToken"

      access_token_settings {
        secret_ref = "adls-access-token"
      }
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Dataflow Endpoint. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the IoT Operations Dataflow Endpoint should exist. Must be between 1-90 characters. Changing this forces a new resource to be created.

* `instance_name` - (Required) The name of the IoT Operations Instance. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `endpoint_type` - (Required) The type of the endpoint. Possible values are `DataExplorer`, `DataLakeStorage`, `FabricOneLake`, `Kafka`, `LocalStorage`, and `Mqtt`. Changing this forces a new resource to be created.

* `extended_location` - (Required) An `extended_location` block as defined below. Changing this forces a new resource to be created.

* `data_explorer_settings` - (Optional) A `data_explorer_settings` block as defined below. Required when `endpoint_type` is `DataExplorer`. Cannot be used with other endpoint settings.

* `data_lake_storage_settings` - (Optional) A `data_lake_storage_settings` block as defined below. Required when `endpoint_type` is `DataLakeStorage`. Cannot be used with other endpoint settings.

* `fabric_one_lake_settings` - (Optional) A `fabric_one_lake_settings` block as defined below. Required when `endpoint_type` is `FabricOneLake`. Cannot be used with other endpoint settings.

* `kafka_settings` - (Optional) A `kafka_settings` block as defined below. Required when `endpoint_type` is `Kafka`. Cannot be used with other endpoint settings.

* `local_storage_settings` - (Optional) A `local_storage_settings` block as defined below. Required when `endpoint_type` is `LocalStorage`. Cannot be used with other endpoint settings.

* `mqtt_settings` - (Optional) An `mqtt_settings` block as defined below. Required when `endpoint_type` is `Mqtt`. Cannot be used with other endpoint settings.

---

An `extended_location` block supports the following:

* `name` - (Required) The extended location name where the IoT Operations Dataflow Endpoint should be deployed.

* `type` - (Required) The extended location type. Must be `CustomLocation`.

---

A `data_explorer_settings` block supports the following:

* `host` - (Required) The Data Explorer cluster URL. Must be between 1-253 characters.

* `database` - (Required) The Data Explorer database name. Must be between 1-253 characters.

* `authentication` - (Required) An `authentication` block as defined below for Data Explorer authentication.

* `batching` - (Optional) A `batching` block as defined below for batching configuration.

---

A `data_lake_storage_settings` block supports the following:

* `host` - (Required) The Data Lake Storage account URL. Must be between 1-253 characters.

* `authentication` - (Required) An `authentication` block as defined below for Data Lake Storage authentication.

* `batching` - (Optional) A `batching` block as defined below for batching configuration.

---

A `fabric_one_lake_settings` block supports the following:

* `host` - (Required) The Fabric OneLake host URL. Must be between 1-253 characters.

* `workspace` - (Required) The Fabric workspace name. Must be between 1-253 characters.

* `names` - (Required) List of lakehouse names to connect to. Each name must be between 1-253 characters.

* `one_lake_path_type` - (Required) The OneLake path type. Possible values are `Files` and `Tables`.

* `authentication` - (Required) An `authentication` block as defined below for Fabric OneLake authentication.

* `batching` - (Optional) A `batching` block as defined below for batching configuration.

---

A `kafka_settings` block supports the following:

* `host` - (Required) The Kafka broker connection string. Must be between 1-253 characters.

* `authentication` - (Required) An `authentication` block as defined below for Kafka authentication.

* `consumer_group_id` - (Optional) The Kafka consumer group ID. Must be between 1-253 characters.

* `compression` - (Optional) The compression type for Kafka messages. Possible values are `None`, `Gzip`, `Snappy`, and `Lz4`.

* `partition_strategy` - (Optional) The partitioning strategy for Kafka messages.

* `kafka_acks` - (Optional) The acknowledgment level for Kafka producers.

* `copy_mqtt_properties` - (Optional) Whether to copy MQTT properties to Kafka headers.

* `cloud_event_attributes` - (Optional) How to handle CloudEvent attributes.

* `batching` - (Optional) A `kafka_batching` block as defined below for Kafka-specific batching configuration.

* `tls_settings` - (Optional) A `tls_settings` block as defined below for TLS configuration.

---

A `local_storage_settings` block supports the following:

* `path` - (Required) The local storage path. Must be between 1-1000 characters.

---

An `mqtt_settings` block supports the following:

* `host` - (Required) The MQTT broker connection string. Must be between 1-253 characters.

* `authentication` - (Required) An `authentication` block as defined below for MQTT authentication.

* `protocol` - (Optional) The MQTT protocol type. Possible values are `Mqtt` and `WebSockets`.

* `keep_alive_seconds` - (Optional) The MQTT keep-alive interval in seconds. Must be between 1-65535.

* `session_expiry_seconds` - (Optional) The MQTT session expiry interval in seconds. Must be between 1-4294967295.

* `max_inflight_messages` - (Optional) The maximum number of in-flight MQTT messages. Must be between 1-65535.

* `qos` - (Optional) The Quality of Service level for MQTT messages. Must be between 0-2.

* `retain` - (Optional) The MQTT retain policy. Possible values are `Keep` and `Never`.

* `client_id_prefix` - (Optional) The prefix for MQTT client IDs. Must be between 1-253 characters.

* `cloud_event_attributes` - (Optional) How to handle CloudEvent attributes in MQTT messages.

* `tls_settings` - (Optional) A `tls_settings` block as defined below for TLS configuration.

---

A `batching` block supports the following:

* `latency_seconds` - (Optional) The maximum latency in seconds before sending a batch. Must be between 1-3600.

* `max_messages` - (Optional) The maximum number of messages in a batch. Must be between 1-1000000.

---

A `kafka_batching` block supports the following:

* `mode` - (Optional) Whether batching is enabled. Possible values are `Enabled` and `Disabled`.

* `latency_ms` - (Optional) The maximum latency in milliseconds before sending a batch. Must be between 0-3600000.

* `max_bytes` - (Optional) The maximum size of a batch in bytes. Must be between 1-1073741824.

* `max_messages` - (Optional) The maximum number of messages in a batch. Must be between 1-1000000.

---

A `tls_settings` block supports the following:

* `mode` - (Required) Whether TLS is enabled. Possible values are `Enabled` and `Disabled`.

* `trusted_ca_certificate_config_map_ref` - (Optional) Reference to the ConfigMap containing trusted CA certificates. Must be between 1-253 characters.

---

An `authentication` block supports the following:

* `method` - (Required) The authentication method. Possible values are `SystemAssignedManagedIdentity`, `UserAssignedManagedIdentity`, `ServiceAccountToken`, `X509Certificate`, `AccessToken`, and `Sasl`.

* `system_assigned_managed_identity_settings` - (Optional) A `system_assigned_managed_identity_settings` block as defined below. Required when `method` is `SystemAssignedManagedIdentity`.

* `user_assigned_managed_identity_settings` - (Optional) A `user_assigned_managed_identity_settings` block as defined below. Required when `method` is `UserAssignedManagedIdentity`.

* `service_account_token_settings` - (Optional) A `service_account_token_settings` block as defined below. Required when `method` is `ServiceAccountToken`.

* `x509_certificate_settings` - (Optional) An `x509_certificate_settings` block as defined below. Required when `method` is `X509Certificate`.

* `access_token_settings` - (Optional) An `access_token_settings` block as defined below. Required when `method` is `AccessToken`.

* `sasl_settings` - (Optional) A `sasl_settings` block as defined below. Required when `method` is `Sasl`.

---

A `system_assigned_managed_identity_settings` block supports the following:

* `audience` - (Required) The audience for the managed identity token. Must be between 1-253 characters.

---

A `user_assigned_managed_identity_settings` block supports the following:

* `client_id` - (Required) The client ID of the user-assigned managed identity.

* `audience` - (Required) The audience for the managed identity token. Must be between 1-253 characters.

---

A `service_account_token_settings` block supports the following:

* `audience` - (Required) The audience for the service account token. Must be between 1-253 characters.

---

An `x509_certificate_settings` block supports the following:

* `secret_ref` - (Required) Reference to the Kubernetes secret containing the X.509 certificate. Must be between 1-253 characters.

---

An `access_token_settings` block supports the following:

* `secret_ref` - (Required) Reference to the Kubernetes secret containing the access token. Must be between 1-253 characters.

---

A `sasl_settings` block supports the following:

* `sasl_type` - (Required) The SASL mechanism type. Possible values are `Plain`, `ScramSha256`, and `ScramSha512`.

* `secret_ref` - (Required) Reference to the Kubernetes secret containing SASL credentials. Must be between 1-253 characters.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ID of the IoT Operations Dataflow Endpoint.

* `provisioning_state` - The provisioning state of the Dataflow Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Operations Dataflow Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Operations Dataflow Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Operations Dataflow Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Operations Dataflow Endpoint.

## Import

An IoT Operations Dataflow Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotoperations_dataflow_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.IoTOperations/instances/instance1/dataflowEndpoints/endpoint1
```
