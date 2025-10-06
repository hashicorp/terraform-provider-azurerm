# IoT Operations Instance

This example shows how to create an Azure IoT Operations instance using Terraform.

## Prerequisites

Before running this example, you need:

1. **Azure CLI** installed and authenticated
2. **Terraform** 1.6 or later
3. **Existing Resource Group** in Azure
4. **Arc-enabled Kubernetes cluster** with a Custom Location
5. **Schema Registry** in Azure Device Registry

## Usage

### Step 1: Set Variables

Create a `terraform.tfvars` file:

```hcl
# Prefix for resource naming
prefix = "mycompany"

# Existing Resource Group
resource_group_name = "existing-resource-group-name"

# Required Resource IDs (replace with your actual values)
custom_location_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ExtendedLocation/customLocations/example-location"
schema_registry_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.DeviceRegistry/schemaRegistries/example-registry"
```

### Step 2: Find Required Resources

Find available Custom Locations:
```bash
az resource list --resource-type "Microsoft.ExtendedLocation/customLocations" --query "[].{Name:name, Id:id}" -o table
```

Find available Schema Registries:
```bash
az resource list --resource-type "Microsoft.DeviceRegistry/schemaRegistries" --query "[].{Name:name, Id:id}" -o table
```

### Step 3: Deploy

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
| `custom_location_id` | ARM ID of Custom Location | `string` | yes |
| `schema_registry_id` | ARM ID of Schema Registry | `string` | yes |

## Outputs

| Name | Description |
|------|-------------|
| `iotoperations_instance_id` | ARM resource ID of the IoT Operations instance |

## Architecture

This example creates:

- **IoT Operations Instance** (named `{prefix}-iotoperations`) via ARM template deployment

The IoT Operations instance requires:
- An existing Resource Group
- An Arc-enabled Kubernetes cluster (Custom Location)
- A Schema Registry for data schemas

## Cleanup

```bash
terraform destroy
```

Note: This will only destroy the IoT Operations instance. The resource group, Custom Location, and Schema Registry will remain.