---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_role_definition"
description: |-
  Gets information about an existing KeyVault Role Definition.
---

# Data Source: azurerm_key_vault_role_definition

Use this data source to access information about an existing KeyVault Role Definition.

## Example Usage

```hcl
data "azurerm_key_vault_role_definition" "example" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
  scope          = "/"
}

output "id" {
  value = data.azurerm_key_vault_role_definition.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name in UUID notation of this KeyVault Role Definition.

* `vault_base_url` - (Required) Specify the base URL of the Managed HSM resource.

* `scope` - (Optional) Specify the scope to retrieve the role definition. Defaults to `/`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the KeyVault Role Definition.

* `assignable_scopes` - A list of assignable role scope. Possible values are `/` and `/keys`.

* `description` - A text description of this role definition.

* `permission` - A `permission` block as defined below.

* `resource_id` - The ID of the role definition resource without base url.

* `role_name` - The role name of the role definition.

* `role_type` - The type of the role definition. Possible values are `AKVBuiltInRole` and `CustomRole`.

---

A `permission` block exports the following:

* `actions` - (Optional) A list of action permission granted.

* `not_actions` - (Optional) A list of action permission excluded (but not denied).

* `data_actions` - (Optional) A list of data action permission granted. Possible values are `Microsoft.KeyVault/managedHsm/keys/read/action`, `Microsoft.KeyVault/managedHsm/keys/write/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/read/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/recover/action`, `Microsoft.KeyVault/managedHsm/keys/backup/action`, `Microsoft.KeyVault/managedHsm/keys/restore/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/delete/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/read/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/write/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/read/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/write/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/delete/action`, `Microsoft.KeyVault/managedHsm/keys/encrypt/action`, `Microsoft.KeyVault/managedHsm/keys/decrypt/action`, `Microsoft.KeyVault/managedHsm/keys/wrap/action`, `Microsoft.KeyVault/managedHsm/keys/unwrap/action`, `Microsoft.KeyVault/managedHsm/keys/sign/action`, `Microsoft.KeyVault/managedHsm/keys/verify/action`, `Microsoft.KeyVault/managedHsm/keys/create`, `Microsoft.KeyVault/managedHsm/keys/delete`, `Microsoft.KeyVault/managedHsm/keys/export/action`, `Microsoft.KeyVault/managedHsm/keys/release/action`, `Microsoft.KeyVault/managedHsm/keys/import/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/delete`, `Microsoft.KeyVault/managedHsm/securitydomain/download/action`, `Microsoft.KeyVault/managedHsm/securitydomain/download/read`, `Microsoft.KeyVault/managedHsm/securitydomain/upload/action`, `Microsoft.KeyVault/managedHsm/securitydomain/upload/read`, `Microsoft.KeyVault/managedHsm/securitydomain/transferkey/read`, `Microsoft.KeyVault/managedHsm/backup/start/action`, `Microsoft.KeyVault/managedHsm/restore/start/action`, `Microsoft.KeyVault/managedHsm/backup/status/action`, `Microsoft.KeyVault/managedHsm/restore/status/action` and `Microsoft.KeyVault/managedHsm/rng/action`.

* `not_data_actions` - (Optional) A list of data action permission granted. Possible values are `Microsoft.KeyVault/managedHsm/keys/read/action`, `Microsoft.KeyVault/managedHsm/keys/write/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/read/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/recover/action`, `Microsoft.KeyVault/managedHsm/keys/backup/action`, `Microsoft.KeyVault/managedHsm/keys/restore/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/delete/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/read/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/write/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/read/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/write/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/delete/action`, `Microsoft.KeyVault/managedHsm/keys/encrypt/action`, `Microsoft.KeyVault/managedHsm/keys/decrypt/action`, `Microsoft.KeyVault/managedHsm/keys/wrap/action`, `Microsoft.KeyVault/managedHsm/keys/unwrap/action`, `Microsoft.KeyVault/managedHsm/keys/sign/action`, `Microsoft.KeyVault/managedHsm/keys/verify/action`, `Microsoft.KeyVault/managedHsm/keys/create`, `Microsoft.KeyVault/managedHsm/keys/delete`, `Microsoft.KeyVault/managedHsm/keys/export/action`, `Microsoft.KeyVault/managedHsm/keys/release/action`, `Microsoft.KeyVault/managedHsm/keys/import/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/delete`, `Microsoft.KeyVault/managedHsm/securitydomain/download/action`, `Microsoft.KeyVault/managedHsm/securitydomain/download/read`, `Microsoft.KeyVault/managedHsm/securitydomain/upload/action`, `Microsoft.KeyVault/managedHsm/securitydomain/upload/read`, `Microsoft.KeyVault/managedHsm/securitydomain/transferkey/read`, `Microsoft.KeyVault/managedHsm/backup/start/action`, `Microsoft.KeyVault/managedHsm/restore/start/action`, `Microsoft.KeyVault/managedHsm/backup/status/action`, `Microsoft.KeyVault/managedHsm/restore/status/action` and `Microsoft.KeyVault/managedHsm/rng/action`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the KeyVault.
