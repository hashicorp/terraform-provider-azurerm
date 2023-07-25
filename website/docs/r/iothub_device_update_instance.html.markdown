---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_device_update_instance"
description: |-
  Manages an IoT Hub Device Update Instance.
---

# azurerm_iothub_device_update_instance

Manages an IoT Hub Device Update Instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_iothub_device_update_account" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_iothub" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_iothub_device_update_instance" "example" {
  name                     = "example"
  device_update_account_id = azurerm_iothub_device_update_account.example.id
  iothub_id                = azurerm_iothub.example.id
  diagnostic_enabled       = true

  diagnostic_storage_account {
    connection_string = azurerm_storage_account.example.primary_connection_string
    id                = azurerm_storage_account.example.id
  }

  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this IoT Hub Device Update Instance. Changing this forces a new resource to be created.

* `device_update_account_id` - (Required) Specifies the ID of the IoT Hub Device Update Account where the IoT Hub Device Update Instance exists. Changing this forces a new resource to be created.

* `iothub_id` - (Required) Specifies the ID of the IoT Hub associated with the IoT Hub Device Update Instance. Changing this forces a new resource to be created.

* `diagnostic_storage_account` - (Optional) A `diagnostic_storage_account` block as defined below.

* `diagnostic_enabled` - (Optional) Whether the diagnostic log collection is enabled. Possible values are `true` and `false`. Defaults to `false`.

* `tags` - (Optional) A mapping of tags which should be assigned to the IoT Hub Device Update Instance.

---

A `diagnostic_storage_account` block supports the following:

* `connection_string` - (Required) Connection String of the Diagnostic Storage Account.

* `id` - (Required) Resource ID of the Diagnostic Storage Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoT Hub Device Update Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Hub Device Update Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Hub Device Update Instance.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Hub Device Update Instance.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Hub Device Update Instance.

## Import

IoT Hub Device Update Instance can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_device_update_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DeviceUpdate/accounts/account1/instances/instance1
```
