# Azure Batch Sample

Sample to deploy a Job in Azure Batch

## Creates

1. A Resource Group
2. A [Storage Account](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#azure-storage-account)
3. A [Batch Account](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#account)
4. A [Pool](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#pool) of compute nodes

## Usage

- Provide values to all variables (credentials and names).
- Create with `terraform apply`
- Destroy all with `terraform destroy --force`
