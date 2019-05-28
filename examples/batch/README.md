## Example: Azure Batch

This example provisions the following Resources:

## Creates

1. A Resource Group
2. A [Storage Account](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#azure-storage-account)
3. A [Batch Account](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#account)
4. Two [Batch pools](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#pool): one with fixed scale and the other with auto-scale.
5. A [Batch pool that uses a custom VM image](https://docs.microsoft.com/en-us/azure/batch/batch-custom-images)

**Note:**: for point #5, it assumes that you have a VM image named `ubuntu1604base-img` in a resource group `batch-custom-img-rg` in the same Azure region than the one you are deploying the Azure Batch account. Check out the [documentation](https://docs.microsoft.com/en-us/azure/batch/batch-custom-images) for more information about how to create a custom VM image for Azure Batch.

## Usage

- Provide values to all variables (credentials and names).
- Create with `terraform apply`
- Destroy all with `terraform destroy --force`
