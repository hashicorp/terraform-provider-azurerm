---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module_role_definition"
description: |-
  Manages a KeyVault Managed Hardware Security Module Role Definition.
---

# azurerm_key_vault_managed_hardware_security_module_role_definition

Manages a KeyVault Managed Hardware Security Module Role Definition. This resource works together with [Managed hardware security module resource](./key_vault_managed_hardware_security_module).

## Example Usage

```hcl
resource "azurerm_key_vault_managed_hardware_security_module" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  sku_name                 = "Standard_B1"
  tenant_id                = data.azurerm_client_config.current.tenant_id
  admin_object_ids         = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled = false

  active_config {
    security_domain_certificate = [
      azurerm_key_vault_certificate.cert[0].id,
      azurerm_key_vault_certificate.cert[1].id,
      azurerm_key_vault_certificate.cert[2].id,
    ]
    security_domain_quorum = 2
  }
}

resource "azurerm_key_vault_managed_hardware_security_module_role_definition" "example" {
  name           = "7d206142-bf01-11ed-80bc-00155d61ee9e"
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.example.hsm_uri
  description    = "desc foo"
  permission {
    data_actions = [
      "Microsoft.KeyVault/managedHsm/keys/read/action",
    ]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this KeyVault Role Definition. Changing this forces a new KeyVault Role Definition to be created.

* `vault_base_url` - (Required) The base URL of the managed hardware security module resource. Changing this forces a new KeyVault Role Definition to be created.

---

* `description` - (Optional) Specifies a text description about this KeyVault Role Definition.

* `permission` - (Optional) One or more `permission` blocks as defined below.

* `role_name` - (Optional) Specify a name for this KeyVault Role Definition.

---

A `permission` block supports the following, more details about permission see [permitted-operations](https://learn.microsoft.com/en-us/azure/key-vault/managed-hsm/built-in-roles#permitted-operations):

* `actions` - (Optional) One or more Allowed Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`. See ['Azure Resource Manager resource provider operations'](https://docs.microsoft.com/azure/role-based-access-control/resource-provider-operations) for details.

* `not_actions` - (Optional) One or more Disallowed Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`. See ['Azure Resource Manager resource provider operations'](https://docs.microsoft.com/azure/role-based-access-control/resource-provider-operations) for details.

* `data_actions` - (Optional) Specifies a list of data action permission to grant. Possible values are `Microsoft.KeyVault/managedHsm/keys/read/action`, `Microsoft.KeyVault/managedHsm/keys/write/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/read/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/recover/action`, `Microsoft.KeyVault/managedHsm/keys/backup/action`, `Microsoft.KeyVault/managedHsm/keys/restore/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/delete/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/read/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/write/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/read/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/write/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/delete/action`, `Microsoft.KeyVault/managedHsm/keys/encrypt/action`, `Microsoft.KeyVault/managedHsm/keys/decrypt/action`, `Microsoft.KeyVault/managedHsm/keys/wrap/action`, `Microsoft.KeyVault/managedHsm/keys/unwrap/action`, `Microsoft.KeyVault/managedHsm/keys/sign/action`, `Microsoft.KeyVault/managedHsm/keys/verify/action`, `Microsoft.KeyVault/managedHsm/keys/create`, `Microsoft.KeyVault/managedHsm/keys/delete`, `Microsoft.KeyVault/managedHsm/keys/export/action`, `Microsoft.KeyVault/managedHsm/keys/release/action`, `Microsoft.KeyVault/managedHsm/keys/import/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/delete`, `Microsoft.KeyVault/managedHsm/securitydomain/download/action`, `Microsoft.KeyVault/managedHsm/securitydomain/download/read`, `Microsoft.KeyVault/managedHsm/securitydomain/upload/action`, `Microsoft.KeyVault/managedHsm/securitydomain/upload/read`, `Microsoft.KeyVault/managedHsm/securitydomain/transferkey/read`, `Microsoft.KeyVault/managedHsm/backup/start/action`, `Microsoft.KeyVault/managedHsm/restore/start/action`, `Microsoft.KeyVault/managedHsm/backup/status/action`, `Microsoft.KeyVault/managedHsm/restore/status/action` and `Microsoft.KeyVault/managedHsm/rng/action`.

* `not_data_actions` - (Optional) Specifies a list of data action permission not to grant. Possible values are `Microsoft.KeyVault/managedHsm/keys/read/action`, `Microsoft.KeyVault/managedHsm/keys/write/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/read/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/recover/action`, `Microsoft.KeyVault/managedHsm/keys/backup/action`, `Microsoft.KeyVault/managedHsm/keys/restore/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/delete/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/read/action`, `Microsoft.KeyVault/managedHsm/roleAssignments/write/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/read/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/write/action`, `Microsoft.KeyVault/managedHsm/roleDefinitions/delete/action`, `Microsoft.KeyVault/managedHsm/keys/encrypt/action`, `Microsoft.KeyVault/managedHsm/keys/decrypt/action`, `Microsoft.KeyVault/managedHsm/keys/wrap/action`, `Microsoft.KeyVault/managedHsm/keys/unwrap/action`, `Microsoft.KeyVault/managedHsm/keys/sign/action`, `Microsoft.KeyVault/managedHsm/keys/verify/action`, `Microsoft.KeyVault/managedHsm/keys/create`, `Microsoft.KeyVault/managedHsm/keys/delete`, `Microsoft.KeyVault/managedHsm/keys/export/action`, `Microsoft.KeyVault/managedHsm/keys/release/action`, `Microsoft.KeyVault/managedHsm/keys/import/action`, `Microsoft.KeyVault/managedHsm/keys/deletedKeys/delete`, `Microsoft.KeyVault/managedHsm/securitydomain/download/action`, `Microsoft.KeyVault/managedHsm/securitydomain/download/read`, `Microsoft.KeyVault/managedHsm/securitydomain/upload/action`, `Microsoft.KeyVault/managedHsm/securitydomain/upload/read`, `Microsoft.KeyVault/managedHsm/securitydomain/transferkey/read`, `Microsoft.KeyVault/managedHsm/backup/start/action`, `Microsoft.KeyVault/managedHsm/restore/start/action`, `Microsoft.KeyVault/managedHsm/backup/status/action`, `Microsoft.KeyVault/managedHsm/restore/status/action` and `Microsoft.KeyVault/managedHsm/rng/action`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the KeyVault Role Definition.

* `resource_manager_id` - The ID of the role definition resource without Key Vault base URL.

* `role_type` - The type of the role definition. Possible values are `AKVBuiltInRole` and `CustomRole`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the KeyVault.
* `read` - (Defaults to 5 minutes) Used when retrieving the KeyVault.
* `update` - (Defaults to 10 minutes) Used when updating the KeyVault.
* `delete` - (Defaults to 10 minutes) Used when deleting the KeyVault.

## Import

KeyVaults can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_hardware_security_module_role_definition.example https://0000.managedhsm.azure.net///RoleDefinition/00000000-0000-0000-0000-000000000000
```
