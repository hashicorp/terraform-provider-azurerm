---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_stream_input_iothub"
description: |-
  Manages a Stream Analytics Stream Input IoTHub.
---

# azurerm_stream_analytics_stream_input_iothub

Manages a Stream Analytics Stream Input IoTHub.

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

resource "azurerm_iothub" "example" {
  name                = "example-iothub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_stream_analytics_stream_input_iothub" "example" {
  name                         = "example-iothub-input"
  stream_analytics_job_name    = data.azurerm_stream_analytics_job.example.name
  resource_group_name          = data.azurerm_stream_analytics_job.example.resource_group_name
  endpoint                     = "messages/events"
  eventhub_consumer_group_name = "$Default"
  iothub_namespace             = azurerm_iothub.example.name
  shared_access_policy_key     = azurerm_iothub.example.shared_access_policy[0].primary_key
  shared_access_policy_name    = "iothubowner"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Input IoTHub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `eventhub_consumer_group_name` - (Required) The name of an Event Hub Consumer Group that should be used to read events from the Event Hub. Specifying distinct consumer group names for multiple inputs allows each of those inputs to receive the same events from the Event Hub.

* `endpoint` - (Required) The IoT Hub endpoint to connect to (ie. messages/events, messages/operationsMonitoringEvents, etc.).

* `iothub_namespace` - (Required) The name or the URI of the IoT Hub.

* `serialization` - (Required) A `serialization` block as defined below.

* `shared_access_policy_key` - (Required) The shared access policy key for the specified shared access policy. Changing this forces a new resource to be created.

* `shared_access_policy_name` - (Required) The shared access policy name for the Event Hub, Service Bus Queue, Service Bus Topic, etc.

---

A `serialization` block supports the following:

* `type` - (Required) The serialization format used for incoming data streams. Possible values are `Avro`, `Csv` and `Json`.

* `encoding` - (Optional) The encoding of the incoming data in the case of input and the encoding of outgoing data in the case of output. It currently can only be set to `UTF8`.

-> **Note:** This is required when `type` is set to `Csv` or `Json`.

* `field_delimiter` - (Optional) The delimiter that will be used to separate comma-separated value (CSV) records. Possible values are ` ` (space), `,` (comma), `	` (tab), `|` (pipe) and `;`.

-> **Note:** This is required when `type` is set to `Csv`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Stream Input IoTHub.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Stream Input IoTHub.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Stream Input IoTHub.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Stream Input IoTHub.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Stream Input IoTHub.

## Import

Stream Analytics Stream Input IoTHub's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_stream_input_iothub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/inputs/input1
```
