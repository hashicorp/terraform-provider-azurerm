---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module_key"
description: |-
  Gets information about an existing Managed Hardware Security Module Key.

---

# Data Source: azurerm_key_vault_managed_hardware_security_module_key

Use this data source to access information about an existing Managed Hardware Security Module Key.

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_key_vault_managed_hardware_security_module_key" "example" {
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.example.id
  name           = azurerm_key_vault_managed_hardware_security_module_key.example.name
}

output "hsm-key-vesrion" {
  value = data.azurerm_key_vault_managed_hardware_security_module_key.example.version
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Managed Hardware Security Module Key.

* `managed_hsm_id` - Specifies the ID of the Managed Hardware Security Module instance where the Secret resides, available on the `azurerm_key_vault_managed_hardware_security_module_key` Data Source / Resource.

-> **Note:** The Managed Hardware Security Module must be in the same subscription as the provider. If the Managed Hardware Security Module is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

The following attributes are exported:

* `id` - The versionless ID of the Managed Hardware Security Module Key.

* `curve` - The EC Curve name of this Managed Hardware Security Module Key.

* `key_type` - Specifies the Key Type of this Managed Hardware Security Module Key

* `key_size` - Specifies the Size of this Managed Hardware Security Module Key.

* `key_opts` - A list of JSON web key operations assigned to this Managed Hardware Security Module Key

* `tags` - A mapping of tags assigned to this Managed Hardware Security Module Key.

* `version` - The current version of the Managed Hardware Security Module Key.

* `versioned_id` - The versioned ID of the Managed Hardware Security Module Key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Hardware Security Module Key.
