---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_managed_hardware_security_module_role_definition"
description: |-
  Gets information about an existing Key Vault Managed Hardware Security Module Role Definition.
---

# Data Source: azurerm_key_vault_managed_hardware_security_module_role_definition

Use this data source to access information about an existing Key Vault Managed Hardware Security Module Role Definition.

## Example Usage

```hcl
data "azurerm_key_vault_managed_hardware_security_module_role_definition" "example" {
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.example.id
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
}

output "id" {
  value = data.azurerm_key_vault_managed_hardware_security_module_role_definition.example.resource_manager_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name in UUID notation of this Key Vault Managed Hardware Security Module Role Definition.

* `managed_hsm_id` - (Required) The ID of the Key Vault Managed Hardware Security Module.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Key Vault Managed Hardware Security Module Role Definition.

* `assignable_scopes` - A list of assignable role scopes. Possible values are `/` and `/keys`.

* `description` - A text description of the Key Vault Managed Hardware Security Module Role Definition.

* `permission` - A `permission` block as defined below.

* `resource_manager_id` - The ID of the Key Vault Managed Hardware Security Module Role Definition resource without base url.

* `role_name` - The display name of the Key Vault Managed Hardware Security Module Role Definition.

* `role_type` - The type of the Key Vault Managed Hardware Security Module Role Definition. Possible values are `AKVBuiltInRole` and `CustomRole`.

---

A `permission` block exports the following:

* `actions` - A list of action permission granted.

* `not_actions` - A list of action permission excluded (but not denied).

* `data_actions` - A list of data action permission granted.

* `not_data_actions` - A list of data action permission granted.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Managed Hardware Security Module Role Definition.
