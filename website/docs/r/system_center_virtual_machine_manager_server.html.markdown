---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_server"
description: |-
  Manages a System Center Virtual Machine Manager Server.
---

# azurerm_system_center_virtual_machine_manager_server

Manages a System Center Virtual Machine Manager Server.

~> **Note:** By request of the service team the provider no longer automatically registering the `Microsoft.ScVmm` Resource Provider for this resource. To register it you can run `az provider register --namespace Microsoft.ScVmm`.

-> **Note:** This resource depends on an existing `System Center Virtual Machine Manager Host Machine`, `Arc Resource Bridge` and `Custom Location`. Installing and configuring these dependencies is outside the scope of this document. See [Virtual Machine Manager documentation](https://learn.microsoft.com/en-us/system-center/vmm/?view=sc-vmm-2022) and [Install VMM](https://learn.microsoft.com/en-us/system-center/vmm/install?view=sc-vmm-2022) for more details of `System Center Virtual Machine Manager Host Machine`. See [What is Azure Arc resource bridge](https://learn.microsoft.com/en-us/azure/azure-arc/resource-bridge/overview) and [Overview of Arc-enabled System Center Virtual Machine Manager](https://learn.microsoft.com/en-us/azure/azure-arc/system-center-virtual-machine-manager/overview) for more details of `Arc Resource Bridge/Appliance`. See [Create and manage custom locations on Azure Arc-enabled Kubernetes](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/custom-locations) for more details of `Custom Location`. If you encounter issues while configuring, we'd recommend opening a ticket with Microsoft Support.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the System Center Virtual Machine Manager should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the System Center Virtual Machine Manager Server should exist. Changing this forces a new resource to be created.

* `custom_location_id` - (Required) The ID of the Custom Location for the System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

* `fqdn` - (Required) The FQDN of the System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

* `username` - (Required) The username that is used to connect to the System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

* `password` - (Required) The password that is used to connect to the System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

* `port` - (Optional) The port on which the System Center Virtual Machine Manager Server is listening. Possible values are between `1` and `65535`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the System Center Virtual Machine Manager Server.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the System Center Virtual Machine Manager Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the System Center Virtual Machine Manager Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the System Center Virtual Machine Manager Server.
* `update` - (Defaults to 3 hours) Used when updating the System Center Virtual Machine Manager Server.
* `delete` - (Defaults to 3 hours) Used when deleting the System Center Virtual Machine Manager Server.

## Import

System Center Virtual Machine Manager Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_system_center_virtual_machine_manager_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.ScVmm/vmmServers/vmmServer1
```
