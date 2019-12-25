---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_encryption_set"
sidebar_current: "docs-azurerm-resource-disk-encryption-set"
description: |-
  Manage Azure DiskEncryptionSet instance.
---

# azurerm_disk_encryption_set

Manage an Azure DiskEncryptionSet instance.

-> **NOTE:** The DiskEncryptionSet service is currently in public preview and are only available in a limited set of regions. 

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "des-example-rg"
  location = "westus2"
}

resource "azurerm_key_vault" "test" {
  name                = "des-example-keyvault"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name                = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.service_principal_object_id

    key_permissions = [
      "create",
      "get",
      "delete",
      "list",
      "wrapkey",
      "unwrapkey",
      "get",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "des-example-key"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "des"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  active_key {
    source_vault_id = azurerm_key_vault.test.id
    key_url         = azurerm_key_vault_key.test.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the disk encryption set that is being created. The name can't be changed after the disk encryption set is created. Supported characters for the name are a-z, A-Z, 0-9, _ and -. The maximum name length is 80 characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Disk Encryption Set should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Disk Encryption Set exists. Changing this forces a new resource to be created.

* `active_key` - (Required) One `active_key` block defined below.

* `identity` - (Optional) A `identity` block defined below.

* `tags` - (Optional) A mapping of tags to assign to the Disk Encryption Set.

---

The `active_key` block supports the following:

* `source_vault_id` - (Required) Specifies the resource ID of the KeyVault containing the key or secret.

* `key_url` - (Required) Specifies the URL pointing to a key or secret in KeyVault.

-> **NOTE** Access to the KeyVault must be granted for this Disk Encryption Set, if you want to further use this Disk Encryption Set in a Managed Disk or Virtual Machine, or Virtual Machine Scale Set. For instructions, please refer to the doc of [Server side encryption of Azure managed disks](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/disk-encryption).

## Attributes Reference

The following attributes are exported:

* `previous_keys` - One or more `previous_key` block defined below.

* `id` - The ID of the Disk Encryption Set.

---

A `identity` block supports the following:

* `type` - (Required) The Managed Service Identity Type of this Disk Encryption Set. The possible value is `SystemAssigned` (where Azure will generate a Service Principal for you).

~> **NOTE:** When `type` is set to `SystemAssigned`, identity the Principal ID can be retrieved after the Disk Encryption Set has been created. See [documentation](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview) for additional information.

---

The `previous_key` block contains the following:

* `source_vault_id` - (Required) Specifies the resource ID of the KeyVault containing the key or secret.

* `key_url` - (Required) Specifies the URL pointing to a key or secret in KeyVault.

## Import

Disk Encryption Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_disk_encryption_set.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/diskEncryptionSets/des1
```