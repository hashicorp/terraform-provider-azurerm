---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_hyperv_replication_policy"
description: |-
    Manages an Azure Site Recovery replication policy for HyperV on Azure.
---

# azurerm_site_recovery_hyperv_replication_policy

Manages a Azure Site Recovery replication policy for HyperV within a Recovery Vault. Replication policies define the frequency at which recovery points are created and how long they are stored.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_site_recovery_hyperv_replication_policy" "policy" {
  name                                               = "policy"
  recovery_vault_id                                  = azurerm_recovery_services_vault.vault.id
  recovery_point_retention_in_hours                  = 2
  application_consistent_snapshot_frequency_in_hours = 1
  replication_interval_in_seconds                    = 300
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the replication policy. Changing this forces a new resource to be created.

* `recovery_vault_id` - (Required) The id of the vault that should be updated. Changing this forces a new resource to be created.

* `recovery_point_retention_in_hours` - (Required) The duration in hours for which the recovery points need to be stored.

* `application_consistent_snapshot_frequency_in_hours` - (Required) Specifies the frequency at which to create application consistent recovery points.

* `replication_interval_in_seconds` - (Required) Specifies how frequently data should be synchronized between source and target locations. Possible values are `30` and `300`.

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
terraform import azurerm_site_recovery_hyperv_replication_policy.mypolicy /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationPolicies/policy-name
```
