---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent"
description: |-
  Manages a System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.
---

# azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent

Manages a System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.

~> **Note:** By request of the service team the provider no longer automatically registering the `Microsoft.ScVmm` Resource Provider for this resource. To register it you can run `az provider register --namespace Microsoft.ScVmm`.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent" "test" {
  scope               = ""
  https_proxy         = ""
  provisioning_action = "install"

  credential {
    username = "testUser"
    password = "H@sh@123!"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `scope` - (Required) The ID of the Hybrid Machine. Changing this forces a new resource to be created.

* `credential` - (Optional) A `credential` block as defined below.

* `https_proxy` - (Optional) The https proxy url. Changing this forces a new resource to be created.

* `provisioning_action` - (Optional) The provisioning action that is used to define the different types of operations for the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent. Possible values are `install`, `repair` and `uninstall`. Changing this forces a new resource to be created.

---

A `credential` block supports the following:

* `username` - (Optional) The username that is used to connect to the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent. Changing this forces a new resource to be created.

* `password` - (Optional) The password that is used to connect to the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.
* `read` - (Defaults to 5 minutes) Used when retrieving this System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.
* `update` - (Defaults to 30 minutes) Used when updating this System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.
* `delete` - (Defaults to 30 minutes) Used when deleting this System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.

## Import

System Center Virtual Machine Manager Virtual Machine Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.ScVmm/virtualMachineInstances/default/guestAgents/default
```
