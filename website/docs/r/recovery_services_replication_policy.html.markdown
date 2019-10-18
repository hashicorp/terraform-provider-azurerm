---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_recovery_services_replication_policy"
sidebar_current: "docs-azurerm-recovery-services-replication-policy"
description: |-
    Manages a site recovery services replication policy on Azure.
---

# azurerm_recovery_services_replication_policy

Manages a Azure recovery vault replication policy.

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

resource "azurerm_recovery_services_replication_policy" "policy" {
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

* `recovery_point_retention_in_minutes` - (Required) Retain the recovery points for given time in minutes.

* `application_consistent_snapshot_frequency_in_minutes` - (Required) Specifies the frequency(in minutes) at which to create application consistent recovery points.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Site recovery recovery vault fabric can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_recovery_services_protection_container.mycontainer /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationPolicies/policy-name
```
