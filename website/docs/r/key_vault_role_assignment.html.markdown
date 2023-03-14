---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_role_assignment"
description: |-
  Manages a KeyVault Role Assignment.
---

# azurerm_key_vault_role_assignment

Manages a KeyVault Role Assignment.

## Example Usage

```hcl
data "azurerm_key_vault_role_definition" "user" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
  scope          = "/"
}

resource "azurerm_key_vault_role_assignment" "example" {
  name               = "a9dbe818-56e7-5878-c0ce-a1477692c1d6"
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.example.hsm_uri
  scope              = "${data.azurerm_key_vault_role_definition.user.scope}"
  role_definition_id = "${data.azurerm_key_vault_role_definition.user.resource_id}"
  principal_id       = "${data.azurerm_client_config.current.object_id}"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name in GUID notation which should be used for this KeyVault Role Assignment. Changing this forces a new KeyVault to be created.

* `principal_id` - (Required) The principal ID to be assigned to this role. It can point to a user, service principal, or security group. Changing this forces a new KeyVault to be created.

* `role_definition_id` - (Required) The resource ID of the role definition to assign. Changing this forces a new KeyVault to be created.

* `scope` - (Required) Specifies the scope to create the role assignment. Changing this forces a new KeyVault to be created.

* `vault_base_url` - (Required) The HSM URI of a Managed Hardware Security Module resource. Changing this forces a new KeyVault to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the KeyVault Role Assignment with Vault Base URL.

* `resource_id` - The resource id of created assignment resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the KeyVault.
* `read` - (Defaults to 5 minutes) Used when retrieving the KeyVault.
* `delete` - (Defaults to 10 minutes) Used when deleting the KeyVault.

## Import

KeyVaults can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_role_assignment.example https://0000.managedhsm.azure.net///RoleAssignment/00000000-0000-0000-0000-000000000000
```
