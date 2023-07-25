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

## Argument Reference

* `name` - The name of the Batch account.

* `resource_group_name` - The Name of the Resource Group where this Batch account exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Batch account ID.

* `name` - The Batch account name.

* `location` - The Azure Region in which this Batch account exists.

* `pool_allocation_mode` - The pool allocation mode configured for this Batch account.

* `storage_account_id` - The ID of the Storage Account used for this Batch account.

* `primary_access_key` - The Batch account primary access key.

* `secondary_access_key` - The Batch account secondary access key.

* `account_endpoint` - The account endpoint used to interact with the Batch service.

* `key_vault_reference` - The `key_vault_reference` block that describes the Azure KeyVault reference to use when deploying the Azure Batch account using the `UserSubscription` pool allocation mode.

* `encryption` - The `encryption` block that describes the Azure KeyVault key reference used to encrypt data for the Azure Batch account.

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Account.
