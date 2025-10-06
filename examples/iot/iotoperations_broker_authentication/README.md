# IoT Operations Broker Authentication

This example shows how to create an Azure IoT Operations broker authentication using Terraform.

## Prerequisites

Before running this example, you need:

1. **Azure CLI** installed and authenticated
2. **Terraform** 1.6 or later
3. **Existing Resource Group** in Azure
4. **Existing IoT Operations Instance**
5. **Existing IoT Operations Broker**

## Usage

### Step 1: Set Variables

Create a `terraform.tfvars` file:

```hcl
# Prefix for resource naming
prefix = "mycompany"

# Existing Resource Group
resource_group_name = "existing-resource-group-name"

# Existing IoT Operations Instance
instance_name = "existing-iotoperations-instance"

# Existing IoT Operations Broker
broker_name = "existing-iotoperations-broker"

# Authentication audience (optional)
audience = "aio-internal"
```

### Step 2: Deploy

```bash
terraform init
terraform plan
terraform apply
```

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| `prefix` | Prefix for resource naming | `string` | n/a | yes |
| `resource_group_name` | Name of existing resource group | `string` | n/a | yes |
| `instance_name` | Name of existing IoT Operations instance | `string` | n/a | yes |
| `broker_name` | Name of existing IoT Operations broker | `string` | n/a | yes |
| `audience` | Authentication audience | `string` | `"aio-internal"` | no |

## Outputs

| Name | Description |
|------|-------------|
| `iotoperations_broker_authentication_id` | ARM resource ID of the IoT Operations broker authentication |

## Architecture

This example creates:

- **IoT Operations Broker Authentication** (named `{prefix}-broker-auth`) within an existing IoT Operations broker

The broker authentication requires:
- An existing Resource Group
- An existing IoT Operations Instance
- An existing IoT Operations Broker

## Authentication Methods

The example configures:
- **ServiceAccountToken** authentication method
- **Custom settings** with audience configuration

## Cleanup

```bash
terraform destroy
```

Note: This will only destroy the broker authentication. The broker, IoT Operations instance, and resource group will remain.