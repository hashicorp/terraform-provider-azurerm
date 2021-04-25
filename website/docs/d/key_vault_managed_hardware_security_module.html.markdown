---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module"
description: |-
  Gets information about an existing Key Vault Managed Hardware Security Module.
---

# Data Source: azurerm_key_vault_managed_hardware_security_module

Use this data source to access information about an existing Key Vault Managed Hardware Security Module.

## Example Usage

```hcl
data "azurerm_key_vault_managed_hardware_security_module" "example" {
  name                = "mykeyvaultHsm"
  resource_group_name = "some-resource-group"
}

output "hsm_uri" {
  value = data.azurerm_key_vault_managed_hardware_security_module.example.hsm_uri
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Key Vault Managed Hardware Security Module.

* `resource_group_name` - The name of the Resource Group in which the Key Vault Managed Hardware Security Module exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Managed Hardware Security Module ID.

* `admin_object_ids` - Specifies a list of administrators object IDs for the key vault Managed Hardware Security Module.
  
* `hsm_uri` - The URI of the Hardware Security Module for performing operations on keys and secrets.

* `location` - The Azure Region in which the Key Vault managed Hardware Security Module exists.

* `purge_protection_enabled` - Is purge protection enabled on this Key Vault Managed Hardware Security Module?

* `sku_name` - The Name of the SKU used for this Key Vault Managed Hardware Security Module.

* `soft_delete_retention_days` - The number of days that items should be retained for soft-deleted.

* `tenant_id` - The Azure Active Directory Tenant ID used for authenticating requests to the Key Vault Managed Hardware Security Module.

* `tags` - A mapping of tags assigned to the Key Vault Managed Hardware Security Module.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Managed Hardware Security Module.
