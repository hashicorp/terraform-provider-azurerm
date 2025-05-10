---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_replication_policy"
description: |-
    Manages an Azure Site Recovery replication policy on Azure.
---

# azurerm_site_recovery_replication_policy

Manages a Azure Site Recovery replication policy within a recovery vault. Replication policies define the frequency at which recovery points are created and how long they are stored.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-network-mapping-secondary"
  location = "East US"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_site_recovery_replication_policy" "policy" {
  name                                                 = "policy"
  resource_group_name                                  = azurerm_resource_group.example.name
  recovery_vault_name                                  = azurerm_recovery_services_vault.vault.name
  recovery_point_retention_in_minutes                  = 24 * 60
  application_consistent_snapshot_frequency_in_minutes = 4 * 60
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the replication policy. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) The name of the vault that should be updated. Changing this forces a new resource to be created.

* `recovery_point_retention_in_minutes` - (Required) The duration in minutes for which the recovery points need to be stored.

* `application_consistent_snapshot_frequency_in_minutes` - (Required) Specifies the frequency(in minutes) at which to create application consistent recovery points.

-> **Note:** The value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Replication Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Site Recovery Replication Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Replication Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Site Recovery Replication Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Site Recovery Replication Policy.

## Import

Site Recovery Replication Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_replication_policy.mypolicy /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationPolicies/policy-name
```
