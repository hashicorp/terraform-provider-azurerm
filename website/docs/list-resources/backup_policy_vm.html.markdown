---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_policy_vm"
description: |-
  Lists Backup VM Backup Policy resources.
---

# List resource: azurerm_backup_policy_vm

Lists Backup VM Backup Policy resources.

## Example Usage

### List all Backup VM Backup Policies in a Recovery Services Vault

```hcl
list "azurerm_backup_policy_vm" "example" {
  provider = azurerm
  config {
    recovery_vault_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.RecoveryServices/vaults/example-vault"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `recovery_vault_id` - (Required) The ID of the Recovery Services Vault to query.
