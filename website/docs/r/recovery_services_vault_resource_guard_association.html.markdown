---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_recovery_services_vault_resource_guard_association"
description: |-
  Manages an association of a Resource Guard and Recovery Services Vault. 
---

# azurerm_recovery_services_vault_resource_guard_association

Manages an association of a Resource Guard and Recovery Services Vault. 

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_protection_resource_guard" "example" {
  name                = "example-resourceguard"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"

  soft_delete_enabled = true
}

resource "azurerm_recovery_services_vault_resource_guard_association" "test" {
  vault_id          = azurerm_recovery_services_vault.test.id
  resource_guard_id = azurerm_data_protection_resource_guard.test.id
}
```

## Arguments Reference

The following arguments are supported:

* `vault_id` - (Required) ID of the Recovery Services Vault which should be associated with. Changing this forces a new resource to be created.

* `resource_guard_id` - (Required) ID of the Resource Guard which should be associated with. Changing this forces a new resource to be created. 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Guard.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Guard.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Guard.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Guard.

## Import

Resource Guards can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_recovery_services_vault_resource_guard_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/backupResourceGuardProxies/proxy1
```
