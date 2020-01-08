---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_encryption_set"
sidebar_current: "docs-azurerm-resource-disk-encryption-set"
description: |-
  Manages a Disk Encryption Set.
---

# azurerm_disk_encryption_set

Manages a Disk Encryption Set.

-> **NOTE:** The Disk Encryption Sets are currently in Public Preview and are only available in a limited set of regions: West Central US, Canada Central and North Europe. 

-> **NOTE:** At this time the Key Vault used to store the Active Key for this Disk Encryption Set must have both Soft Delete & Purge Protection enabled - which are not yet supported by Terraform - instead you can configure this using [a provisioner](https://www.terraform.io/docs/provisioners/local-exec.html) or [the `azurerm_template_deployment` resource](https://www.terraform.io/docs/providers/azurerm/r/template_deployment.html).

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                = "des-example-keyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
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

resource "azurerm_key_vault_key" "example" {
  name         = "des-example-key"
  key_vault_id = azurerm_key_vault.example.id
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

resource "azurerm_disk_encryption_set" "example" {
  name                = "des"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  key_vault_key_uri   = azurerm_key_vault_key.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Disk Encryption Set. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Disk Encryption Set should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Disk Encryption Set exists. Changing this forces a new resource to be created.

* `key_vault_key_uri` - (Required) Specifies the URL to a Key Vault Key (either from a Key Vault Key, or the Key URL for the Key Vault Secret).

-> **NOTE** Access to the KeyVault must be granted for this Disk Encryption Set, if you want to further use this Disk Encryption Set in a Managed Disk or Virtual Machine, or Virtual Machine Scale Set. For instructions, please refer to the doc of [Server side encryption of Azure managed disks](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/disk-encryption).

* `identity` - (Optional) A `identity` block defined below.

* `tags` - (Optional) A mapping of tags to assign to the Disk Encryption Set.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Disk Encryption Set.

---

A `identity` block exports the following:

* `type` - (Required) The Managed Service Identity Type of this Disk Encryption Set. The possible value is `SystemAssigned` (where Azure will generate a Service Principal for you).

~> **NOTE:** When `type` is set to `SystemAssigned`, identity the Principal ID can be retrieved after the Disk Encryption Set has been created. See [documentation](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview) for additional information.

## Import

Disk Encryption Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_disk_encryption_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/diskEncryptionSets/encryptionSet1
```