---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module"
description: |-
  Manages a Key Vault Managed Hardware Security Module.
---

# azurerm_key_vault_managed_hardware_security_module

Manages a Key Vault Managed Hardware Security Module.

~> **Note:** The Azure Provider includes a Feature Toggle which will purge a Key Vault Managed Hardware Security Module resource on destroy, rather than the default soft-delete. See [`purge_soft_deleted_hardware_security_modules_on_destroy`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block#purge_soft_deleted_hardware_security_modules_on_destroy) for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_hardware_security_modules_on_destroy = true
    }
  }
}
data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault_managed_hardware_security_module" "example" {
  name                       = "exampleKVHsm"
  resource_group_name        = azurerm_resource_group.example.name
  location                   = azurerm_resource_group.example.location
  sku_name                   = "Standard_B1"
  purge_protection_enabled   = false
  soft_delete_retention_days = 90
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]

  tags = {
    Env = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Managed Hardware Security Module. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Key Vault Managed Hardware Security Module. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `admin_object_ids` - (Required) Specifies a list of administrators object IDs for the key vault Managed Hardware Security Module. Changing this forces a new resource to be created.

* `sku_name` - (Required) The Name of the SKU used for this Key Vault Managed Hardware Security Module. Possible value is `Standard_B1`. Changing this forces a new resource to be created.

* `tenant_id` - (Required) The Azure Active Directory Tenant ID that should be used for authenticating requests to the key vault Managed Hardware Security Module. Changing this forces a new resource to be created.

* `purge_protection_enabled` - (Optional) Is Purge Protection enabled for this Key Vault Managed Hardware Security Module? Changing this forces a new resource to be created.

* `soft_delete_retention_days` - (Optional) The number of days that items should be retained for once soft-deleted. This value can be between `7` and `90` days. Defaults to `90`. Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Whether traffic from public networks is permitted. Defaults to `true`. Changing this forces a new resource to be created.

* `network_acls` - (Optional) A `network_acls` block as defined below.

* `security_domain_key_vault_certificate_ids` - (Optional) A list of KeyVault certificates resource IDs (minimum of three and up to a maximum of 10) to activate this Managed HSM. More information see [activate-your-managed-hsm](https://learn.microsoft.com/azure/key-vault/managed-hsm/quick-create-cli#activate-your-managed-hsm)

* `security_domain_quorum` - (Optional) Specifies the minimum number of shares required to decrypt the security domain for recovery. This is required when `security_domain_key_vault_certificate_ids` is specified. Valid values are between 2 and 10.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `network_acls` block supports the following:

* `bypass` - (Required) Specifies which traffic can bypass the network rules. Possible values are `AzureServices` and `None`.

* `default_action` - (Required) The Default Action to use. Possible values are `Allow` and `Deny`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Key Vault Secret Managed Hardware Security Module ID.

* `hsm_uri` - The URI of the Key Vault Managed Hardware Security Module, used for performing operations on keys.

* `security_domain_encrypted_data` - This attribute can be used for disaster recovery or when creating another Managed HSM that shares the same security domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Key Vault Managed Hardware Security Module.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Managed Hardware Security Module.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Managed Hardware Security Module.
* `delete` - (Defaults to 1 hour) Used when deleting the Key Vault Managed Hardware Security Module.

## Import

Key Vault Managed Hardware Security Module can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_hardware_security_module.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/managedHSMs/hsm1
```
