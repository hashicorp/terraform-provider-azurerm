---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_customer_managed_key"
description: |-
  Manages a Customer Managed Key for a Cognitive Services Account.
---

# azurerm_cognitive_account_customer_managed_key

Manages a Customer Managed Key for a Cognitive Services Account.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_user_assigned_identity" "example" {
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  name = "example-identity"
}

resource "azurerm_cognitive_account" "example" {
  name                  = "example-account"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  kind                  = "Face"
  sku_name              = "E0"
  custom_subdomain_name = "example-account"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }
}


resource "azurerm_key_vault" "example" {
  name                     = "example-vault"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  soft_delete_enabled      = true
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "cognitive" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_cognitive_account.example.identity.0.tenant_id
  object_id    = azurerm_cognitive_account.example.identity.0.principal_id

  key_permissions    = ["get", "create", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_access_policy" "user" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_user_assigned_identity.example.tenant_id
  object_id    = azurerm_user_assigned_identity.example.principal_id

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_key" "example" {
  name         = "example-key"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.cognitive,
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.user,
  ]
}

resource "azurerm_cognitive_account_customer_managed_key" "example" {
  cognitive_account_id = azurerm_cognitive_account.example.id
  key_source           = "Microsoft.KeyVault"
  key_vault_key_id     = azurerm_key_vault_key.example.id
  identity_client_id   = azurerm_user_assigned_identity.example.client_id
}
```

## Arguments Reference

The following arguments are supported:

* `cognitive_account_id` - (Required) The ID of the Cognitive Account. Changing this forces a new resource to be created.

* `key_source` - (Required) The source for Encryption. Possible values include: `KeySourceMicrosoftCognitiveServices`, `KeySourceMicrosoftKeyVault`.

* `key_vault_key_id` - (Optional) The ID of the Key Vault Key.

* `identity_client_id` - (Optional) The client id of the user assigned identity that has access to the key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cognitive Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Account Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Account Customer Managed Key.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Account Customer Managed Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Account Customer Managed Key.

## Import

Customer Managed Keys for a Cognitive Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1
```
