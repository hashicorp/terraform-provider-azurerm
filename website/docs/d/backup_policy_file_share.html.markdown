---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_backup_policy_file_share"
description: |-
  Gets information about an existing existing File Share Backup Policy.
---

# Data Source: azurerm_backup_policy_file_share

Use this data source to access information about an existing File Share Backup Policy.

## Example Usage

```hcl
data "azurerm_backup_policy_file_share" "policy" {
  name                = "policy"
  recovery_vault_name = "recovery_vault"
  resource_group_name = "resource_group"
}
```

## Argument Reference

The following arguments are supported:

- `name` - Specifies the name of the File Share Backup Policy.

- `recovery_vault_name` - Specifies the name of the Recovery Services Vault.

- `resource_group_name` - The name of the resource group in which the File Share Backup Policy resides.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the File Share Backup Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Services File Share Protection Policy.
