# IoT Operations Broker Authorization

This example shows how to create an Azure IoT Operations broker authorization using Terraform.

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
```

### Step 2: Deploy

```bash
terraform init
terraform plan
terraform apply
```

## Variables

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `prefix` | Prefix for resource naming | `string` | yes |
| `resource_group_name` | Name of existing resource group | `string` | yes |
| `instance_name` | Name of existing IoT Operations instance | `string` | yes |
| `broker_name` | Name of existing IoT Operations broker | `string` | yes |

## Outputs

| Name | Description |
|------|-------------|
| `iotoperations_broker_authorization_id` | ARM resource ID of the IoT Operations broker authorization |

## Architecture

This example creates:

- **IoT Operations Broker Authorization** (named `{prefix}-broker-authz`) within an existing IoT Operations broker

The broker authorization requires:
- An existing Resource Group
- An existing IoT Operations Instance
- An existing IoT Operations Broker

## Authorization Policies

The example configures:
- **Cache**: Enabled for performance
- **Rules**: Authorization rules for broker access
  - **Broker Resources**: `["*"]` (all resources)
  - **Method**: `Connect` (connection method)
  - **Clients**: `["*"]` (all clients)
  - **State Store**: Key-value pairs for additional authorization context

## Resource Hierarchy

```
Resource Group
└── IoT Operations Instance
    └── IoT Operations Broker
        ├── IoT Operations Broker Authentication
        └── IoT Operations Broker Authorization  ← This resource
```

## Cleanup

```bash
terraform destroy
```

Note: This will only destroy the broker authorization. The broker, IoT Operations instance, and resource group will remain.