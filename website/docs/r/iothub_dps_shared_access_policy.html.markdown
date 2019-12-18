---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_dps_shared_access_policy"
sidebar_current: "docs-azurerm-resource-iothub-dps-shared-access-policy-x"
description: |-
  Manages an IotHub Device Provisioning Service Shared Access Policy
---

# azurerm_iothub_dps_shared_access_policy

Manages an IotHub Device Provisioning Service Shared Access Policy

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
    tier     = "Standard"
    capacity = "1"
  }
}

resource "azurerm_iothub_dps_shared_access_policy" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  iothub_dps_name     = azurerm_iothub_dps.example.name

  enrollment_write  = true
  enrollment_read   = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the IotHub Shared Access Policy resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the IotHub Shared Access Policy resource has to be created. Changing this forces a new resource to be created.

* `iothub_dps_name` - (Required) The name of the IoT Hub Device Provisioning service to which this Shared Access Policy belongs. Changing this forces a new resource to be created.

* `enrollment_read` - (Optional) Adds `EnrollmentRead` permission to this Shared Access Account. It allows read access to enrollment data.

-> **NOTE** When `enrollment_read` is set to `true`, `registration_read` must also be set to true. This is a limitation of the Azure REST API

* `enrollment_write` - (Optional) Adds `EnrollmentWrite` permission to this Shared Access Account. It allows write access to enrollment data.

-> **NOTE** When `registration_write` is set to `true`, `enrollment_read`, `registration_read`, and `registration_write` must also be set to true. This is a requirement of the Azure API.

* `registration_read` - (Optional) Adds `RegistrationStatusRead` permission to this Shared Access Account. It allows read access to device registrations.

* `registration_write` - (Optional) Adds `RegistrationStatusWrite` permission to this Shared Access Account. It allows write access to device registrations.

-> **NOTE** When `registration_write` is set to `true`, `registration_read` must also be set to true. This is a requirement of the Azure API.

* `service_config` - (Optional) Adds `ServiceConfig` permission to this Shared Access Account. It allows configuration of the Device Provisioning Service.

-> **NOTE** At least one of `registration_read`, `registration_write`, `service_config`, `enrollment_read`, `enrollment_write` permissions must be set to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Hub Device Provisioning Service Shared Access Policy.

* `primary_key` - The primary key used to create the authentication token.

* `secondary_key` - The secondary key used to create the authentication token.

## Import

IoT Hub Device Provisioning Service Shared Access Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_dps_shared_access_policy.shared_access_policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/provisioningServices/dps1/keys/shared_access_policy1
```