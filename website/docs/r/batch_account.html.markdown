---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_account"
description: |-
  Manages an Azure Batch account.

---

# azurerm_batch_account

Manages an Azure Batch account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "testbatch"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "teststorage"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "example" {
  name                 = "testbatchaccount"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.example.id

  tags = {
    env = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Batch account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Batch account. Changing this forces a new resource to be created.

~> **NOTE:** To work around [a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/5574) this property is currently treated as case-insensitive. A future version of Terraform will require that the casing is correct.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `pool_allocation_mode` - (Optional) Specifies the mode to use for pool allocation. Possible values are `BatchService` or `UserSubscription`. Defaults to `BatchService`.

~> **NOTE:** When using `UserSubscription` mode, an Azure KeyVault reference has to be specified. See `key_vault_reference` below.

~> **NOTE:** When using `UserSubscription` mode, the `Microsoft Azure Batch` service principal has to have `Contributor` role on your subscription scope, as documented [here](https://docs.microsoft.com/en-us/azure/batch/batch-account-create-portal#additional-configuration-for-user-subscription-mode).

* `key_vault_reference` - (Optional) A `key_vault_reference` block that describes the Azure KeyVault reference to use when deploying the Azure Batch account using the `UserSubscription` pool allocation mode.

* `storage_account_id` - (Optional) Specifies the storage account to use for the Batch account. If not specified, Azure Batch will manage the storage.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `key_vault_reference` block supports the following:

* `id` - (Required) The Azure identifier of the Azure KeyVault to use.

* `url` - (Required) The HTTPS URL of the Azure KeyVault to use.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Batch Account.

* `primary_access_key` - The Batch account primary access key.

* `secondary_access_key` - The Batch account secondary access key.

* `account_endpoint` - The account endpoint used to interact with the Batch service.

~> **NOTE:** Primary and secondary access keys are only available when `pool_allocation_mode` is set to `BatchService`. See [documentation](https://docs.microsoft.com/en-us/azure/batch/batch-api-basics) for more information.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Batch Account.
* `update` - (Defaults to 30 minutes) Used when updating the Batch Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Batch Account.

## Import

Batch Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_batch_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Batch/batchAccounts/account1
```
