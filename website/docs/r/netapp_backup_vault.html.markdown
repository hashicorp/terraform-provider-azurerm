---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_backup_vault"
description: |-
  Manages a NetApp Backup Vault.
---

# azurerm_netapp_backup_vault

Manages a NetApp Backup Vault.

## NetApp Backup Vault Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_netapp_account" "example" {
  name                = "example-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_backup_vault" "example" {
  name                = "example-netappbackupvault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  account_name        = azurerm_netapp_account.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Backup Vault. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Backup Vault should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp account in which the NetApp Vault should be created under. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the NetApp Backup Vault.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Backup Vault.
* `update` - (Defaults to 2 hours) Used when updating the NetApp Backup Vault.
* `delete` - (Defaults to 2 hours) Used when deleting the NetApp Backup Vault.

## Import

NetApp Backup Vault can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_backup_vault.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/backupVaults/backupVault1
```
