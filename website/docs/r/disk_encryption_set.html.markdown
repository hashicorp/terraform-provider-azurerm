---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_encryption_set"
description: |-
  Manages a Disk Encryption Set.
---

# azurerm_disk_encryption_set

Manages a Disk Encryption Set.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                        = "des-example-keyvault"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "premium"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_key" "example" {
  name         = "des-example-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  depends_on = [
    azurerm_key_vault_access_policy.example-user
  ]

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_disk_encryption_set" "example" {
  name                = "des"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  key_vault_key_id    = azurerm_key_vault_key.example.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "example-disk" {
  key_vault_id = azurerm_key_vault.example.id

  tenant_id = azurerm_disk_encryption_set.example.identity[0].tenant_id
  object_id = azurerm_disk_encryption_set.example.identity[0].principal_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign",
  ]
}

resource "azurerm_key_vault_access_policy" "example-user" {
  key_vault_id = azurerm_key_vault.example.id

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign",
    "GetRotationPolicy",
  ]
}

resource "azurerm_role_assignment" "example-disk" {
  scope                = azurerm_key_vault.example.id
  role_definition_name = "Key Vault Crypto Service Encryption User"
  principal_id         = azurerm_disk_encryption_set.example.identity[0].principal_id
}

```

## Example Usage with Automatic Key Rotation Enabled

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                        = "des-example-keyvault"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "premium"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_key" "example" {
  name         = "des-example-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  depends_on = [
    azurerm_key_vault_access_policy.example-user
  ]

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_disk_encryption_set" "example" {
  name                = "des"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  key_vault_key_id    = azurerm_key_vault_key.example.versionless_id

  auto_key_rotation_enabled = true

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "example-disk" {
  key_vault_id = azurerm_key_vault.example.id

  tenant_id = azurerm_disk_encryption_set.example.identity[0].tenant_id
  object_id = azurerm_disk_encryption_set.example.identity[0].principal_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign",
  ]
}

resource "azurerm_key_vault_access_policy" "example-user" {
  key_vault_id = azurerm_key_vault.example.id

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign",
    "GetRotationPolicy",
  ]
}

resource "azurerm_role_assignment" "example-disk" {
  scope                = azurerm_key_vault.example.id
  role_definition_name = "Key Vault Crypto Service Encryption User"
  principal_id         = azurerm_disk_encryption_set.example.identity[0].principal_id
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Disk Encryption Set. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Disk Encryption Set should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Disk Encryption Set exists. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below.

* `key_vault_key_id` - (Optional) Specifies the URL to a Key Vault Key (either from a Key Vault Key, or the Key URL for the Key Vault Secret). Exactly one of `managed_hsm_key_id`, `key_vault_key_id` must be specified.

-> **Note:** Access to the KeyVault must be granted for this Disk Encryption Set, if you want to further use this Disk Encryption Set in a Managed Disk or Virtual Machine, or Virtual Machine Scale Set. For instructions, please refer to the doc of [Server side encryption of Azure managed disks](https://docs.microsoft.com/azure/virtual-machines/linux/disk-encryption).

-> **Note:** A KeyVault or Managed HSM using [enable_rbac_authorization](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault#enable_rbac_authorization) requires to use `azurerm_role_assignment` to assign the role `Key Vault Crypto Service Encryption User` to this Disk Encryption Set.
In this case, `azurerm_key_vault_access_policy` is not needed.

* `managed_hsm_key_id` - (Optional) Key ID of a key in a managed HSM.  Exactly one of `managed_hsm_key_id`, `key_vault_key_id` must be specified.

* `auto_key_rotation_enabled` - (Optional) Boolean flag to specify whether Azure Disk Encryption Set automatically rotates the encryption Key to latest version or not. Possible values are `true` or `false`. Defaults to `false`.

-> **Note:** When `auto_key_rotation_enabled` is set to `true` the `key_vault_key_id` or `managed_hsm_key_id` must use the `versionless_id`.

-> **Note:** To validate which Key Vault Key version is currently being used by the service it is recommended that you use the `azurerm_disk_encryption_set` data source or run a `terraform refresh` command and check the value of the exported `key_vault_key_url` or `managed_hsm_key_id` field.

-> **Note:** It may take between 10 to 20 minutes for the service to update the Key Vault Key URL once the keys have been rotated.

* `encryption_type` - (Optional) The type of key used to encrypt the data of the disk. Possible values are `EncryptionAtRestWithCustomerKey`, `EncryptionAtRestWithPlatformAndCustomerKeys` and `ConfidentialVmEncryptedWithCustomerKey`. Defaults to `EncryptionAtRestWithCustomerKey`. Changing this forces a new resource to be created.

* `federated_client_id` - (Optional) Multi-tenant application client id to access key vault in a different tenant.

* `tags` - (Optional) A mapping of tags to assign to the Disk Encryption Set.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that is configured on this Disk Encryption Set. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both). 

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Disk Encryption Set.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Disk Encryption Set.

* `key_vault_key_url` - The URL for the Key Vault Key or Key Vault Secret that is currently being used by the service.

---

An `identity` block exports the following:

* `principal_id` - The (Client) ID of the Service Principal.

* `tenant_id` - The ID of the Tenant the Service Principal is assigned in.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Disk Encryption Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Disk Encryption Set.
* `update` - (Defaults to 1 hour) Used when updating the Disk Encryption Set.
* `delete` - (Defaults to 1 hour) Used when deleting the Disk Encryption Set.

## Import

Disk Encryption Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_disk_encryption_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/diskEncryptionSets/encryptionSet1
```
