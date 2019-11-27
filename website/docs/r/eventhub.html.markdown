---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub"
sidebar_current: "docs-azurerm-resource-messaging-eventhub-x"
description: |-
  Manages a Event Hubs as a nested resource within an Event Hubs namespace.
---

# azurerm_eventhub

Manages a Event Hubs as a nested resource within a Event Hubs namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "acceptanceTestEventHubNamespace"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "Standard"
  capacity            = 1

  tags = {
    environment = "Production"
  }
}

resource "azurerm_eventhub" "example" {
  name                = "acceptanceTestEventHub"
  namespace_name      = "${azurerm_eventhub_namespace.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  partition_count     = 2
  message_retention   = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Namespace resource. Changing this forces a new resource to be created.

* `namespace_name` - (Required) Specifies the name of the EventHub Namespace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventHub's parent Namespace exists. Changing this forces a new resource to be created.

* `partition_count` - (Required) Specifies the current number of shards on the Event Hub. Changing this forces a new resource to be created.

* `message_retention` - (Required) Specifies the number of days to retain the events for this Event Hub. Needs to be between 1 and 7 days; or 1 day when using a Basic SKU for the parent EventHub Namespace.

* `capture_description` - (Optional) A `capture_description` block as defined below.

---

A `capture_description` block supports the following:

* `enabled` - (Required) Specifies if the Capture Description is Enabled.

* `encoding` - (Required) Specifies the Encoding used for the Capture Description. Possible values are `Avro` and `AvroDeflate`.

* `interval_in_seconds` - (Optional) Specifies the time interval in seconds at which the capture will happen. Values can be between `60` and `900` seconds. Defaults to `300` seconds.

* `size_limit_in_bytes` - (Optional) Specifies the amount of data built up in your EventHub before a Capture Operation occurs. Value should be between `10485760` and `524288000`  bytes. Defaults to `314572800` bytes.

* `skip_empty_archives` - (Optional) Specifies if empty files should not be emitted if no events occur during the Capture time window.  Defaults to `false`.

* `destination` - (Required) A `destination` block as defined below.

A `destination` block supports the following:

* `name` - (Required) The Name of the Destination where the capture should take place. At this time the only supported value is `EventHubArchive.AzureBlockBlob`.

-> At this time it's only possible to Capture EventHub messages to Blob Storage. There's [a Feature Request for the Azure SDK to add support for Capturing messages to Azure Data Lake here](https://github.com/Azure/azure-rest-api-specs/issues/2255).

* `archive_name_format` - The Blob naming convention for archiving. e.g. `{Namespace}/{EventHub}/{PartitionId}/{Year}/{Month}/{Day}/{Hour}/{Minute}/{Second}`. Here all the parameters (Namespace,EventHub .. etc) are mandatory irrespective of order

* `blob_container_name` - (Required) The name of the Container within the Blob Storage Account where messages should be archived.

* `storage_account_id` - (Required) The ID of the Blob Storage Account where messages should be archived.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub ID.

* `partition_ids` - The identifiers for partitions created for Event Hubs.

## Import

EventHubs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub.eventhub1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1
```
