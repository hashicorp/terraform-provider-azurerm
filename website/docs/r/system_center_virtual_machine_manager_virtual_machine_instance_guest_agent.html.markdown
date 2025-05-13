---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent"
description: |-
  Manages a System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.
---

# azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent

Manages a System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.

~> **Note:** By request of the service team the provider is no longer automatically registering the `Microsoft.ScVmm` Resource Provider for this resource. To register it you can run `az provider register --namespace Microsoft.ScVmm`.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_arc_machine" "example" {
  name                = "example-arcmachine"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  kind                = "SCVMM"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_system_center_virtual_machine_manager_server" "example" {
  name                = "example-scvmmms"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ExtendedLocation/customLocations/customLocation1"
  fqdn                = "example.labtest"
  username            = "testUser"
  password            = "H@Sh1CoR3!"
}

data "azurerm_system_center_virtual_machine_manager_inventory_items" "example" {
  inventory_type                                  = "Cloud"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.example.id
}

resource "azurerm_system_center_virtual_machine_manager_cloud" "example" {
  name                                                           = "example-scvmmc"
  location                                                       = azurerm_resource_group.example.location
  resource_group_name                                            = azurerm_resource_group.example.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.example.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.example.inventory_items[0].id
}

data "azurerm_system_center_virtual_machine_manager_inventory_items" "example2" {
  inventory_type                                  = "VirtualMachineTemplate"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.example.id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_template" "example" {
  name                                                           = "example-scvmmvmt"
  location                                                       = azurerm_resource_group.example.location
  resource_group_name                                            = azurerm_resource_group.example.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.example.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.example2.inventory_items[0].id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance" "example" {
  scoped_resource_id = azurerm_arc_machine.example.id
  custom_location_id = azurerm_system_center_virtual_machine_manager_server.example.custom_location_id

  infrastructure {
    checkpoint_type                                                 = "Standard"
    system_center_virtual_machine_manager_cloud_id                  = azurerm_system_center_virtual_machine_manager_cloud.example.id
    system_center_virtual_machine_manager_template_id               = azurerm_system_center_virtual_machine_manager_virtual_machine_template.example.id
    system_center_virtual_machine_manager_virtual_machine_server_id = azurerm_system_center_virtual_machine_manager_server.example.id
  }

  operating_system {
    admin_password = "AdminPassword123!"
  }

  lifecycle {
    // Service API always provisions a virtual disk with bus type IDE, hardware, network interface per Virtual Machine Template by default
    ignore_changes = [storage_disk, hardware, network_interface, operating_system.0.computer_name]
  }
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent" "example" {
  scoped_resource_id = azurerm_arc_machine.example.id
  username           = "Administrator"
  password           = "AdminPassword123!"

  depends_on = [azurerm_system_center_virtual_machine_manager_virtual_machine_instance.example]
}
```

## Arguments Reference

The following arguments are supported:

* `scoped_resource_id` - (Required) The ID of the Hybrid Compute Machine where this System Center Virtual Machine Manager Virtual Machine Instance Guest Agent is stored. Changing this forces a new resource to be created.

* `username` - (Required) The username that is used to connect to the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent. Changing this forces a new resource to be created.

* `password` - (Required) The password that is used to connect to the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent. Changing this forces a new resource to be created.

* `provisioning_action` - (Optional) The provisioning action that is used to define the different types of operations for the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent. Possible values are `install`, `repair` and `uninstall`. Defaults to `install`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.
* `read` - (Defaults to 5 minutes) Used when retrieving the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.
* `delete` - (Defaults to 30 minutes) Used when deleting the System Center Virtual Machine Manager Virtual Machine Instance Guest Agent.

## Import

System Center Virtual Machine Manager Virtual Machine Instance Guest Agents can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.ScVmm/virtualMachineInstances/default/guestAgents/default
```
