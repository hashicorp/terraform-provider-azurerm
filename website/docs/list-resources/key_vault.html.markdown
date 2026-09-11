---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault"
description: |-
  Lists Key Vault resources.
---

# List resource: azurerm_key_vault

Lists Key Vault resources.

## Example Usage

### List all Key Vaults in the subscription

```hcl
list "azurerm_key_vault" "example" {
  provider = azurerm
  config {}
}
```

### List all Key Vaults in a specific resource group

```hcl
list "azurerm_key_vault" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
