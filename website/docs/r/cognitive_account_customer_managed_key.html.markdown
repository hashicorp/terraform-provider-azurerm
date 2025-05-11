---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_customer_managed_key"
description: |-
  Manages a Customer Managed Key for a Cognitive Services Account.
---

# azurerm_cognitive_account_customer_managed_key

Manages a Customer Managed Key for a Cognitive Services Account.

~> **Note:** It's possible to define a Customer Managed Key both within [the `azurerm_cognitive_account` resource](cognitive_account.html) via the `customer_managed_key` block and by using [the `azurerm_cognitive_account_customer_managed_key` resource](cognitive_account_customer_managed_key.html). However it's not possible to use both methods to manage a Customer Managed Key for a Cognitive Account, since there'll be conflicts.

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
  name                = "example-identity"
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
  purge_protection_enabled = true

  access_policy {
    tenant_id = azurerm_cognitive_account.example.identity[0].tenant_id
    object_id = azurerm_cognitive_account.example.identity[0].principal_id
    key_permissions = [
      "Get", "Create", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"
    ]
    secret_permissions = [
      "Get",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"
    ]
    secret_permissions = [
      "Get",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.example.tenant_id
    object_id = azurerm_user_assigned_identity.example.principal_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify"
    ]
    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "example-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_cognitive_account_customer_managed_key" "example" {
  cognitive_account_id = azurerm_cognitive_account.example.id
  key_vault_key_id     = azurerm_key_vault_key.example.id
  identity_client_id   = azurerm_user_assigned_identity.example.client_id
}
```

## Arguments Reference

The following arguments are supported:

* `cognitive_account_id` - (Required) The ID of the Cognitive Account. Changing this forces a new resource to be created.

* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Encrypt the data in this Cognitive Account.

* `identity_client_id` - (Optional) The Client ID of the User Assigned Identity that has access to the key. This property only needs to be specified when there're multiple identities attached to the Cognitive Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Account Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Account Customer Managed Key.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Account Customer Managed Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Account Customer Managed Key.

## Import

Customer Managed Keys for a Cognitive Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1
```
