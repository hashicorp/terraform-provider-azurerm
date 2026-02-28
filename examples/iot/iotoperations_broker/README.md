# IoT Operations Broker

This example shows how to create an Azure IoT Operations broker using Terraform.

## Prerequisites

Before running this example, you need:

1. **Azure CLI** installed and authenticated
2. **Terraform** 1.6 or later
3. **Existing Resource Group** in Azure
4. **Existing IoT Operations Instance**
5. **Arc-enabled Kubernetes cluster** with a Custom Location

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

# Custom Location (Arc-enabled Kubernetes cluster)
custom_location_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ExtendedLocation/customLocations/example-location"
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
| `custom_location_id` | ARM ID of Custom Location | `string` | yes |

## Outputs

| Name | Description |
|------|-------------|
| `iotoperations_broker_id` | ARM resource ID of the IoT Operations broker |

## Architecture

This example creates:

- **IoT Operations Broker** (named `{prefix}-broker`) within an existing IoT Operations instance

The broker requires:
- An existing Resource Group
- An existing IoT Operations Instance
- An Arc-enabled Kubernetes cluster (Custom Location)

## Cleanup

```bash
terraform destroy
```

Note: This will only destroy the broker. The IoT Operations instance, resource group, and Custom Location will remain.