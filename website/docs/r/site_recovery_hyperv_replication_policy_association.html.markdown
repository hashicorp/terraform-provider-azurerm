---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_hyperv_replication_policy_association"
description: |-
    Manages an Azure Site Recovery replication policy association for HyperV on Azure.
---

# azurerm_site_recovery_hyperv_replication_policy_association

Manages an Azure Site Recovery replication policy for HyperV within a Recovery Vault. 

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "East US"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_site_recovery_services_vault_hyperv_site" "example" {
  recovery_vault_id = azurerm_recovery_services_vault.example.id
  name              = "example-site"
}

resource "azurerm_site_recovery_hyperv_replication_policy" "example" {
  name                                               = "policy"
  recovery_vault_id                                  = azurerm_recovery_services_vault.example.id
  recovery_point_retention_in_hours                  = 2
  application_consistent_snapshot_frequency_in_hours = 1
  replication_interval_in_seconds                    = 300
}

resource "azurerm_site_recovery_hyperv_replication_policy_association" "example" {
  name           = "example-association"
  hyperv_site_id = azurerm_site_recovery_services_vault_hyperv_site.example.id
  policy_id      = azurerm_site_recovery_hyperv_replication_policy.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the replication policy association. Changing this forces a new association to be created.

* `hyperv_site_id` - (Required) The ID of the HyperV site to which the policy should be associated. Changing this forces a new association to be created.

* `policy_id` - (Required) The ID of the HyperV replication policy which to be associated. Changing this forces a new association to be created.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Replication Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Site Recovery Replication Policy Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Replication Policy Association.
* `delete` - (Defaults to 1 hour) Used when deleting the Site Recovery Replication Policy Association.

## Import

Site Recovery Replication Policies can be imported using the `resource id`, e.g.

```shell
terraform import  azurerm_site_recovery_hyperv_replication_policy_association.mypolicy /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/site-name/replicationProtectionContainers/container-name/replicationProtectionContainerMappings/mapping-name
```
