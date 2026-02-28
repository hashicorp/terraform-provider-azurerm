# IoT Operations Dataflow Endpoints Example

This example demonstrates how to create various Azure IoT Operations Dataflow Endpoints using Terraform.

## Overview

This Terraform configuration creates:
- A Resource Group
- An IoT Operations Instance
- Multiple Dataflow Endpoints for different data sources and destinations:
  - **MQTT Endpoint** - For IoT device data ingestion
  - **Azure Data Explorer (ADX) Endpoint** - For real-time analytics
  - **Azure Blob Storage Endpoint** - For data archival
  - **Local Storage Endpoint** - For temporary buffering
  - **Fabric OneLake Endpoint** - For Microsoft Fabric analytics (optional)

## Prerequisites

- Azure subscription
- Terraform installed
- Azure CLI installed and authenticated
- An Azure Kubernetes Service (AKS) cluster for IoT Operations deployment
- Required Azure services configured:
  - Azure Data Explorer cluster (for ADX endpoint)
  - Azure Storage Account (for blob storage endpoint)
  - Microsoft Fabric workspace (for Fabric endpoint, optional)

## Usage

1. Clone this repository and navigate to this example directory:
   ```bash
   cd examples/iot/iotoperations_dataflow_endpoint
   ```

2. Copy the example variables file:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

3. Edit `terraform.tfvars` with your specific values:
   ```hcl
   resource_group_name = "my-iotops-rg"
   location           = "East US 2"
   instance_name      = "my-iotops-instance"
   
   # Update endpoints with your actual service URLs
   adx_cluster_uri        = "https://myiotcluster.eastus2.kusto.windows.net"
   storage_account_host   = "https://myiotdatalake.blob.core.windows.net"
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

## Endpoint Types and Configuration

### MQTT Endpoint

The MQTT endpoint is used for ingesting data from IoT devices via MQTT protocol:

```hcl
endpoint_type = "Mqtt"
mqtt_settings {
  host = "mqtt-broker.iotoperations.svc.cluster.local"
  port = 1883
  
  authentication {
    method = "UsernamePassword"
    username = "iot-client"
    password_secret_name = "mqtt-credentials"
  }
  
  tls {
    mode = "Enabled"
    trusted_ca_certificate_config_map = "mqtt-ca-cert"
  }
}
```

**Features:**
- Username/password authentication
- TLS encryption support
- Configurable QoS levels
- Keep-alive and session management
- Message retention options

### Azure Data Explorer (ADX) Endpoint

The ADX endpoint sends processed data to Azure Data Explorer for real-time analytics:

```hcl
endpoint_type = "DataExplorer"
data_explorer_settings {
  host     = "https://mycluster.eastus2.kusto.windows.net"
  database = "iottelemetry"
  
  authentication {
    method = "SystemAssignedManagedIdentity"
    system_assigned_managed_identity_audience = "https://kusto.windows.net"
  }
  
  batching {
    latency_seconds = 5
    max_messages    = 1000
  }
}
```

**Features:**
- Managed Identity authentication
- Configurable batching for performance
- Direct integration with KQL queries
- Real-time data ingestion

### Azure Blob Storage Endpoint

The storage endpoint archives data to Azure Blob Storage:

```hcl
endpoint_type = "DataLakeStorage"
data_lake_storage_settings {
  host           = "https://mystorageaccount.blob.core.windows.net"
  container_name = "iotdata"
  
  authentication {
    method = "SystemAssignedManagedIdentity"
    system_assigned_managed_identity_audience = "https://storage.azure.com"
  }
  
  batching {
    latency_seconds = 60
    max_messages    = 10000
  }
}
```

**Features:**
- Long-term data archival
- Large batch processing
- Cost-effective storage
- Integration with Azure Data Lake

### Local Storage Endpoint

The local storage endpoint provides temporary data buffering:

```hcl
endpoint_type = "LocalStorage"
local_storage_settings {
  persistent_volume_claim_ref = "iot-local-storage"
}
```

**Features:**
- High-speed local caching
- Edge computing scenarios
- Kubernetes PVC integration
- Data resilience during network outages

### Fabric OneLake Endpoint (Optional)

The Fabric endpoint integrates with Microsoft Fabric for advanced analytics:

```hcl
endpoint_type = "FabricOneLake"
fabric_one_lake_settings {
  host             = "https://onelake.dfs.fabric.microsoft.com"
  workspace_id     = "12345678-1234-5678-9abc-123456789012"
  lakehouse_name   = "iotlakehouse"
  
  authentication {
    method = "SystemAssignedManagedIdentity"
    system_assigned_managed_identity_audience = "https://onelake.dfs.fabric.microsoft.com"
  }
}
```

**Features:**
- Microsoft Fabric integration
- Advanced analytics and AI capabilities
- Delta Lake format support
- Power BI integration

## Authentication and Security

All endpoints support Managed Identity authentication for secure, passwordless connections:

- **System Assigned Managed Identity**: Automatically created and managed by Azure
- **Audience-specific authentication**: Each service has its specific audience URL
- **TLS encryption**: Secure data transmission
- **Secret management**: Kubernetes secrets for sensitive data

## Batching Configuration

Endpoints support batching for optimal performance:

- **Latency**: Maximum time to wait before sending a batch
- **Max Messages**: Maximum number of messages per batch
- **Performance tuning**: Balance between latency and throughput

## Resource Hierarchy

The IoT Operations resources follow this hierarchy:
1. **Instance** - Top-level IoT Operations instance
2. **Dataflow Endpoints** - Data source and destination definitions
3. **Dataflows** - Processing pipelines that connect endpoints

## Example Data Flow Architecture

```
IoT Devices → MQTT Endpoint → Dataflow Processing → {
                                                     ├── ADX Endpoint (Real-time analytics)
                                                     ├── Storage Endpoint (Long-term archive)
                                                     ├── Local Endpoint (Edge caching)
                                                     └── Fabric Endpoint (Advanced analytics)
                                                   }
```

## Outputs

This configuration provides comprehensive outputs including:
- Individual endpoint IDs and names
- Endpoint configuration summaries
- Connection details for integration with dataflows

## Clean Up

To destroy the resources:
```bash
terraform destroy
```

## Notes

- Ensure all referenced Azure services exist before deployment
- Managed Identity requires appropriate RBAC permissions on target services
- Test endpoint connectivity before creating dataflows
- Monitor endpoint performance and adjust batching settings as needed
- Local storage requires properly configured Kubernetes PVCs

## Troubleshooting

Common issues and solutions:

1. **Authentication failures**: Verify Managed Identity permissions
2. **Connection timeouts**: Check network connectivity and firewall rules
3. **Batching issues**: Adjust latency and message count settings
4. **Storage full**: Monitor local storage usage and implement cleanup policies

## More Information

For more details about Azure IoT Operations Dataflow Endpoints, visit the [official documentation](https://docs.microsoft.com/azure/iot-operations/connect-to-cloud/).