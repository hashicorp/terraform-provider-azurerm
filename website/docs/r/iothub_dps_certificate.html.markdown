---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_dps_certificate"
description: |-
  Manages an IoT Device Provisioning Service Certificate.
---

# azurerm_iothub_dps_certificate

Manages an IotHub Device Provisioning Service Certificate. 

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iothub_dps" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_iothub_dps_certificate" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  iot_dps_name        = azurerm_iothub_dps.example.name

  certificate_content = filebase64("example.cer")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Iot Device Provisioning Service Certificate resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the Iot Device Provisioning Service Certificate resource has to be created. Changing this forces a new resource to be created.

* `iot_dps_name` - (Required) The name of the IoT Device Provisioning Service that this certificate will be attached to. Changing this forces a new resource to be created.

* `certificate_content` - (Required) The Base-64 representation of the X509 leaf certificate .cer file or just a .pem file content.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTHub Device Provisioning Service Certificate.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub Device Provisioning Service Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub Device Provisioning Service Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Device Provisioning Service Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub Device Provisioning Service Certificate.

## Import

IoTHub Device Provisioning Service Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_dps_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/provisioningServices/example/certificates/example
```
