# IoT Operations Broker Listener Example

This example demonstrates how to create an Azure IoT Operations Broker Listener using Terraform.

## Overview

This Terraform configuration creates:
- A Resource Group
- An IoT Operations Instance
- An IoT Operations Broker
- An IoT Operations Broker Listener with configurable options

## Prerequisites

- Azure subscription
- Terraform installed
- Azure CLI installed and authenticated
- An Azure Kubernetes Service (AKS) cluster for IoT Operations deployment

## Usage

1. Clone this repository and navigate to this example directory:
   ```bash
   cd examples/iot/iotoperations_broker_listener
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
   broker_name        = "my-broker"
   listener_name      = "mqtt-listener"
   listener_port      = 1883
   service_type       = "LoadBalancer"
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

### Broker Listener

The broker listener is configured with the following key parameters:

- **Port**: The port number on which the listener will accept connections (default: 1883 for MQTT)
- **Service Type**: The Kubernetes service type (ClusterIP, NodePort, LoadBalancer)
- **Service Name**: The name of the Kubernetes service

### TLS Configuration

TLS can be enabled for secure communication:

```hcl
enable_tls = true
tls_mode  = "Automatic"  # or "Manual"
```

When TLS is enabled, the listener will use cert-manager to automatically provision and manage certificates.

### Authentication and Authorization

The listener can reference authentication and authorization policies:

```hcl
authentication_ref_name = "my-auth-policy"
authorization_ref_name  = "my-authz-policy"
```

These references link to separately created authentication and authorization resources.

## Resource Hierarchy

The IoT Operations resources follow this hierarchy:
1. **Instance** - Top-level IoT Operations instance
2. **Broker** - MQTT broker within the instance
3. **Listener** - Network endpoint for the broker

Each listener belongs to a specific broker, which belongs to a specific instance.

## Outputs

This configuration provides the following outputs:
- Resource Group ID
- IoT Operations Instance ID
- Broker ID
- Broker Listener ID and configuration details

## Clean Up

To destroy the resources:
```bash
terraform destroy
```

## Notes

- The listener port should not conflict with other services in your cluster
- When using LoadBalancer service type, ensure your cluster supports external load balancers
- TLS configuration requires cert-manager to be installed in your cluster
- Authentication and authorization policies must be created separately if referenced

## More Information

For more details about Azure IoT Operations, visit the [official documentation](https://docs.microsoft.com/azure/iot-operations/).