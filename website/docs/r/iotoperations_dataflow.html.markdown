---
subcategory: "IoT Operations"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotoperations_dataflow"
description: |-
  Manages an Azure IoT Operations Dataflow.
---

# azurerm_iotoperations_dataflow

Manages an Azure IoT Operations Dataflow.

A Dataflow defines data processing pipelines that move and transform data between sources and destinations in IoT Operations. It supports complex data transformations including filtering, mapping, and dataset creation with multiple serialization formats and built-in transformation capabilities.

## Example Usage

### Basic Data Pipeline

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
    type = azurerm_iotoperations_instance.example.extended_location_type
  }
}

resource "azurerm_iotoperations_dataflow" "basic_pipeline" {
  name                    = "basic-pipeline"
  resource_group_name     = azurerm_resource_group.example.name
  instance_name          = azurerm_iotoperations_instance.example.name
  dataflow_profile_name  = azurerm_iotoperations_dataflow_profile.example.name
  mode                   = "Enabled"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  # Source operation - read from MQTT endpoint
  operations {
    name           = "mqtt-source"
    operation_type = "Source"

    source_settings {
      data_sources = [
        "azure-iot-operations/data/thermostat"
      ]
      endpoint_ref         = "mqtt-endpoint"
      serialization_format = "Json"
    }
  }

  # Destination operation - write to Event Hub
  operations {
    name           = "eventhub-destination"
    operation_type = "Destination"

    destination_settings {
      data_destination = "processed-telemetry"
      endpoint_ref     = "eventhub-endpoint"
    }
  }
}
```

### Complex Transformation Pipeline

```hcl
resource "azurerm_iotoperations_dataflow" "transformation_pipeline" {
  name                    = "transformation-pipeline"
  resource_group_name     = azurerm_resource_group.example.name
  instance_name          = azurerm_iotoperations_instance.example.name
  dataflow_profile_name  = azurerm_iotoperations_dataflow_profile.example.name
  mode                   = "Enabled"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  # Source operation
  operations {
    name           = "sensor-data-source"
    operation_type = "Source"

    source_settings {
      data_sources = [
        "sensors/temperature/+",
        "sensors/humidity/+",
        "sensors/pressure/+"
      ]
      endpoint_ref         = "mqtt-broker"
      asset_ref           = "factory-sensors"
      schema_ref          = "sensor-schema"
      serialization_format = "Json"
    }
  }

  # Transformation operation
  operations {
    name           = "data-processing"
    operation_type = "BuiltInTransformation"

    built_in_transformation_settings {
      schema_ref           = "processed-schema"
      serialization_format = "Json"

      # Filter out invalid readings
      filter {
        description = "Remove readings with invalid temperature values"
        expression  = "temperature >= -50 && temperature <= 150"
        inputs      = ["temperature"]
        type        = "Filter"
      }

      # Filter out test devices
      filter {
        description = "Exclude test devices from production data"
        expression  = "!deviceId.startsWith('test-')"
        inputs      = ["deviceId"]
        type        = "Filter"
      }

      # Convert temperature units
      map {
        description = "Convert Celsius to Fahrenheit"
        expression  = "(temperature * 9/5) + 32"
        inputs      = ["temperature"]
        output      = "temperatureF"
        type        = "Compute"
      }

      # Add timestamp
      map {
        description = "Add processing timestamp"
        expression  = "now()"
        inputs      = []
        output      = "processedAt"
        type        = "BuiltInFunction"
      }

      # Rename fields
      map {
        description = "Standardize field names"
        inputs      = ["deviceId"]
        output      = "device_identifier"
        type        = "Rename"
      }

      # Pass through existing fields
      map {
        description = "Keep humidity as-is"
        inputs      = ["humidity"]
        output      = "humidity"
        type        = "PassThrough"
      }

      # Add new calculated properties
      map {
        description = "Add device metadata"
        expression  = "{'location': location, 'type': 'environmental_sensor'}"
        inputs      = ["location"]
        output      = "metadata"
        type        = "NewProperties"
      }

      # Create datasets for different consumers
      datasets {
        key         = "alerts"
        description = "High-priority sensor alerts"
        expression  = "temperature > 100 || humidity > 90"
        inputs      = ["temperature", "humidity", "device_identifier"]
        schema_ref  = "alert-schema"
      }

      datasets {
        key         = "metrics"
        description = "Aggregated sensor metrics"
        expression  = "{'avg_temp': avg(temperature), 'max_humidity': max(humidity)}"
        inputs      = ["temperature", "humidity"]
        schema_ref  = "metrics-schema"
      }
    }
  }

  # Destination operation
  operations {
    name           = "cloud-storage"
    operation_type = "Destination"

    destination_settings {
      data_destination = "processed-sensor-data"
      endpoint_ref     = "adls-endpoint"
    }
  }
}
```

### Multi-Source Analytics Pipeline

```hcl
resource "azurerm_iotoperations_dataflow" "analytics_pipeline" {
  name                    = "analytics-pipeline"
  resource_group_name     = azurerm_resource_group.example.name
  instance_name          = azurerm_iotoperations_instance.example.name
  dataflow_profile_name  = azurerm_iotoperations_dataflow_profile.example.name
  mode                   = "Enabled"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  # Manufacturing line data source
  operations {
    name           = "manufacturing-source"
    operation_type = "Source"

    source_settings {
      data_sources = [
        "factory/line1/machines/+/telemetry",
        "factory/line1/machines/+/status",
        "factory/line1/quality/+"
      ]
      endpoint_ref         = "factory-mqtt"
      asset_ref           = "production-line-1"
      serialization_format = "Json"
    }
  }

  # Quality control data source
  operations {
    name           = "quality-source"
    operation_type = "Source"

    source_settings {
      data_sources = [
        "quality/inspections/+",
        "quality/defects/+",
        "quality/reports/+"
      ]
      endpoint_ref         = "quality-mqtt"
      asset_ref           = "quality-systems"
      serialization_format = "Json"
    }
  }

  # Advanced transformation
  operations {
    name           = "advanced-analytics"
    operation_type = "BuiltInTransformation"

    built_in_transformation_settings {
      serialization_format = "Parquet"

      # Filter for production hours only
      filter {
        description = "Include only production shift data"
        expression  = "hour(timestamp) >= 6 && hour(timestamp) <= 22"
        inputs      = ["timestamp"]
        type        = "Filter"
      }

      # Filter out maintenance modes
      filter {
        description = "Exclude machines in maintenance mode"
        expression  = "machineStatus != 'maintenance'"
        inputs      = ["machineStatus"]
        type        = "Filter"
      }

      # Calculate efficiency metrics
      map {
        description = "Calculate Overall Equipment Effectiveness (OEE)"
        expression  = "(actualOutput / plannedOutput) * (operatingTime / plannedTime) * (goodParts / totalParts)"
        inputs      = ["actualOutput", "plannedOutput", "operatingTime", "plannedTime", "goodParts", "totalParts"]
        output      = "oee_score"
        type        = "Compute"
      }

      # Categorize performance
      map {
        description = "Categorize machine performance"
        expression  = "oee_score >= 0.85 ? 'excellent' : oee_score >= 0.65 ? 'good' : oee_score >= 0.40 ? 'fair' : 'poor'"
        inputs      = ["oee_score"]
        output      = "performance_category"
        type        = "Compute"
      }

      # Add analytics metadata
      map {
        description = "Add analytics processing metadata"
        expression  = "{'pipeline_version': '2.1', 'processed_by': 'advanced-analytics'}"
        inputs      = []
        output      = "analytics_metadata"
        type        = "NewProperties"
      }

      # Normalize machine IDs
      map {
        description = "Standardize machine identifiers"
        expression  = "upper(replace(machineId, '-', '_'))"
        inputs      = ["machineId"]
        output      = "normalized_machine_id"
        type        = "Compute"
      }

      # Create performance dataset
      datasets {
        key         = "performance_summary"
        description = "Machine performance summary for dashboards"
        expression  = "{'machine': normalized_machine_id, 'oee': oee_score, 'category': performance_category, 'shift': shift}"
        inputs      = ["normalized_machine_id", "oee_score", "performance_category", "shift"]
        schema_ref  = "performance-schema"
      }

      # Create quality dataset
      datasets {
        key         = "quality_metrics"
        description = "Quality control metrics for analysis"
        expression  = "{'defect_rate': defectCount / totalParts, 'quality_grade': qualityScore}"
        inputs      = ["defectCount", "totalParts", "qualityScore"]
        schema_ref  = "quality-schema"
      }

      # Create alerts dataset
      datasets {
        key         = "production_alerts"
        description = "Production issues requiring immediate attention"
        expression  = "oee_score < 0.40 || defectCount > 10"
        inputs      = ["oee_score", "defectCount", "normalized_machine_id"]
        schema_ref  = "alert-schema"
      }
    }
  }

  # Data lake destination
  operations {
    name           = "datalake-analytics"
    operation_type = "Destination"

    destination_settings {
      data_destination = "manufacturing-analytics"
      endpoint_ref     = "datalake-endpoint"
    }
  }
}
```

### Real-time Streaming Pipeline

```hcl
resource "azurerm_iotoperations_dataflow" "streaming_pipeline" {
  name                    = "streaming-pipeline"
  resource_group_name     = azurerm_resource_group.example.name
  instance_name          = azurerm_iotoperations_instance.example.name
  dataflow_profile_name  = azurerm_iotoperations_dataflow_profile.example.name
  mode                   = "Enabled"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  # High-frequency sensor data
  operations {
    name           = "streaming-source"
    operation_type = "Source"

    source_settings {
      data_sources = [
        "realtime/vibration/+",
        "realtime/temperature/+",
        "realtime/pressure/+"
      ]
      endpoint_ref         = "edge-mqtt"
      serialization_format = "Json"
    }
  }

  # Real-time processing
  operations {
    name           = "realtime-transform"
    operation_type = "BuiltInTransformation"

    built_in_transformation_settings {
      serialization_format = "Delta"

      # Filter for anomalous readings
      filter {
        description = "Detect vibration anomalies"
        expression  = "vibrationLevel > threshold * 1.5"
        inputs      = ["vibrationLevel", "threshold"]
        type        = "Filter"
      }

      # Calculate moving averages
      map {
        description = "Calculate 5-minute moving average"
        expression  = "avg(temperature, 300)"  # 5-minute window
        inputs      = ["temperature"]
        output      = "temperature_5min_avg"
        type        = "BuiltInFunction"
      }

      # Detect trends
      map {
        description = "Detect increasing temperature trend"
        expression  = "slope(temperature, 600) > 0.1"  # 10-minute trend
        inputs      = ["temperature"]
        output      = "temperature_rising"
        type        = "BuiltInFunction"
      }

      # Add severity levels
      map {
        description = "Assign alert severity"
        expression  = "vibrationLevel > criticalThreshold ? 'critical' : vibrationLevel > warningThreshold ? 'warning' : 'info'"
        inputs      = ["vibrationLevel", "criticalThreshold", "warningThreshold"]
        output      = "alert_severity"
        type        = "Compute"
      }

      # Create streaming datasets
      datasets {
        key         = "critical_alerts"
        description = "Critical condition alerts for immediate response"
        expression  = "alert_severity == 'critical' || temperature_rising == true"
        inputs      = ["alert_severity", "temperature_rising", "deviceId", "timestamp"]
      }

      datasets {
        key         = "trend_analysis"
        description = "Trend data for predictive maintenance"
        expression  = "{'device': deviceId, 'trend': slope(temperature, 3600), 'variance': variance(vibrationLevel, 1800)}"
        inputs      = ["deviceId", "temperature", "vibrationLevel"]
      }
    }
  }

  # Real-time destination
  operations {
    name           = "realtime-output"
    operation_type = "Destination"

    destination_settings {
      data_destination = "realtime-alerts"
      endpoint_ref     = "stream-analytics"
    }
  }
}
```

### Multi-Format Data Integration

```hcl
resource "azurerm_iotoperations_dataflow" "integration_pipeline" {
  name                    = "integration-pipeline"
  resource_group_name     = azurerm_resource_group.example.name
  instance_name          = azurerm_iotoperations_instance.example.name
  dataflow_profile_name  = azurerm_iotoperations_dataflow_profile.example.name
  mode                   = "Enabled"

  extended_location {
    name = azurerm_iotoperations_instance.example.extended_location_name
    type = "CustomLocation"
  }

  # JSON data source
  operations {
    name           = "json-source"
    operation_type = "Source"

    source_settings {
      data_sources = [
        "legacy/systems/+/data"
      ]
      endpoint_ref         = "legacy-endpoint"
      schema_ref          = "legacy-json-schema"
      serialization_format = "Json"
    }
  }

  # Format transformation and standardization
  operations {
    name           = "format-standardization"
    operation_type = "BuiltInTransformation"

    built_in_transformation_settings {
      schema_ref           = "unified-schema"
      serialization_format = "Parquet"

      # Standardize timestamps
      map {
        description = "Convert various timestamp formats to ISO 8601"
        expression  = "parseTimestamp(timestamp_field, timestamp_format)"
        inputs      = ["timestamp_field", "timestamp_format"]
        output      = "standardized_timestamp"
        type        = "BuiltInFunction"
      }

      # Normalize measurement units
      map {
        description = "Convert all temperatures to Celsius"
        expression  = "unit == 'F' ? (value - 32) * 5/9 : value"
        inputs      = ["value", "unit"]
        output      = "temperature_celsius"
        type        = "Compute"
      }

      # Standardize device identifiers
      map {
        description = "Create unified device ID format"
        expression  = "concat(location, '_', deviceType, '_', serialNumber)"
        inputs      = ["location", "deviceType", "serialNumber"]
        output      = "unified_device_id"
        type        = "BuiltInFunction"
      }

      # Data quality scoring
      map {
        description = "Calculate data quality score"
        expression  = "(timestamp_valid ? 25 : 0) + (value_in_range ? 25 : 0) + (device_id_valid ? 25 : 0) + (schema_compliant ? 25 : 0)"
        inputs      = ["timestamp_valid", "value_in_range", "device_id_valid", "schema_compliant"]
        output      = "data_quality_score"
        type        = "Compute"
      }

      # Filter high-quality data
      filter {
        description = "Include only high-quality data records"
        expression  = "data_quality_score >= 75"
        inputs      = ["data_quality_score"]
        type        = "Filter"
      }

      # Create integration datasets
      datasets {
        key         = "master_data"
        description = "Master dataset with all standardized records"
        inputs      = ["unified_device_id", "standardized_timestamp", "temperature_celsius", "data_quality_score"]
        schema_ref  = "master-data-schema"
      }

      datasets {
        key         = "data_quality_report"
        description = "Data quality metrics for monitoring"
        expression  = "{'avg_quality': avg(data_quality_score), 'record_count': count(), 'processing_time': now()}"
        inputs      = ["data_quality_score"]
        schema_ref  = "quality-report-schema"
      }
    }
  }

  # Unified destination
  operations {
    name           = "unified-destination"
    operation_type = "Destination"

    destination_settings {
      data_destination = "enterprise-data-hub"
      endpoint_ref     = "enterprise-endpoint"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Dataflow. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the IoT Operations Dataflow should exist. Must be between 1-90 characters. Changing this forces a new resource to be created.

* `instance_name` - (Required) The name of the IoT Operations Instance. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `dataflow_profile_name` - (Required) The name of the Dataflow Profile that this Dataflow belongs to. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `mode` - (Optional) The operational mode of the dataflow. Possible values are `Enabled` and `Disabled`. Defaults to `Enabled`.

* `operations` - (Required) A list of `operations` blocks as defined below. At least one operation must be configured.

* `extended_location` - (Required) An `extended_location` block as defined below. Changing this forces a new resource to be created.

---

An `extended_location` block supports the following:

* `name` - (Required) The extended location name where the IoT Operations Dataflow should be deployed.

* `type` - (Required) The extended location type. Must be `CustomLocation`.

---

An `operations` block supports the following:

* `operation_type` - (Required) The type of operation. Possible values are `Source`, `Destination`, and `BuiltInTransformation`.

* `name` - (Optional) The name of the operation. Must be between 1-63 characters.

* `source_settings` - (Optional) A `source_settings` block as defined below. Required when `operation_type` is `Source`.

* `destination_settings` - (Optional) A `destination_settings` block as defined below. Required when `operation_type` is `Destination`.

* `built_in_transformation_settings` - (Optional) A `built_in_transformation_settings` block as defined below. Required when `operation_type` is `BuiltInTransformation`.

---

A `source_settings` block supports the following:

* `data_sources` - (Required) List of data source identifiers. At least one data source must be specified. Must be between 1-253 characters each.

* `endpoint_ref` - (Required) Reference to the endpoint configuration. Must be between 1-253 characters.

* `asset_ref` - (Optional) Reference to the asset configuration. Must be between 1-253 characters.

* `schema_ref` - (Optional) Reference to the schema definition. Must be between 1-253 characters.

* `serialization_format` - (Optional) The data serialization format. Currently only `Json` is supported.

---

A `destination_settings` block supports the following:

* `data_destination` - (Required) The data destination identifier. Must be between 1-253 characters.

* `endpoint_ref` - (Required) Reference to the endpoint configuration. Must be between 1-253 characters.

---

A `built_in_transformation_settings` block supports the following:

* `schema_ref` - (Optional) Reference to the output schema definition. Must be between 1-253 characters.

* `serialization_format` - (Optional) The output serialization format. Possible values are `Delta`, `Json`, and `Parquet`.

* `datasets` - (Optional) A list of `datasets` blocks as defined below for creating named datasets.

* `filter` - (Optional) A list of `filter` blocks as defined below for data filtering operations.

* `map` - (Optional) A list of `map` blocks as defined below for data transformation operations.

---

A `datasets` block supports the following:

* `key` - (Required) The unique key identifier for the dataset. Must be between 1-253 characters.

* `inputs` - (Required) List of input field names for the dataset.

* `description` - (Optional) Description of the dataset. Must be between 1-500 characters.

* `expression` - (Optional) Expression for creating the dataset. Must be between 1-1000 characters.

* `schema_ref` - (Optional) Reference to the dataset schema definition. Must be between 1-253 characters.

---

A `filter` block supports the following:

* `expression` - (Required) The filter expression. Must be between 1-1000 characters.

* `inputs` - (Required) List of input field names for the filter.

* `description` - (Optional) Description of the filter operation. Must be between 1-500 characters.

* `type` - (Optional) The type of filter. Currently only `Filter` is supported.

---

A `map` block supports the following:

* `inputs` - (Required) List of input field names for the mapping operation.

* `output` - (Required) The output field name. Must be between 1-253 characters.

* `description` - (Optional) Description of the mapping operation. Must be between 1-500 characters.

* `expression` - (Optional) The transformation expression. Must be between 1-1000 characters.

* `type` - (Optional) The type of mapping operation. Possible values are `BuiltInFunction`, `Compute`, `NewProperties`, `PassThrough`, and `Rename`.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ID of the IoT Operations Dataflow.

* `provisioning_state` - The provisioning state of the Dataflow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Operations Dataflow.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Operations Dataflow.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Operations Dataflow.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Operations Dataflow.

## Import

An IoT Operations Dataflow can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotoperations_dataflow.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.IoTOperations/instances/instance1/dataflowProfiles/profile1/dataflows/dataflow1
```
