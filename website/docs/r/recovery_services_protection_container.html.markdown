---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_recovery_services_protection_container"
description: |-
    Manages a site recovery services protection container on Azure.
---

# azurerm_recovery_services_protection_container

~> **NOTE:** This resource has been deprecated in favour of the `azurerm_site_recovery_protection_container` resource and will be removed in the next major version of the AzureRM Provider. The new resource shares the same fields as this one, and information on migrating across [can be found in this guide](../guides/migrating-between-renamed-resources.html).

Manages a Azure recovery vault protection container.

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

resource "azurerm_recovery_services_fabric" "fabric" {
  name                = "primary-fabric"
  resource_group_name = azurerm_resource_group.secondary.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  location            = azurerm_resource_group.primary.location
}

resource "azurerm_recovery_services_protection_container" "protection-container" {
  name                 = "protection-container"
  resource_group_name  = azurerm_resource_group.secondary.name
  recovery_vault_name  = azurerm_recovery_services_vault.vault.name
  recovery_fabric_name = azurerm_recovery_services_fabric.fabric.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network mapping.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located.

* `recovery_vault_name` - (Required) The name of the vault that should be updated.

* `recovery_fabric_name` - (Required) Name of fabric that should contain this protection container.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Recovery Services Protection Container.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Recovery Services Protection Container.
* `update` - (Defaults to 30 minutes) Used when updating the Recovery Services Protection Container.
* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Services Protection Container.
* `delete` - (Defaults to 30 minutes) Used when deleting the Recovery Services Protection Container.

## Import

Recovery Services Protection Containers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_recovery_services_protection_container.mycontainer /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/fabric-name/replicationProtectionContainers/protection-container-name
```
