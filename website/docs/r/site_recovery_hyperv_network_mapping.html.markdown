---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_hyperv_network_mapping"
description: |-
    Manages a HyperV site recovery network mapping on Azure.
---

# azurerm_site_recovery_hyperv_network_mapping

Manages a HyperV site recovery network mapping on Azure. A HyperV network mapping decides how to translate connected networks when a VM is migrated from HyperV VMM Center to Azure.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "target" {
  name     = "tfex-network-mapping"
  location = "East US"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.target.location
  resource_group_name = azurerm_resource_group.target.name
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "target" {
  name                = "network"
  resource_group_name = azurerm_resource_group.target.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_resource_group.target.location
}

resource "azurerm_site_recovery_hyperv_network_mapping" "recovery-mapping" {
  name                                              = "recovery-network-mapping"
  recovery_vault_id                                 = azurerm_recovery_services_vault.vault.id
  source_system_center_virtual_machine_manager_name = "my-vmm-server"
  source_network_name                               = "my-vmm-network"
  target_network_id                                 = azurerm_virtual_network.target.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the HyperV network mapping. Changing this forces a new resource to be created.

* `recovery_vault_id` - (Required) The ID of the Recovery Services Vault where the HyperV network mapping should be created. Changing this forces a new resource to be created.

* `source_system_center_virtual_machine_manager_name` - (Required) Specifies the name of source System Center Virtual Machine Manager where the source network exists. Changing this forces a new resource to be created. 

* `source_network_name` - (Required) The Name of the primary network. Changing this forces a new resource to be created.

* `target_network_id` - (Required) The id of the recovery network. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery HyperV Network Mapping.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Site Recovery HyperV Network Mapping.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery HyperV Network Mapping.
* `delete` - (Defaults to 30 minutes) Used when deleting the Site Recovery HyperV Network Mapping.

## Import

Site Recovery Network Mapping can be imported using the `resource id`, e.g.

```shell
terraform import  azurerm_site_recovery_hyperv_network_mapping.mymapping /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/primary-fabric-name/replicationNetworks/azureNetwork/replicationNetworkMappings/mapping-name
```
