---
subcategory: "IoT Hub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_endpoint_storage_container"
description: |-
  Manages an IotHub Storage Container Endpoint
---

# azurerm_iothub_endpoint_storage_container

Manages an IotHub Storage Container Endpoint

~> **Note:** Endpoints can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_endpoint_*` resources - but the two ways of defining the endpoints cannot be used together. If both are used against the same IoTHub, spurious changes will occur. Also, defining a `azurerm_iothub_endpoint_*` resource and another endpoint of a different type directly on the `azurerm_iothub` resource is not supported.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "acctestcont"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_iothub" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "example" {
  resource_group_name = azurerm_resource_group.example.name
  iothub_id           = azurerm_iothub.example.id
  name                = "acctest"

  container_name    = "acctestcont"
  connection_string = azurerm_storage_account.example.primary_blob_connection_string

  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  encoding                   = "JSON"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the endpoint. The name must be unique across endpoint types. The following names are reserved: `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the Storage Container has been created. Changing this forces a new resource to be created.

* `container_name` - (Required) The name of storage container in the storage account.

* `iothub_id` - (Required) The IoTHub ID for the endpoint. Changing this forces a new resource to be created.

* `authentication_type` - (Optional) Type used to authenticate against the storage endpoint. Possible values are `keyBased` and `identityBased`. Defaults to `keyBased`.

* `identity_id` - (Optional) ID of the User Managed Identity used to authenticate against the storage endpoint.

-> **Note:** `identity_id` can only be specified when `authentication_type` is `identityBased`. It must be one of the `identity_ids` of the Iot Hub. If not specified when `authentication_type` is `identityBased`, System Assigned Managed Identity of the Iot Hub will be used.

* `endpoint_uri` - (Optional) URI of the Storage Container endpoint. This corresponds to the `primary_blob_endpoint` of the parent storage account. This attribute can only be specified and is mandatory when `authentication_type` is `identityBased`.

* `connection_string` - (Optional) The connection string for the endpoint. This attribute can only be specified and is mandatory when `authentication_type` is `keyBased`.

* `batch_frequency_in_seconds` - (Optional) Time interval at which blobs are written to storage. Value should be between 60 and 720 seconds. Default value is 300 seconds.

* `max_chunk_size_in_bytes` - (Optional) Maximum number of bytes for each blob written to storage. Value should be between 10485760(10MB) and 524288000(500MB). Default value is 314572800(300MB).

* `encoding` - (Optional) Encoding that is used to serialize messages to blobs. Supported values are `Avro`, `AvroDeflate` and `JSON`. Default value is `Avro`. Changing this forces a new resource to be created.

* `file_name_format` - (Optional) File name format for the blob. All parameters are mandatory but can be reordered. Defaults to `{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoTHub Storage Container Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub Storage Container Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Storage Container Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub Storage Container Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub Storage Container Endpoint.

## Import

IoTHub Storage Container Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_endpoint_storage_container.storage_container1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/iotHubs/hub1/endpoints/storage_container_endpoint1
```
