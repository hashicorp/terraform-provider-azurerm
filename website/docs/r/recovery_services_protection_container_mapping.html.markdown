---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_recovery_services_protection_container_mapping"
sidebar_current: "docs-azurerm-recovery-services-protection-container-mapping"
description: |-
    Manages a site recovery services protection container mappings on Azure.
---

# azurerm_recovery_services_protection_container_mapping

~> **NOTE:** This resource has been deprecated in favour of the `azurerm_site_recovery_protection_container_mapping` resource and will be removed in the next major version of the AzureRM Provider. The new resource shares the same fields as this one, and information on migrating across [can be found in this guide](../guides/migrating-between-renamed-resources.html).

Manages a Azure recovery vault protection container mapping. A network protection container mapping decides how to translate the protection container when a VM is migrated from one region to another.

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
  location            = "${azurerm_resource_group.secondary.location}"
  resource_group_name = "${azurerm_resource_group.secondary.name}"
  sku                 = "Standard"
}

resource "azurerm_recovery_services_fabric" "primary" {
  name                = "primary-fabric"
  resource_group_name = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.vault.name}"
  location            = "${azurerm_resource_group.primary.location}"
}

resource "azurerm_recovery_services_fabric" "secondary" {
  name                = "secondary-fabric"
  resource_group_name = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.vault.name}"
  location            = "${azurerm_resource_group.secondary.location}"
}

resource "azurerm_recovery_services_protection_container" "primary" {
  name                 = "primary-protection-container"
  resource_group_name  = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name  = "${azurerm_recovery_services_vault.vault.name}"
  recovery_fabric_name = "${azurerm_recovery_services_fabric.primary.name}"
}

resource "azurerm_recovery_services_protection_container" "secondary" {
  name                 = "secondary-protection-container"
  resource_group_name  = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name  = "${azurerm_recovery_services_vault.vault.name}"
  recovery_fabric_name = "${azurerm_recovery_services_fabric.secondary.name}"
}

resource "azurerm_recovery_services_replication_policy" "policy" {
  name                                                 = "policy"
  resource_group_name                                  = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name                                  = "${azurerm_recovery_services_vault.vault.name}"
  recovery_point_retention_in_minutes                  = "${24 * 60}"
  application_consistent_snapshot_frequency_in_minutes = "${4 * 60}"
}

resource "azurerm_recovery_services_protection_container_mapping" "container-mapping" {
  name                                      = "container-mapping"
  resource_group_name                       = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name                       = "${azurerm_recovery_services_vault.vault.name}"
  recovery_fabric_name                      = "${azurerm_recovery_services_fabric.primary.name}"
  recovery_source_protection_container_name = "${azurerm_recovery_services_protection_container.primary.name}"
  recovery_target_protection_container_id   = "${azurerm_recovery_services_protection_container.secondary.id}"
  recovery_replication_policy_id            = "${azurerm_recovery_services_replication_policy.policy.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network mapping.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located.

* `recovery_vault_name` - (Required) The name of the vault that should be updated.

* `recovery_fabric_name` - (Required) Name of fabric that should contains the protection container to map.

* `recovery_source_protection_container_name` - (Required) Name of the protection container to map.

* `recovery_target_protection_container_id` - (Required) Id of protection container to map to.

* `recovery_replication_policy_id` - (Required) Id of the policy to use for this mapping.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Site recovery recovery vault fabric can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_recovery_services_protection_container_mapping.mymapping /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/
```
