---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_scale_set_standby_pool"
description: |-
  Manages a Standby Pool for Virtual Machine Scale Sets.
---

# azurerm_virtual_machine_scale_set_standby_pool

Manages a Standby Pool for Virtual Machine Scale Sets.
~> **Note:** please follow the prerequisites mentioned in this [article](https://learn.microsoft.com/azure/virtual-machine-scale-sets/standby-pools-create?tabs=portal#prerequisites) before using this resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "example" {
  name                        = "example-ovmss"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  platform_fault_domain_count = 1
  zones                       = ["1"]
}

resource "azurerm_virtual_machine_scale_set_standby_pool" "example" {
  name                                  = "example-spsvmp"
  resource_group_name                   = azurerm_resource_group.example.name
  location                              = "West Europe"
  attached_virtual_machine_scale_set_id = azurerm_orchestrated_virtual_machine_scale_set.example.id
  virtual_machine_state                 = "Running"

  elasticity_profile {
    max_ready_capacity = 10
    min_ready_capacity = 5
  }

  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Standby Pool. Changing this forces a new Standby Pool to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Standby Pool should exist. Changing this forces a new Standby Pool to be created.

* `location` - (Required) Specifies the Azure Region where the Standby Pool should exist. Changing this forces a new Standby Pool to be created.

* `attached_virtual_machine_scale_set_id` - (Required) Specifies the fully qualified resource ID of a virtual machine scale set the pool is attached to.

* `elasticity_profile` - (Required) An `elasticity_profile` block as defined below.

* `virtual_machine_state` - (Required) Specifies the desired state of virtual machines in the pool. Possible values are `Running` and `Deallocated`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Standby Pool.

---

An `elasticity_profile` block supports the following:

* `max_ready_capacity` - (Required) Specifies the maximum number of virtual machines in the standby pool.

* `min_ready_capacity` - (Required) Specifies the desired minimum number of virtual machines in the standby pool.

~> **Note:** `min_ready_capacity` cannot exceed `max_ready_capacity`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Standby Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Standby Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Standby Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Standby Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Standby Pool.

## Import

Standby Pool can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_scale_set_standby_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StandbyPool/standbyVirtualMachinePools/standbyVirtualMachinePool1
```
