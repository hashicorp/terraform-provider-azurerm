---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_dps_shared_access_policy"
description: |-
  Gets information about an existing IotHub Device Provisioning Service Shared Access Policy
---

# Data Source: azurerm_iothub_dps_shared_access_policy

Use this data source to access information about an existing IotHub Device Provisioning Service Shared Access Policy

## Example Usage

```hcl
data "azurerm_iothub_dps_shared_access_policy" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  iothub_dps_name     = azurerm_iothub_dps.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the IotHub Shared Access Policy.

* `resource_group_name` - Specifies the name of the resource group under which the IotHub Shared Access Policy resource exists.

* `iothub_dps_name` - Specifies the name of the IoT Hub Device Provisioning service to which the Shared Access Policy belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Hub Device Provisioning Service Shared Access Policy.

* `primary_key` - The primary key used to create the authentication token.

* `primary_connection_string` - The primary connection string of the Shared Access Policy.

* `secondary_key` - The secondary key used to create the authentication token.

* `secondary_connection_string` - The secondary connection string of the Shared Access Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Device Provisioning Service Shared Access Policy.
