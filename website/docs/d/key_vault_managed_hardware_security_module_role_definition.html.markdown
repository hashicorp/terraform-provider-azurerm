---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_managed_hardware_security_module_role_definition"
description: |-
  Gets information about an existing KeyVault Role Definition.
---

# Data Source: azurerm_key_vault_managed_hardware_security_module_role_definition

Use this data source to access information about an existing KeyVault Role Definition.

## Example Usage

```hcl
data "azurerm_key_vault_managed_hardware_security_module_role_definition" "example" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
}

output "id" {
  value = data.azurerm_key_vault_managed_hardware_security_module_role_definition.example.resource_manager_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name in UUID notation of this KeyVault Role Definition.

* `vault_base_url` - (Required) Specify the base URL of the Managed HSM resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the KeyVault Role Definition.

* `assignable_scopes` - A list of assignable role scope. Possible values are `/` and `/keys`.

* `description` - A text description of this role definition.

* `permission` - A `permission` block as defined below.

* `resource_manager_id` - The ID of the role definition resource without base url.

* `role_name` - The role name of the role definition.

* `role_type` - The type of the role definition. Possible values are `AKVBuiltInRole` and `CustomRole`.

---

A `permission` block exports the following:

* `actions` - A list of action permission granted.

* `not_actions` - A list of action permission excluded (but not denied).

* `data_actions` - A list of data action permission granted.

* `not_data_actions` - (Optional) A list of data action permission granted.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the KeyVault.
