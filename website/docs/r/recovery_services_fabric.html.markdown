---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_recovery_services_fabric"
sidebar_current: "docs-azurerm-recovery-services-fabric"
description: |-
    Manages a site recovery services fabric on Azure.
---

# azurerm_recovery_services_fabric

Manages a Azure recovery vault fabric.

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

resource "azurerm_recovery_services_fabric" "fabric" {
  name                         = "primary-fabric"
  resource_group_name          = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name          = "${azurerm_recovery_services_vault.vault.name}"
  location                     = "${azurerm_resource_group.primary.location}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network mapping.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located.

* `recovery_vault_name` - (Required) The name of the vault that should be updated.

* `location` - (Required) In what region should the fabric be located.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Site recovery recovery vault fabric can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_recovery_services_fabric.myfabric /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/fabric-name
```
