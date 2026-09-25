---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_vault"
description: |-
    Lists Data Protection Backup Vault resources.
---

# List resource: azurerm_data_protection_backup_vault

Lists Data Protection Backup Vault resources.

## Example Usage

### List all Data Protection Backup Vaults in the subscription

```hcl
list "azurerm_data_protection_backup_vault" "example" {
  provider = azurerm
  config {}
}
```

### List all Data Protection Backup Vaults in a specific resource group

```hcl
list "azurerm_data_protection_backup_vault" "example" {
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
