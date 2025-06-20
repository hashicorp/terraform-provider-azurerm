---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_certificate"
description: |-
  Manages an IoTHub Certificate.
---

# azurerm_iothub_certificate

Manages an IotHub Certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iothub" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "B1"
    capacity = "1"
  }
}

resource "azurerm_iothub_certificate" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  iothub_name         = azurerm_iothub.example.name
  is_verified         = true

  certificate_content = filebase64("example.cer")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the IotHub Certificate resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the IotHub Certificate resource has to be created. Changing this forces a new resource to be created.

* `iothub_name` - (Required) The name of the IoTHub that this certificate will be attached to. Changing this forces a new resource to be created.

* `certificate_content` - (Required) The Base-64 representation of the X509 leaf certificate .cer file or just a .pem file content.

* `is_verified` - (Optional) Is the certificate verified? Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoTHub Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub Certificate.

## Import

IoTHub Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/iotHubs/example/certificates/example
```
