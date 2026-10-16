---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_access_policy"
description: |-
  Lists Key Vault Access Policy resources.
---

# List resource: azurerm_key_vault_access_policy

Lists Key Vault Access Policy resources for a given Key Vault.

## Example Usage

### List all Access Policies for a Key Vault

```hcl
list "azurerm_key_vault_access_policy" "example" {
  provider = azurerm
  config {
    key_vault_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.KeyVault/vaults/example-kv"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `key_vault_id` - (Required) The ID of the Key Vault whose Access Policies should be listed.
