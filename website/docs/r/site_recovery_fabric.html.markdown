---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_fabric"
description: |-
    Manages a Site Recovery Replication Fabric on Azure.
---

# azurerm_site_recovery_fabric

Manages a Azure Site Recovery Replication Fabric within a Recovery Services vault. Only Azure fabrics are supported at this time. Replication Fabrics serve as a container within an Azure region for other Site Recovery resources such as protection containers, protected items, network mappings.

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

resource "azurerm_site_recovery_fabric" "fabric" {
  name                = "primary-fabric"
  resource_group_name = azurerm_resource_group.secondary.name
  recovery_vault_name = azurerm_recovery_services_vault.vault.name
  location            = azurerm_resource_group.primary.location
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network mapping. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located. Changing this forces a new resource to be created.

* `recovery_vault_name` - (Required) The name of the vault that should be updated. Changing this forces a new resource to be created.

* `location` - (Required) In what region should the fabric be located. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Fabric.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Site Recovery Fabric.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Fabric.
* `delete` - (Defaults to 30 minutes) Used when deleting the Site Recovery Fabric.

## Import

Site Recovery Fabric can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_fabric.myfabric /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/fabric-name
```
