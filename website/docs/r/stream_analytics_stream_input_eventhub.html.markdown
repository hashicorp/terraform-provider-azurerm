---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_stream_input_eventhub"
description: |-
  Manages a Stream Analytics Stream Input EventHub.
---

# azurerm_stream_analytics_stream_input_eventhub

~> **Note:** This resource creates a Stream Input of type `Microsoft.ServiceBus/EventHub`, to create a Stream Input of type `Microsoft.EventHub/EventHub` please use the resource azurerm_stream_analytics_stream_input_eventhub_v2.

Manages a Stream Analytics Stream Input EventHub.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_stream_analytics_job" "example" {
  name                = "example-job"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "example" {
  name                = "example-eventhub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "example" {
  name                = "example-consumergroup"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_stream_analytics_stream_input_eventhub" "example" {
  name                         = "eventhub-stream-input"
  stream_analytics_job_name    = data.azurerm_stream_analytics_job.example.name
  resource_group_name          = data.azurerm_stream_analytics_job.example.resource_group_name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.example.name
  eventhub_name                = azurerm_eventhub.example.name
  servicebus_namespace         = azurerm_eventhub_namespace.example.name
  shared_access_policy_key     = azurerm_eventhub_namespace.example.default_primary_key
  shared_access_policy_name    = "RootManageSharedAccessKey"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Input EventHub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `eventhub_name` - (Required) The name of the Event Hub.

* `servicebus_namespace` - (Required) The namespace that is associated with the desired Event Hub, Service Bus Queue, Service Bus Topic, etc.

* `serialization` - (Required) A `serialization` block as defined below.

* `partition_key` - (Optional) The property the input Event Hub has been partitioned by.

* `authentication_mode` - (Optional) The authentication mode for the Stream Output. Possible values are `Msi` and `ConnectionString`. Defaults to `ConnectionString`.

* `eventhub_consumer_group_name` - (Optional) The name of an Event Hub Consumer Group that should be used to read events from the Event Hub. Specifying distinct consumer group names for multiple inputs allows each of those inputs to receive the same events from the Event Hub. If not set the input will use the Event Hub's default consumer group.

* `shared_access_policy_key` - (Optional) The shared access policy key for the specified shared access policy.

* `shared_access_policy_name` - (Optional) The shared access policy name for the Event Hub, Service Bus Queue, Service Bus Topic, etc.

---

A `serialization` block supports the following:

* `type` - (Required) The serialization format used for incoming data streams. Possible values are `Avro`, `Csv` and `Json`.

* `encoding` - (Optional) The encoding of the incoming data in the case of input and the encoding of outgoing data in the case of output. It currently can only be set to `UTF8`.

-> **Note:** This is required when `type` is set to `Csv` or `Json`.

* `field_delimiter` - (Optional) The delimiter that will be used to separate comma-separated value (CSV) records. Possible values are ` ` (space), `,` (comma), `	` (tab), `|` (pipe) and `;`.

-> **Note:** This is required when `type` is set to `Csv`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Stream Input EventHub.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Stream Input EventHub.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Stream Input EventHub.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Stream Input EventHub.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Stream Input EventHub.

## Import

Stream Analytics Stream Input EventHub's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_stream_input_eventhub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/inputs/input1
```
