---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub"
sidebar_current: "docs-azurerm-resource-iothub"
description: |-
  Manages a IotHub resource 
---

# azurerm_iothub

Manages a IotHub

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_iothub" "test" {
  name                = "test"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku {
    name = "S1"
    tier = "Standard"
    capacity = "1"
  }

  tags {
    "purpose" = "testing"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the IotHub resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the IotHub resource has to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource has to be createc. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below. 

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) The name of the sku. Possible values are `F1`, `S1`, `S2`, and `S3`.

* `tier` - (Required) The billing tier for the IoT Hub. Possible values are `Free` or `Standard`.

* `capacity` - (Required) The number of provisioned IoT Hub units. 

## Attributes Reference

The following attributes are exported:

* `id` - The IotHub ID.

* `hostname` - The hostname of the IotHub Resource.

* `shared_access_policy` - A list of `shared_access_policy` blocks as defined below.

---

A `shared access policy` block contains the following:

* `key_name` - The name of the shared access policy.

* `primary_key` - The primary key.

* `secondary_key` - The secondary key.

* `permissions` - The permissions assigned to the shared access policy.

## Import

IoTHubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub.hub1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1
```
