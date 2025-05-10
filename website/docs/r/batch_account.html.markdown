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
  name                                = "testbatchaccount"
  resource_group_name                 = azurerm_resource_group.example.name
  location                            = azurerm_resource_group.example.location
  pool_allocation_mode                = "BatchService"
  storage_account_id                  = azurerm_storage_account.example.id
  storage_account_authentication_mode = "StorageKeys"

  tags = {
    env = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Batch account. Only lowercase Alphanumeric characters allowed. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Batch account. Changing this forces a new resource to be created.

~> **Note:** To work around [a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/5574) this property is currently treated as case-insensitive. A future version of Terraform will require that the casing is correct.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `network_profile` - (Optional) A `network_profile` block as defined below.

* `pool_allocation_mode` - (Optional) Specifies the mode to use for pool allocation. Possible values are `BatchService` or `UserSubscription`. Defaults to `BatchService`.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for this server. Defaults to `true`.

~> **Note:** When using `UserSubscription` mode, an Azure KeyVault reference has to be specified. See `key_vault_reference` below.

~> **Note:** When using `UserSubscription` mode, the `Microsoft Azure Batch` service principal has to have `Contributor` role on your subscription scope, as documented [here](https://docs.microsoft.com/azure/batch/batch-account-create-portal#additional-configuration-for-user-subscription-mode).

* `key_vault_reference` - (Optional) A `key_vault_reference` block, as defined below, that describes the Azure KeyVault reference to use when deploying the Azure Batch account using the `UserSubscription` pool allocation mode.

* `storage_account_id` - (Optional) Specifies the storage account to use for the Batch account. If not specified, Azure Batch will manage the storage.

~> **Note:** When using `storage_account_id`, the `storage_account_authentication_mode` must be specified as well.

* `storage_account_authentication_mode` - (Optional) Specifies the storage account authentication mode. Possible values include `StorageKeys`, `BatchAccountManagedIdentity`.

~> **Note:** When using `BatchAccountManagedIdentity` mod, the `identity.type` must set to `UserAssigned` or `SystemAssigned`.

* `storage_account_node_identity` - (Optional) Specifies the user assigned identity for the storage account.

* `allowed_authentication_modes` - (Optional) Specifies the allowed authentication mode for the Batch account. Possible values include `AAD`, `SharedKey` or `TaskAuthenticationToken`.

* `encryption` - (Optional) Specifies if customer managed key encryption should be used to encrypt batch account data. One `encryption` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Batch Account. Possible values are `SystemAssigned` or `UserAssigned`.

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Batch Account.

~> **Note:** This is required when `type` is set to `UserAssigned`.

---

A `network_profile` block supports the following:

* `account_access` - (Optional) An `account_access` block as defined below.

* `node_management_access` - (Optional) A `node_management_access` block as defined below.

~> **Note:** At least one of `account_access` or `node_management_access` must be specified.

---

An `account_access` block supports the following:

* `default_action` - (Optional) Specifies the default action for the account access. Possible values are `Allow` and `Deny`. Defaults to `Deny`.

* `ip_rule` - (Optional) One or more `ip_rule` blocks as defined below.
---

A `node_management_access` block supports the following:

* `default_action` - (Optional) Specifies the default action for the node management access. Possible values are `Allow` and `Deny`. Defaults to `Deny`.

* `ip_rule` - (Optional) One or more `ip_rule` blocks as defined below.

---

An `ip_rule` block supports the following:

* `ip_range` - (Required) The CIDR block from which requests will match the rule.

* `action` - (Optional) Specifies the action of the ip rule. The only possible value is `Allow`. Defaults to `Allow`.

---

A `key_vault_reference` block supports the following:

* `id` - (Required) The Azure identifier of the Azure KeyVault to use.

* `url` - (Required) The HTTPS URL of the Azure KeyVault to use.

---

A `encryption` block supports the following:

* `key_vault_key_id` - (Required) The full URL path to the Azure key vault key id that should be used to encrypt data, as documented [here](https://docs.microsoft.com/azure/batch/batch-customer-managed-key). Both versioned and versionless keys are supported.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Batch Account.

* `identity` - An `identity` block as defined below.

* `primary_access_key` - The Batch account primary access key.

* `secondary_access_key` - The Batch account secondary access key.

* `account_endpoint` - The account endpoint used to interact with the Batch service.

~> **Note:** Primary and secondary access keys are only available when `pool_allocation_mode` is set to `BatchService` and `allowed_authentication_modes` contains `SharedKey`. See [documentation](https://docs.microsoft.com/azure/batch/batch-api-basics) for more information.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Batch Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Account.
* `update` - (Defaults to 30 minutes) Used when updating the Batch Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Batch Account.

## Import

Batch Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_batch_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Batch/batchAccounts/account1
```
