---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_recovery_services_vault_guard_proxy"
description: |-
  Manages a Recovery Services Vault Guard Proxy.
---

# azurerm_recovery_services_vault_guard_proxy

Manages a Recovery Services Vault Guard Proxy.

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

resource "azurerm_recovery_services_vault_guard_proxy" "test" {
  name              = "example-guard-proxy"
  vault_id          = azurerm_recovery_services_vault.test.id
  resource_guard_id = azurerm_data_protection_resource_guard.test.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Recovery Services Vault Guard Proxy. Changing this forces a new resource to be created.

* `vault_id` - (Required) ID of the Recovery Services Vault where the Resource Guard Proxy should exist. Changing this forces a new resource to be created.

* `resource_guard_id` - (Required) ID of the Resource Guard which the Resource Guard Proxy should be associated with. Changing this forces a new resource to be created. 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Guard.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Guard.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Guard.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Guard.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Guard.

## Import

Resource Guards can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_recovery_services_vault_guard_proxy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/backupResourceGuardProxies/proxy1
```
