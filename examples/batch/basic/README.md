## Example: Azure Batch

This example provisions the following Resources:

## Creates

1. A Resource Group
2. A [Storage Account](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#azure-storage-account)
3. A [Batch Account](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#account)
4. Two [Batch pools](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#pool): one with fixed scale and the other with auto-scale

## Usage

- Provide values to all variables (credentials and names).
- Create with `terraform apply`
- Destroy all with `terraform destroy --force`

## Example Usage with User Subscription mode

It's also possible to deploy Azure Batch Account in User Subscription mode. In this mode, all the machines that will be created by batch pools will be created in the user Azure subscription. In this mode, you need to specify a reference to an Azure Key Vault that will be used by Azure Batch to store and retrieve sensitive information. You can read more about User Subscription mode in Azure Batch on [this page](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics#account).

~> **NOTE:** the script below uses the also the [AzureAD provider](https://www.terraform.io/docs/providers/azuread/) to retrieve the "Microsoft Azure Batch" service principal information.


```hcl
resource "azurerm_resource_group" "example" {
  name     = "batch-rg"
  location = "westeurope"
}

# Get Microsoft Azure Batch service principal reference using the AzureAD provider
data "azuread_service_principal" "batchsp" {
  display_name = "Microsoft Azure Batch"
}

resource "azurerm_key_vault" "example" {
  name                            = "batchkv"
  location                        = "${azurerm_resource_group.example.location}"
  resource_group_name             = "${azurerm_resource_group.example.name}"
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  tenant_id                       = "00000000-0000-0000-0000-000000000002"

  sku_name = "standard"

  access_policy {
    tenant_id = "00000000-0000-0000-0000-000000000002"
    object_id = "${data.azuread_service_principal.batchsp.object_id}"

    secret_permissions = [
      "get",
      "list",
      "set",
      "delete"
    ]

  }
}

resource "azurerm_role_assignment" "contribrole" {
  scope                = "/subscriptions/00000000-0000-0000-0000-000000000001"
  role_definition_name = "Contributor"
  principal_id         = "${data.azuread_service_principal.batchsp.object_id}"
}

resource "azurerm_batch_account" "example" {
  name                 = "batchaccount"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  location             = "${azurerm_resource_group.example.location}"
  storage_account_id   = "${azurerm_storage_account.example.id}"
  pool_allocation_mode = "UserSubscription"
  
  # reference the Azure KeyVault
  key_vault_reference {
    id  = "${azurerm_key_vault.example.id}"
    url = "${azurerm_key_vault.example.vault_uri}"
  }
}
```
