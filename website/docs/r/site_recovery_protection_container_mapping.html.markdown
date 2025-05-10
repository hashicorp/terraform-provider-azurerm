---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_protection_container_mapping"
description: |-
    Manages a Site Recovery protection container mapping on Azure.
---

# azurerm_site_recovery_protection_container_mapping

Manages a Azure recovery vault protection container mapping. A protection container mapping decides how to translate the protection container when a VM is migrated from one region to another.

## Example Usage

```hcl
resource "azurerm_resource_group" "primary" {
  name     = "tfex-network-mapping-primary"
  location = "West US"
}

resource "azurerm_resource_group" "secondary" {
  name     = "tfex-network-mapping-secondary"
  location = "East US"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.secondary.location
  resource_group_name = azurerm_resource_group.secondary.name
  sku                 = "Standard"
}

resource "azurerm_site_recovery_fabric" "primary" {
  name                = "primary-fabric"
  resource_group_name = azurerm_resource_group.secondary.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  location            = azurerm_resource_group.primary.location
}

resource "azurerm_site_recovery_fabric" "secondary" {
  name                = "secondary-fabric"
  resource_group_name = azurerm_resource_group.secondary.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  location            = azurerm_resource_group.secondary.location
}

resource "azurerm_site_recovery_protection_container" "primary" {
  name                 = "primary-protection-container"
  resource_group_name  = azurerm_resource_group.secondary.name
  recovery_vault_name  = azurerm_recovery_services_vault.vault.name
  recovery_fabric_name = azurerm_site_recovery_fabric.primary.name
}

resource "azurerm_site_recovery_protection_container" "secondary" {
  name                 = "secondary-protection-container"
  resource_group_name  = azurerm_resource_group.secondary.name
  recovery_vault_name  = azurerm_recovery_services_vault.vault.name
  recovery_fabric_name = azurerm_site_recovery_fabric.secondary.name
}

resource "azurerm_site_recovery_replication_policy" "policy" {
  name                                                 = "policy"
  resource_group_name                                  = azurerm_resource_group.secondary.name
  recovery_vault_name                                  = azurerm_recovery_services_vault.vault.name
  recovery_point_retention_in_minutes                  = 24 * 60
  application_consistent_snapshot_frequency_in_minutes = 4 * 60
}

resource "azurerm_site_recovery_protection_container_mapping" "container-mapping" {
  name                                      = "container-mapping"
  resource_group_name                       = azurerm_resource_group.secondary.name
  recovery_vault_name                       = azurerm_recovery_services_vault.vault.name
  recovery_fabric_name                      = azurerm_site_recovery_fabric.primary.name
  recovery_source_protection_container_name = azurerm_site_recovery_protection_container.primary.name
  recovery_target_protection_container_id   = azurerm_site_recovery_protection_container.secondary.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.policy.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the protection container mapping. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) The name of the vault that should be updated. Changing this forces a new resource to be created.

* `recovery_fabric_name` - (Required) Name of fabric that should contains the protection container to map. Changing this forces a new resource to be created.

* `recovery_source_protection_container_name` - (Required) Name of the source protection container to map. Changing this forces a new resource to be created.

* `recovery_target_protection_container_id` - (Required) Id of target protection container to map to. Changing this forces a new resource to be created.

* `recovery_replication_policy_id` - (Required) Id of the policy to use for this mapping. Changing this forces a new resource to be created.

* `automatic_update` - (Optional) a `automatic_update` block defined as below.

---

An `automatic_update` block supports the following:

* `enabled` - (Optional) Should the Mobility service installed on Azure virtual machines be automatically updated. Defaults to `false`.

~> **Note:** The setting applies to all Azure VMs protected in the same container. For more details see [this document](https://learn.microsoft.com/en-us/azure/site-recovery/azure-to-azure-autoupdate#enable-automatic-updates)

* `automation_account_id` - (Optional) The automation account ID which holds the automatic update runbook and authenticates to Azure resources.

~> **Note:** `automation_account_id` is required when `enabled` is specified.

* `authentication_type` - (Optional) The authentication type used for automation account. Possible values are `RunAsAccount` and `SystemAssignedIdentity`. Defaults to `SystemAssignedIdentity`.

~> **Note:** `RunAsAccount` of `authentication_type` is deprecated and will retire on September 30, 2023. Details could be found [here](https://learn.microsoft.com/en-us/azure/automation/whats-new#support-for-run-as-accounts).

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Protection Container Mapping.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Site Recovery Protection Container Mapping.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Protection Container Mapping.
* `update` - (Defaults to 30 minutes) Used when updating the Site Recovery Protection Container Mapping.
* `delete` - (Defaults to 30 minutes) Used when deleting the Site Recovery Protection Container Mapping.

## Import

Site Recovery Protection Container Mappings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_protection_container_mapping.mymapping /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/fabric1/replicationProtectionContainers/container1/replicationProtectionContainerMappings/mapping1
```
