# IoT Operations Dataflow Profiles Example

This example demonstrates how to create multiple Azure IoT Operations Dataflow Profiles using Terraform for different processing scenarios.

## Overview

This Terraform configuration creates:
- A Resource Group
- An IoT Operations Instance
- Multiple Dataflow Profiles optimized for different use cases:
  - **High Performance Profile** - For real-time, high-throughput processing
  - **Standard Profile** - For regular batch processing workloads
  - **Edge Profile** - For resource-constrained edge environments
  - **Development Profile** - For testing and development (optional)

## Prerequisites

- Azure subscription
- Terraform installed
- Azure CLI installed and authenticated
- An Azure Kubernetes Service (AKS) cluster for IoT Operations deployment

## Usage

1. Clone this repository and navigate to this example directory:
   ```bash
   cd examples/iot/iotoperations_dataflow_profile
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
   
   # Configure profiles based on your workload requirements
   high_performance_instance_count = 4
   standard_instance_count        = 2
   edge_instance_count           = 1
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

## Dataflow Profile Types

### High Performance Profile

Optimized for real-time, high-throughput data processing:

```hcl
instance_count = 4
diagnostics {
  logs {
    level = "warn"  # Reduced logging for performance
  }
  metrics {
    prometheus_port = 9090
  }
  self_check {
    mode                = "Enabled"
    interval_seconds    = 30    # Frequent health checks
    timeout_seconds     = 15    # Quick timeout
  }
}
```

**Use Cases:**
- Real-time analytics
- High-frequency sensor data processing
- Low-latency alerting systems
- Stream processing applications

**Characteristics:**
- Multiple instances for horizontal scaling
- Reduced logging to minimize performance impact
- Frequent health checks for reliability
- Optimized for throughput over resource usage

### Standard Profile

Balanced configuration for regular batch processing:

```hcl
instance_count = 2
diagnostics {
  logs {
    level = "info"  # Standard logging level
  }
  metrics {
    prometheus_port = 9091
  }
  self_check {
    mode                = "Enabled"
    interval_seconds    = 60    # Regular health checks
    timeout_seconds     = 30    # Standard timeout
  }
}
```

**Use Cases:**
- Batch data processing
- ETL operations
- Scheduled data transformations
- General-purpose data flows

**Characteristics:**
- Moderate instance count for balanced performance
- Standard logging for operational visibility
- Regular health checks
- Good balance of performance and resource usage

### Edge Profile

Optimized for resource-constrained edge environments:

```hcl
instance_count = 1
diagnostics {
  logs {
    level = "error"  # Minimal logging
  }
  metrics {
    prometheus_port = 9092
  }
  self_check {
    mode                = "Enabled"
    interval_seconds    = 120   # Less frequent checks
    timeout_seconds     = 60    # Longer timeout
  }
}
```

**Use Cases:**
- Edge computing scenarios
- IoT gateway processing
- Remote site data processing
- Resource-limited environments

**Characteristics:**
- Single instance to minimize resource usage
- Error-only logging to reduce I/O
- Less frequent health checks to save resources
- Optimized for minimal resource consumption

### Development Profile (Optional)

Configuration for testing and development:

```hcl
instance_count = 1
diagnostics {
  logs {
    level = "debug"  # Verbose logging for debugging
  }
  metrics {
    prometheus_port = 9093
  }
  self_check {
    mode                = "Enabled"
    interval_seconds    = 30    # Frequent checks for development
    timeout_seconds     = 15    # Quick feedback
  }
}
```

**Use Cases:**
- Development and testing
- Debugging data flows
- Proof of concept implementations
- Learning and experimentation

**Characteristics:**
- Single instance for simplicity
- Debug-level logging for detailed insights
- Frequent health checks for quick feedback
- Easy to enable/disable via variable

## Configuration Options

### Instance Count

Controls the number of dataflow processing instances:
- **1**: Minimal resource usage, suitable for light workloads
- **2-3**: Balanced performance for moderate workloads
- **4+**: High performance for demanding workloads

### Log Levels

Available logging levels in order of verbosity:
- **trace**: Most verbose, includes all operations
- **debug**: Detailed information for debugging
- **info**: General operational information
- **warn**: Warning messages and above
- **error**: Error messages only (minimal)

### Self-Check Configuration

Health monitoring settings:
- **Mode**: Enable/disable health checks
- **Interval**: How often to perform health checks
- **Timeout**: Maximum time to wait for health check response

### Metrics

Prometheus metrics configuration:
- **Port**: Different ports for each profile to avoid conflicts
- **Endpoint**: Accessible at `http://localhost:{port}/metrics`

## Monitoring and Observability

Each profile exposes metrics on different ports:
- High Performance: `:9090/metrics`
- Standard: `:9091/metrics`
- Edge: `:9092/metrics`
- Development: `:9093/metrics`

Monitor key metrics:
- Processing throughput
- Error rates
- Resource utilization
- Health check status

## Resource Hierarchy

The IoT Operations resources follow this hierarchy:
1. **Instance** - Top-level IoT Operations instance
2. **Dataflow Profiles** - Compute resource managers
3. **Dataflows** - Processing pipelines that use profiles
4. **Dataflow Endpoints** - Data sources and destinations

## Profile Selection Strategy

Choose profiles based on your requirements:

| Requirement | Recommended Profile |
|-------------|-------------------|
| Real-time processing | High Performance |
| Batch processing | Standard |
| Edge computing | Edge |
| Development/Testing | Development |
| Mixed workloads | Multiple profiles |

## Best Practices

1. **Resource Planning**: Size instance counts based on expected load
2. **Monitoring**: Use different Prometheus ports for each profile
3. **Log Management**: Balance logging verbosity with performance needs
4. **Health Checks**: Adjust intervals based on environment reliability
5. **Environment Separation**: Use different profiles for dev/test/prod

## Outputs

This configuration provides comprehensive outputs including:
- Individual profile IDs and names
- Configuration summaries for all profiles
- Prometheus metrics endpoints
- Resource hierarchy information

## Clean Up

To destroy the resources:
```bash
terraform destroy
```

## Notes

- Profiles define compute resources but don't process data by themselves
- Create dataflows that reference these profiles for actual data processing
- Monitor resource usage to optimize instance counts
- Adjust self-check intervals based on your reliability requirements
- Use appropriate log levels for your operational needs

## Troubleshooting

Common issues and solutions:

1. **Resource constraints**: Reduce instance counts or use edge profile
2. **Port conflicts**: Ensure each profile uses a unique Prometheus port
3. **Performance issues**: Increase instance count or use high-performance profile
4. **Monitoring gaps**: Verify Prometheus metrics endpoints are accessible

## More Information

For more details about Azure IoT Operations Dataflow Profiles, visit the [official documentation](https://docs.microsoft.com/azure/iot-operations/process-data/).