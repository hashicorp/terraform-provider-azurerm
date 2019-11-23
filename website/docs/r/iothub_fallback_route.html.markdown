---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_fallback_route"
sidebar_current: "docs-azurerm-resource-messaging-iothub-fallback-route-x"
description: |-
  Manages an IotHub Fallback Route
---

# azurerm_iothub_fallback_route

Manages an IotHub Fallback Route

~> **NOTE:** Fallback route can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_fallback_route` resource - but the two cannot be used together. If both are used against the same IoTHub, spurious changes will occur.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West US"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "example"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  storage_account_name  = "${azurerm_storage_account.example.name}"
  container_access_type = "private"
}

resource "azurerm_iothub" "example" {
  name                = "exampleIothub"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  iothub_name         = "${azurerm_iothub.example.name}"
  name                = "example"

  connection_string          = "${azurerm_storage_account.example.primary_blob_connection_string}"
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  container_name             = "${azurerm_storage_container.example.name}"
  encoding                   = "Avro"
  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
}

resource "azurerm_iothub_fallback_route" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  iothub_name         = "${azurerm_iothub.example.name}"

  condition      = "true"
  endpoint_names = ["${azurerm_iothub_endpoint_storage_container.example.name}"]
  enabled        = true
}

```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group under which the IotHub Storage Container Endpoint resource has to be created. Changing this forces a new resource to be created.

* `iothub_name` - (Required) The name of the IoTHub to which this Fallback Route belongs. Changing this forces a new resource to be created.

* `enabled` - (Required) Used to specify whether the fallback route is enabled.

* `endpoint_names` - (Required) The endpoints to which messages that satisfy the condition are routed. Currently only 1 endpoint is allowed.

* `condition` - (Optional) The condition that is evaluated to apply the routing rule. If no condition is provided, it evaluates to `true` by default. For grammar, see: https://docs.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTHub FallbackRoute.

## Import

IoTHub Fallback Route can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_route.route1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1/FallbackRoute/default
```
~> **NOTE:** As there may only be a single fallback route per IoTHub, the id always ends with `/FallbackRoute/default`.
