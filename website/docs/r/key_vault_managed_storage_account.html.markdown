---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_storage_account"
description: |-
  Manages a Key Vault Managed Storage Account.
---

# azurerm_key_vault_managed_storage_account

Manages a Key Vault Managed Storage Account.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_key_vault" "example" {
  name                = "keyvaultname"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Delete"
    ]

    storage_permissions = [
      "Get",
      "List",
      "Set",
      "SetSAS",
      "GetSAS",
      "DeleteSAS",
      "Update",
      "RegenerateKey"
    ]
  }
}

resource "azurerm_key_vault_managed_storage_account" "example" {
  name                         = "examplemanagedstorage"
  key_vault_id                 = azurerm_key_vault.example.id
  storage_account_id           = azurerm_storage_account.example.id
  storage_account_key          = "key1"
  regenerate_key_automatically = false
  regeneration_period          = "P1D"
}
```

## Example Usage (automatically regenerate Storage Account access key)

```hcl
data "azurerm_client_config" "current" {}

data "azuread_service_principal" "test" {
  # display_name = "Azure Key Vault"
  application_id = "cfa8b339-82a2-471a-a3c9-0fc0be7a4093"
}
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_key_vault" "example" {
  name                = "keyvaultname"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "Delete"
    ]

    storage_permissions = [
      "Get",
      "List",
      "Set",
      "SetSAS",
      "GetSAS",
      "DeleteSAS",
      "Update",
      "RegenerateKey"
    ]
  }
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Account Key Operator Service Role"
  principal_id         = data.azuread_service_principal.test.id
}

resource "azurerm_key_vault_managed_storage_account" "example" {
  name                         = "examplemanagedstorage"
  key_vault_id                 = azurerm_key_vault.example.id
  storage_account_id           = azurerm_storage_account.example.id
  storage_account_key          = "key1"
  regenerate_key_automatically = true
  regeneration_period          = "P1D"

  depends_on = [
    azurerm_role_assignment.example,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Key Vault Managed Storage Account. Changing this forces a new Key Vault Managed Storage Account to be created.

* `key_vault_id` - (Required) The ID of the Key Vault where the Managed Storage Account should be created. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) The ID of the Storage Account.

* `storage_account_key` - (Required) Which Storage Account access key that is managed by Key Vault. Possible values are `key1` and `key2`.

---

* `regenerate_key_automatically` - (Optional) Should Storage Account access key be regenerated periodically?

~> **Note:** Azure Key Vault application needs to have access to Storage Account for auto regeneration to work. Example can be found above.

* `regeneration_period` - (Optional) How often Storage Account access key should be regenerated. Value needs to be in [ISO 8601 duration format](https://en.wikipedia.org/wiki/ISO_8601#Durations).

* `tags` - (Optional) A mapping of tags which should be assigned to the Key Vault Managed Storage Account. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Key Vault Managed Storage Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Managed Storage Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Managed Storage Account.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Managed Storage Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Managed Storage Account.

## Import

Key Vault Managed Storage Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_storage_account.example https://example-keyvault.vault.azure.net/storage/exampleStorageAcc01
```
