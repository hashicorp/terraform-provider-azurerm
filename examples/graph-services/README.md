## Example: Graph Services Account

This example provisions a Graph Services Account using the current authenticated Azure client configuration.

### Prerequisites

Before running this example, you need:

1. An Azure subscription
2. Appropriate Azure CLI or PowerShell authentication with permissions to create Graph Services resources
3. The authenticated client must have a valid application ID (client ID)

### Variables

- `prefix` - A prefix used for all resources in this example (default: "example")
- `location` - The Azure region where all resources in this example should be created (default: "West Europe")

### Running the Example

1. Ensure you are authenticated to Azure:
   ```bash
   # Using Azure CLI
   az login
   
   # Or using PowerShell
   Connect-AzAccount
   ```

2. Optional: Set variables in `terraform.tfvars` or export them as environment variables:
   ```bash
   export TF_VAR_prefix="my-prefix"
   export TF_VAR_location="East US"
   ```

3. Initialize and apply the Terraform configuration:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

### What's Created

- **Resource Group** - A resource group to contain the Graph Services resources
- **Graph Services Account** - A Graph Services Account associated with the current authenticated client application

### Outputs

- `graph_services_account_id` - The ID of the created Graph Services Account
- `billing_plan_id` - The billing plan ID of the Graph Services Account