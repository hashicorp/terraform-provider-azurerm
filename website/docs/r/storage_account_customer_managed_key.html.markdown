---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_customer_managed_key"
description: |-
  Manages a Customer Managed Key for a Storage Account.
---

# azurerm_storage_account_customer_managed_key

Manages a Customer Managed Key for a Storage Account.

~> **Note:** It's possible to define a Customer Managed Key both within [the `azurerm_storage_account` resource](storage_account.html) via the `customer_managed_key` block and by using [the `azurerm_storage_account_customer_managed_key` resource](storage_account_customer_managed_key.html). However it's not possible to use both methods to manage a Customer Managed Key for a Storage Account, since there'll be conflicts.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                = "examplekv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "storage" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_storage_account.example.identity[0].principal_id

  secret_permissions = ["Get"]
  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey"
  ]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  secret_permissions = ["Get"]
  key_permissions = [
    "Get",
    "Create",
    "Delete",
    "List",
    "Restore",
    "Recover",
    "UnwrapKey",
    "WrapKey",
    "Purge",
    "Encrypt",
    "Decrypt",
    "Sign",
    "Verify",
    "GetRotationPolicy",
    "SetRotationPolicy"
  ]
}


resource "azurerm_key_vault_key" "example" {
  name         = "tfex-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey"
  ]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.storage
  ]
}


resource "azurerm_storage_account" "example" {
  name                     = "examplestor"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  identity {
    type = "SystemAssigned"
  }

  lifecycle {
    ignore_changes = [
      customer_managed_key
    ]
  }
}

resource "azurerm_storage_account_customer_managed_key" "example" {
  storage_account_id = azurerm_storage_account.example.id
  key_vault_id       = azurerm_key_vault.example.id
  key_name           = azurerm_key_vault_key.example.name
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the Storage Account. Changing this forces a new resource to be created.

* `key_name` - (Required) The name of Key Vault Key.

* `key_vault_id` - (Optional) The ID of the Key Vault. Exactly one of `managed_hsm_key_id`, `key_vault_id`, or `key_vault_uri` must be specified.

~> **Note:** When the principal running Terraform has access to the subscription containing the Key Vault, it's recommended to use the `key_vault_id` property for maximum compatibility, rather than the `key_vault_uri` property. However if the Key Vault is in a different subscription from the Storage Account `key_vault_uri` will need to be used instead otherwise there will be a diff.

* `key_vault_uri` - (Optional) URI pointing at the Key Vault. Required when using `federated_identity_client_id`. Exactly one of `managed_hsm_key_id`, `key_vault_id`, or `key_vault_uri` must be specified.

* `managed_hsm_key_id` - (Optional) Key ID of a key in a managed HSM.  Exactly one of `managed_hsm_key_id`, `key_vault_id`, or `key_vault_uri` must be specified.

* `key_version` - (Optional) The version of Key Vault Key. Remove or omit this argument to enable Automatic Key Rotation.

* `user_assigned_identity_id` - (Optional) The ID of a user assigned identity.

* `federated_identity_client_id` - (Optional) The Client ID of the multi-tenant application to be used in conjunction with the user-assigned identity for cross-tenant customer-managed-keys server-side encryption on the storage account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Account.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Account Customer Managed Keys.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Customer Managed Keys.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Account Customer Managed Keys.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Account Customer Managed Keys.

## Import

Customer Managed Keys for a Storage Account can be imported using the `resource id` of the Storage Account, e.g.

```shell
terraform import azurerm_storage_account_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
