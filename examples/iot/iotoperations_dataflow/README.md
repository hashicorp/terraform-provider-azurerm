# IoT Operations Dataflow Example

This example demonstrates how to create an Azure IoT Operations Dataflow using Terraform.

## Overview

This Terraform configuration creates:
- A Resource Group
- An IoT Operations Instance
- An IoT Operations Dataflow Profile
- An IoT Operations Dataflow with sources, destinations, transformations, and operations

## Prerequisites

- Azure subscription
- Terraform installed
- Azure CLI installed and authenticated
- An Azure Kubernetes Service (AKS) cluster for IoT Operations deployment
- IoT Operations endpoints (MQTT broker, Azure Data Explorer, etc.) configured

## Usage

1. Clone this repository and navigate to this example directory:
   ```bash
   cd examples/iot/iotoperations_dataflow
   ```

2. Copy the example variables file:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

3. Edit `terraform.tfvars` with your desired values:
   ```hcl
   resource_group_name = "my-iotops-rg"
   location           = "East US 2"
   instance_name      = "my-iotops-instance"
   dataflow_profile_name = "my-dataflow-profile"
   dataflow_name      = "temperature-processing-flow"
   ```

4. Initialize Terraform:
   ```bash
   terraform init
   ```

5. Plan the deployment:
   ```bash
   terraform plan
   ```

6. Apply the configuration:
   ```bash
   terraform apply
   ```

## Configuration

### Dataflow Profile

The dataflow profile manages the compute resources for dataflow operations:

- **Instance Count**: Number of dataflow instances to run
- **Log Level**: Logging verbosity (trace, debug, info, warn, error)
- **Metrics**: Prometheus metrics configuration

### Dataflow

The dataflow defines the data processing pipeline with these components:

#### Sources
Data input endpoints that can include:
- MQTT topics from broker endpoints
- Asset references for structured data
- Schema references for data validation
- Serialization format (JSON, Avro, etc.)

#### Destinations
Data output endpoints such as:
- Azure Data Explorer (ADX)
- Azure Blob Storage
- Event Hubs
- Custom endpoints

#### Transformations
Data processing operations:
- **Filter**: Conditional logic to filter data
- **Map**: Transform data structure and add computed fields
- **Aggregate**: Group and summarize data

#### Operations
High-level data flow operations that connect sources to destinations through transformations.

## Resource Hierarchy

The IoT Operations resources follow this hierarchy:
1. **Instance** - Top-level IoT Operations instance
2. **Dataflow Profile** - Compute resource manager for dataflows
3. **Dataflow** - Data processing pipeline definition

## Example Data Flow

This example creates a data flow that:
1. Reads temperature and humidity data from MQTT topics
2. Filters data based on temperature and humidity thresholds
3. Transforms the data by adding metadata and computed fields
4. Writes processed data to Azure Data Explorer
5. Archives raw data to Blob Storage

## Advanced Configuration

### Multiple Sources and Destinations

You can configure multiple data sources and destinations:

```hcl
dataflow_sources = [
  {
    name         = "sensor-data"
    endpoint_ref = "mqtt-endpoint"
    asset_ref    = "temperature-asset"
  },
  {
    name         = "machine-data"
    endpoint_ref = "opcua-endpoint"
    asset_ref    = "pressure-asset"
  }
]
```

### Complex Transformations

Build sophisticated data transformations:

```hcl
dataflow_transformations = [
  {
    type = "filter"
    filter = {
      expression = "temperature > 20 && status == 'active'"
      type      = "condition"
    }
  },
  {
    type = "map"
    map = {
      expression = "{ temp_f: temperature * 9/5 + 32, alert: temperature > 30 }"
      type      = "newProperties"
    }
  }
]
```

## Outputs

This configuration provides the following outputs:
- Resource Group ID
- IoT Operations Instance ID
- Dataflow Profile ID and configuration
- Dataflow ID and configuration details

## Clean Up

To destroy the resources:
```bash
terraform destroy
```

## Notes

- Ensure your IoT Operations instance is properly configured with the necessary endpoints
- Schema references must point to valid schema registry entries
- Endpoint references must match configured IoT Operations endpoints
- Transformation expressions use JSONPath and mathematical operations
- Monitor dataflow performance through Prometheus metrics

## More Information

For more details about Azure IoT Operations Dataflows, visit the [official documentation](https://docs.microsoft.com/azure/iot-operations/process-data/).