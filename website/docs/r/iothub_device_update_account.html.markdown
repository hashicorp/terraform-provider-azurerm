---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_device_update_account"
description: |-
  Manages an IoT Hub Device Update Account.
---

# azurerm_iothub_device_update_account

Manages an IoT Hub Device Update Account.

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

  identity {
    type = "SystemAssigned"
  }

  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this IoT Hub Device Update Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the IoT Hub Device Update Account should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the IoT Hub Device Update Account should exist. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `public_network_access_enabled` - (Optional) Specifies whether the public network access is enabled for the IoT Hub Device Update Account. Possible values are `true` and `false`. Defaults to `true`.

* `sku` - (Optional) Sku of the IoT Hub Device Update Account. Possible values are `Free` and `Standard`. Defaults to `Standard`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the IoT Hub Device Update Account.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this IoT Hub Device Update Account. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this IoT Hub Device Update Account.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoT Hub Device Update Account.

* `host_name` - The API host name of the IoT Hub Device Update Account.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this IoT Hub Device Update Account.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this IoT Hub Device Update Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Hub Device Update Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Hub Device Update Account.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Hub Device Update Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Hub Device Update Account.

## Import

IoT Hub Device Update Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_device_update_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DeviceUpdate/accounts/account1
```
