---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module"
description: |-
  Manages a Key Vault Managed Hardware Security Module.
---

# azurerm_key_vault_managed_hardware_security_module

Manages a Key Vault Managed Hardware Security Module.

## Example Usage

```hcl
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

* `purge_protection_enabled` - (Optional) Is Purge Protection enabled for this Key Vault Managed Hardware Security Module? Defaults to `false`. Changing this forces a new resource to be created.

* `soft_delete_retention_days` - (Optional) The number of days that items should be retained for once soft-deleted. This value can be between `7` and `90` days. Defaults to `90`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Secret Managed Hardware Security Module ID.

* `hsm_uri` - The URI of the Key Vault Managed Hardware Security Module, used for performing operations on keys.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Managed Hardware Security Module.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Managed Hardware Security Module.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Managed Hardware Security Module.

## Import

Key Vault Managed Hardware Security Module can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_hardware_security_module.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/managedHSMs/hsm1
```
