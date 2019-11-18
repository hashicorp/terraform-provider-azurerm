---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_dps_certificate"
sidebar_current: "docs-azurerm-resource-messaging-iot-dps_certificate"
description: |-
  Manages an IoT Device Provisioning Service Certificate.
---

# azurerm_iot_dps_certificate

Manages an IoT Device Provisioning Service Certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_iot_dps" "example" {
  name                = "example"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
}

resource "azurerm_iot_dps_certificate" "example" {
  name                = "example"
  resource_group_name = "${azurerm_resource_group.example.name}"
  iot_dps_name        = "${azurerm_iot_dps.example.name}"

  certificate_content = "${filebase64("example.cer")}"
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

* `id` - The ID of the IoT Device Provisioning Service Certificate.

## Import

IoT Device Provisioning Service Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_dps_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/provisioningServices/example/certificates/example
```
