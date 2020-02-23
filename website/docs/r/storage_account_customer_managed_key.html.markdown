---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_customer_managed_key"
description: |-
  Manages a Customer Managed Key for a Storage Account.
---

# azurerm_storage_account_customer_managed_key

Manages a Customer Managed Key for a Storage Account.

## Example Usage

```hcl
data "azurerm_key_vault" "example" {
  name                = "mykeyvault"
  resource_group_name = "some-resource-group"
}

data "azurerm_key_vault_key" "example" {
  name         = "secret-sauce"
  key_vault_id = data.azurerm_key_vault.example.id
}

resource "azurerm_storage_account" "tfex" {
  name                     = "exampleaccount"
  resource_group_name      = data.azurerm_key_vault.example.name
  location                 = data.azurerm_key_vault.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account_customer_managed_key" "example" {
  storage_account_id = azurerm_storage_account.example.id
  key_vault_id       = data.azurerm_key_vault.example.id
  key_name           = data.azurerm_key_vault_key.example.name
  key_version        = data.azurerm_key_vault_key.example.version
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the Storage Account. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of the Key Vault. Changing this forces a new resource to be created.

* `key_name` - (Required) The name of Key Vault Key.

* `key_version` - (Required) The version of Key Vault Key.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Account.

---

## Import

Customer Managed Keys for a Storage Account can be imported using the `resource id` of the Storage Account, e.g.

```shell
terraform import azurerm_storage_account_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```