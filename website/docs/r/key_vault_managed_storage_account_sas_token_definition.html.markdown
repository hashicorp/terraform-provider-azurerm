---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_storage_account_sas_token_definition"
description: |-
  Manages a Key Vault Managed Storage Account SAS Definition.
---

# azurerm_key_vault_managed_storage_account_sas_token_definition

Manages a Key Vault Managed Storage Account SAS Definition.

## Example Usage

```hcl
data "azurerm_client_config" "example" {}

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

data "azurerm_storage_account_sas" "example" {
  connection_string = azurerm_storage_account.example.primary_connection_string
  https_only        = true

  resource_types {
    service   = true
    container = false
    object    = false
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2021-04-30T00:00:00Z"
  expiry = "2023-04-30T00:00:00Z"

  permissions {
    read    = true
    write   = true
    delete  = false
    list    = false
    add     = true
    create  = true
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

resource "azurerm_key_vault" "example" {
  name                = "example-keyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.example.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.example.tenant_id
    object_id = data.azurerm_client_config.example.object_id

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

resource "azurerm_key_vault_managed_storage_account_sas_token_definition" "example" {
  name                       = "examplesasdefinition"
  validity_period            = "P1D"
  managed_storage_account_id = azurerm_key_vault_managed_storage_account.example.id
  sas_template_uri           = data.azurerm_storage_account_sas.example.sas
  sas_type                   = "account"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAS Definition.

* `managed_storage_account_id` - (Required) The ID of the Managed Storage Account.

* `sas_template_uri` - (Required) The SAS definition token template signed with an arbitrary key. Tokens created according to the SAS definition will have the same properties as the template, but regenerated with a new validity period.

* `sas_type` - (Required) The type of SAS token the SAS definition will create. Possible values are `account` and `service`.

* `validity_period` - (Required) Validity period of SAS token. Value needs to be in [ISO 8601 duration format](https://en.wikipedia.org/wiki/ISO_8601#Durations).

---

* `tags` - (Optional) A mapping of tags which should be assigned to the SAS Definition. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Storage Account SAS Definition.

* `secret_id` - The ID of the Secret that is created by Managed Storage Account SAS Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault.

## Import

Key Vaults can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_storage_account_sas_token_definition.example https://example-keyvault.vault.azure.net/storage/exampleStorageAcc01/sas/exampleSasDefinition01
```
