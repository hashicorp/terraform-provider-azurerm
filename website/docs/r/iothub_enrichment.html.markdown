---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iothub_enrichment"
description: |-
  Manages an IotHub Enrichment
---

# azurerm_iothub_enrichment

Manages an IotHub Enrichment

~> **NOTE:** Enrichment can be defined either directly on the `azurerm_iothub` resource, or using the `azurerm_iothub_enrichment` resources - but the two cannot be used together. If both are used against the same IoTHub, spurious changes will occur.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "example"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_iothub" "example" {
  name                = "exampleIothub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "example" {
  resource_group_name = azurerm_resource_group.example.name
  iothub_name         = azurerm_iothub.example.name
  name                = "example"

  connection_string          = azurerm_storage_account.example.primary_blob_connection_string
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  container_name             = azurerm_storage_container.example.name
  encoding                   = "Avro"
  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
}

resource "azurerm_iothub_route" "example" {
  resource_group_name = azurerm_resource_group.example.name
  iothub_name         = azurerm_iothub.example.name
  name                = "example"

  source         = "DeviceMessages"
  condition      = "true"
  endpoint_names = [azurerm_iothub_endpoint_storage_container.example.name]
  enabled        = true
}

resource "azurerm_iothub_enrichment" "example" {
  resource_group_name = azurerm_resource_group.example.name
  iothub_name         = azurerm_iothub.example.name
  key                 = "example"

  value          = "my value"
  endpoint_names = [azurerm_iothub_endpoint_storage_container.example.name]
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) The key of the enrichment.

* `value` - (Required) The value of the enrichment. Value can be any static string, the name of the IoT hub sending the message (use `$iothubname`) or information from the device twin (ex: `$twin.tags.latitude`)

* `endpoint_names` - (Required) The list of endpoints which will be enriched.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTHub Enrichment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IotHub Enrichment.
* `update` - (Defaults to 30 minutes) Used when updating the IotHub Enrichment.
* `read` - (Defaults to 5 minutes) Used when retrieving the IotHub Enrichment.
* `delete` - (Defaults to 30 minutes) Used when deleting the IotHub Enrichment.

## Import

IoTHub Enrichment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iothub_enrichment.enrichment1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Devices/IotHubs/hub1/Enrichments/enrichment1
```
