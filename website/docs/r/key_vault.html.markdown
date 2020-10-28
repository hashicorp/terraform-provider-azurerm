---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault"
description: |-
  Manages a Key Vault.
---

# azurerm_key_vault

Manages a Key Vault.

## Disclaimers

~> **Note:** It's possible to define Key Vault Access Policies both within [the `azurerm_key_vault` resource](key_vault.html) via the `access_policy` block and by using [the `azurerm_key_vault_access_policy` resource](key_vault_access_policy.html). However it's not possible to use both methods to manage Access Policies within a KeyVault, since there'll be conflicts.

~> **Note:** Terraform will automatically recover a soft-deleted Key Vault during Creation if one is found - you can opt out of this using the `features` block within the Provider block.

## Example Usage

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_key_vault" "example" {
  name                        = "testvault"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_enabled         = true
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "get",
      "ManageContacts",
    ]

    secret_permissions = [
      "get",
    ]

    storage_permissions = [
      "get",
    ]
  }

  network_acls {
    default_action = "Deny"
    bypass         = "AzureServices"
  }

  contact {
    email = "example@example.com"
    name  = "example"
    phone = "0123456789"
  }

  tags = {
    environment = "Testing"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Key Vault. Changing this forces a new resource to be created.

* `sku_name` - (Required) The Name of the SKU used for this Key Vault. Possible values are `standard` and `premium`.

* `tenant_id` - (Required) The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault.

---

* `access_policy` - (Optional) [A list](/docs/configuration/attr-as-blocks.html) of up to 16 objects describing access policies, as described below.

-> **NOTE** Since `access_policy` can be configured both inline and via the separate `azurerm_key_vault_access_policy` resource, we have to explicitly set it to empty slice (`[]`) to remove it.

* `enabled_for_deployment` - (Optional) Boolean flag to specify whether Azure Virtual Machines are permitted to retrieve certificates stored as secrets from the key vault. Defaults to `false`.

* `enabled_for_disk_encryption` - (Optional) Boolean flag to specify whether Azure Disk Encryption is permitted to retrieve secrets from the vault and unwrap keys. Defaults to `false`.

* `enabled_for_template_deployment` - (Optional) Boolean flag to specify whether Azure Resource Manager is permitted to retrieve secrets from the key vault. Defaults to `false`.

* `enable_rbac_authorization` - (Optional) Boolean flag to specify whether Azure Key Vault uses Role Based Access Control (RBAC) for authorization of data actions. Defaults to `false`.

* `network_acls` - (Optional) A `network_acls` block as defined below.

* `purge_protection_enabled` - (Optional) Is Purge Protection enabled for this Key Vault? Defaults to `false`.

!> **Note:** Once Purge Protection has been Enabled it's not possible to Disable it. Support for [disabling purge protection is being tracked in this Azure API issue](https://github.com/Azure/azure-rest-api-specs/issues/8075). Deleting the Key Vault with Purge Protection Enabled will schedule the Key Vault to be deleted (which will happen by Azure in the configured number of days, currently 90 days - which will be configurable in Terraform in the future).

* `soft_delete_enabled` - (Optional) Should Soft Delete be enabled for this Key Vault? Defaults to `false`.

!> **Note:** Once Soft Delete has been Enabled it's not possible to Disable it.

~> **Note:** Terraform will check when creating a Key Vault for a previous soft-deleted Key Vault and recover it if one exists. You can configure this behaviour using the `features` block within the `provider` block.  

* `soft_delete_retention_days` - (Optional) The number of days that items should be retained for once soft-deleted.

~> **Note:** This field can only be set once Soft Delete is Enabled and cannot be updated.

* `contact` - (Optional) One or more `contact` block as defined below.

~> **Note:** This field can only be set once user has `ManageContacts` permission.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `access_policy` block supports the following:

* `tenant_id` - (Required) The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault. Must match the `tenant_id` used above.

* `object_id` - (Required) The object ID of a user, service principal or security group in the Azure Active Directory tenant for the vault. The object ID must be unique for the list of access policies.

* `application_id` - (Optional) The object ID of an Application in Azure Active Directory.

* `certificate_permissions` - (Optional) List of certificate permissions, must be one or more from the following: `backup`, `create`, `delete`, `deleteissuers`, `get`, `getissuers`, `import`, `list`, `listissuers`, `managecontacts`, `manageissuers`, `purge`, `recover`, `restore`, `setissuers` and `update`.

* `key_permissions` - (Optional) List of key permissions, must be one or more from the following: `backup`, `create`, `decrypt`, `delete`, `encrypt`, `get`, `import`, `list`, `purge`, `recover`, `restore`, `sign`, `unwrapKey`, `update`, `verify` and `wrapKey`.

* `secret_permissions` - (Optional) List of secret permissions, must be one or more from the following: `backup`, `delete`, `get`, `list`, `purge`, `recover`, `restore` and `set`.

* `storage_permissions` - (Optional) List of storage permissions, must be one or more from the following: `backup`, `delete`, `deletesas`, `get`, `getsas`, `list`, `listsas`, `purge`, `recover`, `regeneratekey`, `restore`, `set`, `setsas` and `update`.

---

A `network_acls` block supports the following:

* `bypass` - (Required) Specifies which traffic can bypass the network rules. Possible values are `AzureServices` and `None`.

* `default_action` - (Required) The Default Action to use when no rules match from `ip_rules` / `virtual_network_subnet_ids`. Possible values are `Allow` and `Deny`.

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the Key Vault.

* `virtual_network_subnet_ids` - (Optional) One or more Subnet ID's which should be able to access this Key Vault.

---

A `contact` block supports the following:

* `email` - (Required) E-mail address of the contact.

* `name` - (Optional) Name of the contact.

* `phone` - (Optional) Phone number of the contact.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Key Vault.

* `vault_uri` - The URI of the Key Vault, used for performing operations on keys and secrets.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault.

## Import

Key Vault's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/vault1
```
