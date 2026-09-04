---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_account"
description: |-
  Get information about an existing Batch Account

---

# Data Source: azurerm_batch_account

Use this data source to access information about an existing Batch Account.

## Example Usage

```hcl
data "azurerm_batch_account" "example" {
  name                = "testbatchaccount"
  resource_group_name = "test"
}

output "pool_allocation_mode" {
  value = data.azurerm_batch_account.example.pool_allocation_mode
}
```

## Arguments Reference

* `name` - The name of the Batch account.

* `resource_group_name` - The Name of the Resource Group where this Batch account exists.

## Attributes Reference

The following attributes are exported:

* `allowed_authentication_modes` -  Specifies the allowed authentication mode for the Batch account. Possible values include `AAD`, `SharedKey` or `TaskAuthenticationToken`.

* `account_endpoint` - The account endpoint used to interact with the Batch service.

* `encryption` - The `encryption` block that describes the Azure KeyVault key reference used to encrypt data for the Azure Batch account.

* `id` - The Batch account ID.

* `identity` - An `identity` block as defined below.

* `key_vault_reference` - The `key_vault_reference` block that describes the Azure KeyVault reference to use when deploying the Azure Batch account using the `UserSubscription` pool allocation mode.

* `location` - The Azure Region in which this Batch account exists.

* `network_profile` - A `network_profile` block as defined below.

* `name` - The Batch account name.

* `pool_allocation_mode` - The pool allocation mode configured for this Batch account.

* `primary_access_key` - The Batch account primary access key.

* `public_network_access_enabled` - Whether public network access is allowed for this server. Defaults to `true`.

* `storage_account_id` - The ID of the Storage Account used for this Batch account.

* `storage_account_authentication_mode` - (Optional) Specifies the storage account authentication mode. Possible values include `StorageKeys`, `BatchAccountManagedIdentity`.

* `storage_account_node_identity` - (Optional) Specifies the user assigned identity for the storage account.

* `secondary_access_key` - The Batch account secondary access key.

* `tags` - A map of tags assigned to the Batch account.

~> **Note:** Primary and secondary access keys are only available when `pool_allocation_mode` is set to `BatchService`. See [documentation](https://docs.microsoft.com/azure/batch/batch-api-basics) for more information.

---

A `key_vault_reference` block exports the following:

* `id` - The Azure identifier of the Azure KeyVault reference.

* `url` - The HTTPS URL of the Azure KeyVault reference.

---

An `encryption` block exports the following:

* `key_vault_key_id` - The full URL path of the Key Vault Key used to encrypt data for this Batch account.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Account.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Batch` - 2024-07-01
