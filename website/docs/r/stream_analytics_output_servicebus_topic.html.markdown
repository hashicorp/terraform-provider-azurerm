---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_servicebus_topic"
description: |-
  Manages a Stream Analytics Output to a ServiceBus Topic.
---

# azurerm_stream_analytics_output_servicebus_topic

Manages a Stream Analytics Output to a ServiceBus Topic.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

data "azurerm_stream_analytics_job" "example" {
  name                = "example-job"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "example-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "example" {
  name                = "example-topic"
  namespace_id        = azurerm_servicebus_namespace.example.id
  enable_partitioning = true
}

resource "azurerm_stream_analytics_output_servicebus_topic" "example" {
  name                      = "service-bus-topic-output"
  stream_analytics_job_name = data.azurerm_stream_analytics_job.example.name
  resource_group_name       = data.azurerm_stream_analytics_job.example.resource_group_name
  topic_name                = azurerm_servicebus_topic.example.name
  servicebus_namespace      = azurerm_servicebus_namespace.example.name
  shared_access_policy_key  = azurerm_servicebus_namespace.example.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"
  property_columns          = ["col1", "col2"]

  serialization {
    type   = "Csv"
    format = "Array"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Output. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `topic_name` - (Required) The name of the Service Bus Topic.

* `servicebus_namespace` - (Required) The namespace that is associated with the desired Event Hub, Service Bus Topic, Service Bus Topic, etc.

* `shared_access_policy_key` - (Optional) The shared access policy key for the specified shared access policy. Required if `authentication_mode` is `ConnectionString`.

* `shared_access_policy_name` - (Optional) The shared access policy name for the Event Hub, Service Bus Queue, Service Bus Topic, etc. Required if `authentication_mode` is `ConnectionString`.

* `serialization` - (Required) A `serialization` block as defined below.

* `property_columns` - (Optional) A list of property columns to add to the Service Bus Topic output.

* `authentication_mode` - (Optional) The authentication mode for the Stream Output. Possible values are `Msi` and `ConnectionString`. Defaults to `ConnectionString`.

* `system_property_columns` - (Optional) A key-value pair of system property columns that will be attached to the outgoing messages for the Service Bus Topic Output.

-> **Note:** The acceptable keys are `ContentType`, `CorrelationId`, `Label`, `MessageId`, `PartitionKey`, `ReplyTo`, `ReplyToSessionId`, `ScheduledEnqueueTimeUtc`, `SessionId`, `TimeToLive` and `To`.

---

A `serialization` block supports the following:

* `type` - (Required) The serialization format used for outgoing data streams. Possible values are `Avro`, `Csv`, `Json` and `Parquet`.

* `encoding` - (Optional) The encoding of the incoming data in the case of input and the encoding of outgoing data in the case of output. It currently can only be set to `UTF8`.

-> **Note:** This is required when `type` is set to `Csv` or `Json`.

* `field_delimiter` - (Optional) The delimiter that will be used to separate comma-separated value (CSV) records. Possible values are ` ` (space), `,` (comma), `	` (tab), `|` (pipe) and `;`.

-> **Note:** This is required when `type` is set to `Csv`.

* `format` - (Optional) Specifies the format of the JSON the output will be written in. Possible values are `Array` and `LineSeparated`.

-> **Note:** This is Required and can only be specified when `type` is set to `Json`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics Output ServiceBus Topic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics Output ServiceBus Topic.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics Output ServiceBus Topic.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics Output ServiceBus Topic.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics Output ServiceBus Topic.

## Import

Stream Analytics Output ServiceBus Topic's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_servicebus_topic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/outputs/output1
```
