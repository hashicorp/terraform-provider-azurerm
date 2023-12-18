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

  credential {
    username = "testUser"
    password = "H@Sh1CoR3!"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the System Center Virtual Machine Manager should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the System Center Virtual Machine Manager Server should exist. Changing this forces a new resource to be created.

* `credential` - (Required) A `credential` block as defined below.

* `custom_location_id` - (Required) The ID of the Custom Location for the System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

* `fqdn` - (Required) The FQDN of the System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

* `port` - (Optional) The port that is used to listened by the System Center Virtual Machine Manager Server. Possible values are between `1` and `65535`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the System Center Virtual Machine Manager Server.

---

A `credential` block supports the following:

* `username` - (Required) The username that is used to connect to the System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

* `password` - (Required) The password that is used to connect to the System Center Virtual Machine Manager Server. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the System Center Virtual Machine Manager Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this System Center Virtual Machine Manager Server.
* `read` - (Defaults to 5 minutes) Used when retrieving this System Center Virtual Machine Manager Server.
* `update` - (Defaults to 30 minutes) Used when updating this System Center Virtual Machine Manager Server.
* `delete` - (Defaults to 30 minutes) Used when deleting this System Center Virtual Machine Manager Server.

## Import

System Center Virtual Machine Manager Servers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_system_center_virtual_machine_manager_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.ScVmm/vmmServers/vmmServer1
```
