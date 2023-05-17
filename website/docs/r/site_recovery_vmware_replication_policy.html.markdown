---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_vmware_replication_policy"
description: |-
  Manages a VMWare Replication Policy.
---

# azurerm_site_recovery_vmware_replication_policy

Manages a VMWare Replication Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "eastus"
}

resource "azurerm_recovery_services_vault" "example" {
  name                               = "example-vault"
  location                           = azurerm_resource_group.example.location
  resource_group_name                = azurerm_resource_group.example.name
  sku                                = "Standard"
  classic_vmware_replication_enabled = true

  soft_delete_enabled = false
}

resource "azurerm_site_recovery_vmware_replication_policy" "example" {
  name                                                 = "example-policy"
  recovery_vault_id                                    = azurerm_recovery_services_vault.example.id
  recovery_point_retention_in_minutes                  = 1440
  application_consistent_snapshot_frequency_in_minutes = 240
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Classic Replication Policy. Changing this forces a new Replication Policy to be created.

* `recovery_vault_id` - (Required) ID of the Recovery Services Vault. Changing this forces a new Replication Policy to be created.

* `recovery_point_retention_in_minutes` - (Required) Specifies the period up to which the recovery points will be retained. Must between `0` to `21600`.

* `application_consistent_snapshot_frequency_in_minutes` - (Required) Specifies the frequency at which to create application consistent recovery points. Must between `0` to `720`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Classic Replication Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Classic Replication Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Classic Replication Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Classic Replication Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Classic Replication Policy.

## Import

VMWare Replication Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_vmware_replication_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/vault1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationPolicies/policy1
```
