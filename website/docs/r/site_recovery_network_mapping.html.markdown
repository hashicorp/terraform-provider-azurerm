---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_network_mapping"
sidebar_current: "docs-azurerm-site-recovery-network-mapping"
description: |-
    Manages a site recovery network mapping on Azure.
---

# azurerm_site_recovery_network_mapping

Manages a site recovery network mapping on Azure. A network mapping decides how to translate connected netwroks when a VM is migrated from one region to another.

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

resource "azurerm_site_recovery_fabric" "primary" {
  name                = "primary-fabric"
  resource_group_name = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.vault.name}"
  location            = "${azurerm_resource_group.primary.location}"
}

resource "azurerm_site_recovery_fabric" "secondary" {
  name                = "secondary-fabric"
  resource_group_name = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.vault.name}"
  location            = "${azurerm_resource_group.secondary.location}"
  depends_on          = ["azurerm_site_recovery_fabric.primary"] # Avoids issues with crearing fabrics simultainusly
}

resource "azurerm_virtual_network" "primary" {
  name                = "network1"
  resource_group_name = "${azurerm_resource_group.primary.name}"
  address_space       = ["192.168.1.0/24"]
  location            = "${azurerm_resource_group.primary.location}"
}

resource "azurerm_virtual_network" "secondary" {
  name                = "network2"
  resource_group_name = "${azurerm_resource_group.secondary.name}"
  address_space       = ["192.168.2.0/24"]
  location            = "${azurerm_resource_group.secondary.location}"
}

resource "azurerm_site_recovery_network_mapping" "recovery-mapping" {
  name                        = "recovery-network-mapping-1"
  resource_group_name         = "${azurerm_resource_group.secondary.name}"
  recovery_vault_name         = "${azurerm_recovery_services_vault.vault.name}"
  source_recovery_fabric_name = "primary-fabric"
  target_recovery_fabric_name = "secondary-fabric"
  source_network_id           = "${azurerm_virtual_network.primary.id}"
  target_network_id           = "${azurerm_virtual_network.secondary.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network mapping.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located.

* `recovery_vault_name` - (Required) The name of the vault that should be updated.

* `source_recovery_fabric_name` - (Required) Specifies the ASR fabric where mapping should be created.

* `target_recovery_fabric_name` - (Required) The Azure Site Recovery fabric object corresponding to the recovery Azure region.

* `source_network_id` - (Required) The id of the primary network.

* `target_network_id` - (Required) The id of the recovery network.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Site recovery network mapping can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_network_mapping.mymapping /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/primary-fabric-name/replicationNetworks/azureNetwork/replicationNetworkMappings/mapping-name
```
