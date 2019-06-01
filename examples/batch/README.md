## Example: Azure Batch

This example provisions the following Resources:

## Creates

1. A Resource Group
2. A [Storage Account](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#azure-storage-account)
3. A [Batch Account](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#account)
4. Two [Batch pools](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#pool): one with fixed scale and the other with auto-scale.

## Usage

- Provide values to all variables (credentials and names).
- Create with `terraform apply`
- Destroy all with `terraform destroy --force`
