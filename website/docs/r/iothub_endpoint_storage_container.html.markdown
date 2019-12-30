---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_endpoint_storage_container"
sidebar_current: "docs-azurerm-resource-messaging-iothub-endpoint-storage-container"
description: |-
  Manages an IotHub Storage Container Endpoint
---

# azurerm_iothub_endpoint_storage_container

Manages an IotHub Storage Container Endpoint

~> **NOTE:** Endpoints can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_endpoint_*` resources - but the two ways of defining the endpoints cannot be used together. If both are used against the same IoTHub, spurious changes will occur. Also, defining a `azurerm_iothub_endpoint_*` resource and another endpoint of a different type directly on the `azurerm_iothub` resource is not supported.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "acctestcont"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  storage_account_name  = "${azurerm_storage_account.example.name}"
  container_access_type = "private"
}

resource "azurerm_iothub" "example" {
  name                = "example"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  iothub_name         = "${azurerm_iothub.example.name}"
  name                = "acctest"

  container_name    = "acctestcont"
  connection_string = "${azurerm_storage_account.example.primary_blob_connection_string}"

  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  encoding                   = "JSON"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the endpoint. The name must be unique across endpoint types. The following names are reserved:  `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`.

* `resource_group_name` - (Required) The name of the resource group under which the IotHub Storage Container Endpoint resource has to be created. Changing this forces a new resource to be created.

* `iothub_name` - (Required) The name of the IoTHub to which this Storage Container Endpoint belongs. Changing this forces a new resource to be created.

* `connection_string` - (Required) The connection string for the endpoint.

* `batch_frequency_in_seconds` - (Optional) Time interval at which blobs are written to storage. Value should be between 60 and 720 seconds. Default value is 300 seconds. 

* `max_chunk_size_in_bytes` - (Optional) Maximum number of bytes for each blob written to storage. Value should be between 10485760(10MB) and 524288000(500MB). Default value is 314572800(300MB).

* `container_name` - (Required) The name of storage container in the storage account.
* 
* `encoding` - (Optional) Encoding that is used to serialize messages to blobs. Supported values are 'avro' and 'avrodeflate'. Default value is 'avro'.

* `file_name_format` - (Optional) File name format for the blob. Default format is ``{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}``. All parameters are mandatory but can be reordered.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTHub Storage Container Endpoint.

## Import

IoTHub Storage Container Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_endpoint_storage_container.storage_container1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1/Endpoints/storage_container_endpoint1
```
