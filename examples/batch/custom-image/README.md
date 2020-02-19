## Example: Azure Batch with custom image

This example provisions the following Resources:

## Creates

1. A Resource Group
2. A [Batch Account](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#account)
3. A Custom Virtual Machine image to be used by the Azure Batch Pool
4. A [Batch pool that uses a custom VM image for virtual machines](https://docs.microsoft.com/en-us/azure/batch/batch-custom-images)

## Usage

- Provide values to all variables (credentials and names).
- Create with `terraform apply`
- Destroy all with `terraform destroy --force`
