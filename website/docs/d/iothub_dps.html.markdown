---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_dps"
description: |-
  Gets information about an existing IoT Hub Device Provisioning Service.
---

# Data Source: azurerm_iothub_dps

Use this data source to access information about an existing IotHub Device Provisioning Service.

## Example Usage

```hcl
data "azurerm_iothub_dps" "example" {
  name                = "iot_hub_dps_test"
  resource_group_name = "iothub_dps_rg"
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Iot Device Provisioning Service resource.

* `resource_group_name` - The name of the resource group under which the Iot Device Provisioning Service is located in.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Device Provisioning Service.

* `location` - Specifies the supported Azure location where the IoT Device Provisioning Service exists.

* `allocation_policy` - The allocation policy of the IoT Device Provisioning Service.

* `device_provisioning_host_name` - The device endpoint of the IoT Device Provisioning Service.

* `id_scope` - The unique identifier of the IoT Device Provisioning Service.

* `service_operations_host_name` - The service endpoint of the IoT Device Provisioning Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Hub Device Provisioning Service.
