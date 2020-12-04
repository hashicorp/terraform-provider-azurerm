---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub"
description: |-
  Manages an IotHub
---

# azurerm_iothub

Manages an IotHub

~> **NOTE:** Endpoints can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_endpoint_*` resources - but the two ways of defining the endpoints cannot be used together. If both are used against the same IoTHub, spurious changes will occur. Also, defining a `azurerm_iothub_endpoint_*` resource and another endpoint of a different type directly on the `azurerm_iothub` resource is not supported.

~> **NOTE:** Routes can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_route` resource - but the two cannot be used together. If both are used against the same IoTHub, spurious changes will occur.

~> **NOTE:** Fallback route can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_fallback_route` resource - but the two cannot be used together. If both are used against the same IoTHub, spurious changes will occur.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "Canada Central"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorage"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "examplecontainer"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-namespace"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Basic"
}

resource "azurerm_eventhub" "example" {
  name                = "example-eventhub"
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_eventhub_namespace.example.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "example" {
  resource_group_name = azurerm_resource_group.example.name
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  name                = "acctest"
  send                = true
}

resource "azurerm_iothub" "example" {
  name                = "Example-IoTHub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  endpoint {
    type                       = "AzureIotHub.StorageContainer"
    connection_string          = azurerm_storage_account.example.primary_blob_connection_string
    name                       = "export"
    batch_frequency_in_seconds = 60
    max_chunk_size_in_bytes    = 10485760
    container_name             = azurerm_storage_container.example.name
    encoding                   = "Avro"
    file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
  }

  endpoint {
    type              = "AzureIotHub.EventHub"
    connection_string = azurerm_eventhub_authorization_rule.example.primary_connection_string
    name              = "export2"
  }

  route {
    name           = "export"
    source         = "DeviceMessages"
    condition      = "true"
    endpoint_names = ["export"]
    enabled        = true
  }

  route {
    name           = "export2"
    source         = "DeviceMessages"
    condition      = "true"
    endpoint_names = ["export2"]
    enabled        = true
  }

  tags = {
    purpose = "testing"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the IotHub resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the IotHub resource has to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource has to be createc. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `event_hub_partition_count` - (Optional) The number of device-to-cloud partitions used by backing event hubs. Must be between `2` and `128`.

* `event_hub_retention_in_days` - (Optional) The event hub retention to use in days. Must be between `1` and `7`.

* `endpoint` - (Optional) An `endpoint` block as defined below.

* `fallback_route` - (Optional) A `fallback_route` block as defined below. If the fallback route is enabled, messages that don't match any of the supplied routes are automatically sent to this route. Defaults to messages/events.

~> **NOTE:** If `fallback_route` isn't explicitly specified, the fallback route wouldn't be enabled by default.

* `file_upload` - (Optional) A `file_upload` block as defined below.

* `ip_filter_rule` - (Optional) One or more `ip_filter_rule` blocks as defined below.

* `route` - (Optional) A `route` block as defined below.

* `public_network_access_enabled` - (Optional) Is the IotHub resource accessible from a public network? 

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) The name of the sku. Possible values are `B1`, `B2`, `B3`, `F1`, `S1`, `S2`, and `S3`.

* `capacity` - (Required) The number of provisioned IoT Hub units.

~> **NOTE:** Only one IotHub can be on the `Free` tier per subscription.

---

An `endpoint` block supports the following:

* `type` - (Required) The type of the endpoint. Possible values are `AzureIotHub.StorageContainer`, `AzureIotHub.ServiceBusQueue`, `AzureIotHub.ServiceBusTopic` or `AzureIotHub.EventHub`.

* `connection_string` - (Required) The connection string for the endpoint.

* `name` - (Required) The name of the endpoint. The name must be unique across endpoint types. The following names are reserved:  `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`.

* `batch_frequency_in_seconds` - (Optional) Time interval at which blobs are written to storage. Value should be between 60 and 720 seconds. Default value is 300 seconds. This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

* `max_chunk_size_in_bytes` - (Optional) Maximum number of bytes for each blob written to storage. Value should be between 10485760(10MB) and 524288000(500MB). Default value is 314572800(300MB). This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

* `container_name` - (Optional) The name of storage container in the storage account. This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

* `encoding` - (Optional) Encoding that is used to serialize messages to blobs. Supported values are 'avro' and 'avrodeflate'. Default value is 'avro'. This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

* `file_name_format` - (Optional) File name format for the blob. Default format is ``{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}``. All parameters are mandatory but can be reordered. This attribute is mandatory for endpoint type `AzureIotHub.StorageContainer`.

* `resource_group_name` - (Optional) The resource group in which the endpoint will be created.

---

An `ip_filter_rule` block supports the following:

* `name` - (Required) The name of the filter.

* `ip_mask` - (Required) The IP address range in CIDR notation for the rule.

* `action` - (Required) The desired action for requests captured by this rule. Possible values are  `Accept`, `Reject`

---

A `route` block supports the following:

* `name` - (Required) The name of the route.

* `source` - (Required) The source that the routing rule is to be applied to, such as `DeviceMessages`. Possible values include: `RoutingSourceInvalid`, `RoutingSourceDeviceMessages`, `RoutingSourceTwinChangeEvents`, `RoutingSourceDeviceLifecycleEvents`, `RoutingSourceDeviceJobLifecycleEvents`.

* `condition` - (Optional) The condition that is evaluated to apply the routing rule. If no condition is provided, it evaluates to true by default. For grammar, see: https://docs.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language.

* `endpoint_names` - (Required) The list of endpoints to which messages that satisfy the condition are routed.

* `enabled` - (Required) Used to specify whether a route is enabled.

---

A `fallback_route` block supports the following:

* `source` - (Optional) The source that the routing rule is to be applied to, such as `DeviceMessages`. Possible values include: `RoutingSourceInvalid`, `RoutingSourceDeviceMessages`, `RoutingSourceTwinChangeEvents`, `RoutingSourceDeviceLifecycleEvents`, `RoutingSourceDeviceJobLifecycleEvents`.

* `condition` - (Optional) The condition that is evaluated to apply the routing rule. If no condition is provided, it evaluates to true by default. For grammar, see: https://docs.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language.

* `endpoint_names` - (Optional) The endpoints to which messages that satisfy the condition are routed. Currently only 1 endpoint is allowed.

* `enabled` - (Optional) Used to specify whether the fallback route is enabled.

---

A `file_upload` block supports the following:

* `connection_string` - (Required) The connection string for the Azure Storage account to which files are uploaded.

* `container_name` - (Required) The name of the root container where you upload files. The container need not exist but should be creatable using the connection_string specified.

* `sas_ttl` - (Optional) The period of time for which the SAS URI generated by IoT Hub for file upload is valid, specified as an [ISO 8601 timespan duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This value must be between 1 minute and 24 hours, and evaluates to 'PT1H' by default.

* `notifications` - (Optional) Used to specify whether file notifications are sent to IoT Hub on upload. It evaluates to false by default.

* `lock_duration` - (Optional) The lock duration for the file upload notifications queue, specified as an [ISO 8601 timespan duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This value must be between 5 and 300 seconds, and evaluates to 'PT1M' by default.

* `default_ttl` - (Optional) The period of time for which a file upload notification message is available to consume before it is expired by the IoT hub, specified as an [ISO 8601 timespan duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). This value must be between 1 minute and 48 hours, and evaluates to 'PT1H' by default.

* `max_delivery_count` - (Optional) The number of times the IoT hub attempts to deliver a file upload notification message. It evaluates to 10 by default.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTHub.

* `event_hub_events_endpoint` -  The EventHub compatible endpoint for events data
* `event_hub_events_path` -  The EventHub compatible path for events data
* `event_hub_operations_endpoint` -  The EventHub compatible endpoint for operational data
* `event_hub_operations_path` -  The EventHub compatible path for operational data

-> **NOTE:** These fields can be used in conjunction with the `shared_access_policy` block to build a connection string

* `hostname` - The hostname of the IotHub Resource.

* `shared_access_policy` - One or more `shared_access_policy` blocks as defined below.

---

A `shared access policy` block contains the following:

* `key_name` - The name of the shared access policy.

* `primary_key` - The primary key.

* `secondary_key` - The secondary key.

* `permissions` - The permissions assigned to the shared access policy.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub.

## Import

IoTHubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub.hub1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1
```
