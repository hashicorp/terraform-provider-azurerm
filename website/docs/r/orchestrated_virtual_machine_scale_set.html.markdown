---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_orchestrated_virtual_machine_scale_set"
description: |-
  Manages an Orchestrated Virtual Machine Scale Set.
---

# azurerm_orchestrated_virtual_machine_scale_set

Manages an Orchestrated Virtual Machine Scale Set.

-> **Note:** Orchestrated Virtual Machine Scale Sets are in Public Preview - [more details can be found in the Azure Documentation](https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/orchestration-modes).

-> **Note:** Azure has deprecated the `single_placement_group` attribute in the Orchestrated Virtual Machine Scale Set since api-version 2019-12-01 and there is a breaking change in the Orchestrated Virtual Machine Scale Set. If you have an Orchestrated Virtual Machine Scale Set created using `azurerm` provider version `<=2.13.0` you will have to remove the `single_placement_group` attribute in your config and recreate the resource to have it managed by terraform.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "example" {
  name                = "example-VMSS"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  platform_fault_domain_count = 1

  zones = ["1"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Orchestrated Virtual Machine Scale Set. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Orchestrated Virtual Machine Scale Set should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Orchestrated Virtual Machine Scale Set should exist. Changing this forces a new resource to be created.

* `platform_fault_domain_count` - (Required) Specifies the number of fault domains that are used by this Orchestrated Virtual Machine Scale Set. Changing this forces a new resource to be created.

~> **NOTE:** The number of Fault Domains varies depending on which Azure Region you're using - a list can be found [here](https://github.com/MicrosoftDocs/azure-docs/blob/master/includes/managed-disks-common-fault-domain-region-list.md).

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group which the Virtual Machine should be assigned to. Changing this forces a new resource to be created.

* `single_placement_group` - (Optional / **Deprecated**) Should the Orchestrated Virtual Machine Scale Set use single placement group?

~> **NOTE:** Due to a limitation of the Azure API at this time, you can only assign `single_placement_group` to `false`.

* `zones` - (Optional) A list of Availability Zones in which the Virtual Machines in this Scale Set should be created in. Changing this forces a new resource to be created.

~> **Note:** Due to a limitation of the Azure API at this time only one Availability Zone can be defined.

* `tags` - (Optional) A mapping of tags which should be assigned to this Orchestrated Virtual Machine Scale Set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Orchestrated Virtual Machine Scale Set.

* `unique_id` - The Unique ID for the Orchestrated Virtual Machine Scale Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Orchestrated Virtual Machine Scale Set.
* `update` - (Defaults to 30 minutes) Used when updating the Orchestrated Virtual Machine Scale Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Orchestrated Virtual Machine Scale Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the Orchestrated Virtual Machine Scale Set.

## Import

An Orchestrated Virtual Machine Scale Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_orchestrated_virtual_machine_scale_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/Microsoft.Compute/virtualMachineScaleSets/scaleset1
```
