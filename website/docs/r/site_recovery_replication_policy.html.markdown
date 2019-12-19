---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_replication_policy"
sidebar_current: "docs-azurerm-site-recovery-replication-policy"
description: |-
    Manages an Azure Site Recovery replication policy on Azure.
---

# azurerm_site_recovery_replication_policy

Manages a Azure Site Recovery replication policy within a recovery vault. Replication policies define the frequency at which recovery points are created and how long they are stored.

## Example Usage

```hcl
resource "azurerm_resource_group" "secondary" {
  name     = "tfex-network-mapping-secondary"
  location = "East US"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = "${azurerm_resource_group.secondary.location}"
  resource_group_name = "${azurerm_resource_group.secondary.name}"
  sku                 = "Standard"
}

resource "azurerm_site_recovery_replication_policy" "policy" {
  name                                                 = "policy"
  resource_group_name                                  = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name                                  = "${azurerm_recovery_services_vault.vault.name}"
  recovery_point_retention_in_minutes                  = "${24 * 60}"
  application_consistent_snapshot_frequency_in_minutes = "${4 * 60}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network mapping.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located.

* `recovery_vault_name` - (Required) The name of the vault that should be updated.

* `recovery_point_retention_in_minutes` - (Required) The duration in minutes for which the recovery points need to be stored.

* `application_consistent_snapshot_frequency_in_minutes` - (Required) Specifies the frequency(in minutes) at which to create application consistent recovery points.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Site Recovery replication policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_replication_policy.mypolicy /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationPolicies/policy-name
```
