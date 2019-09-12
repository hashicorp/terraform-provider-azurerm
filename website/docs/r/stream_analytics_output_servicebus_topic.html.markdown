---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_servicebus_topic"
sidebar_current: "docs-azurerm-resource-stream-analytics-output-servicebus-topic"
description: |-
  Manages a Stream Analytics Output to a ServiceBus Topic.
---

# azurerm_stream_analytics_output_servicebus_topic

Manages a Stream Analytics Output to a ServiceBus Topic.

## Example Usage

```hcl
data "azurerm_resource_group" "example" {
  name = "example-resources"
}

data "azurerm_stream_analytics_job" "example" {
  name                = "example-job"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "example-namespace"
  location            = "${data.azurerm_resource_group.example.location}"
  resource_group_name = "${data.azurerm_resource_group.example.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "example" {
  name                = "example-topic"
  resource_group_name = "${data.azurerm_resource_group.example.name}"
  namespace_name      = "${azurerm_servicebus_namespace.example.name}"
  enable_partitioning = true
}

resource "azurerm_stream_analytics_output_servicebus_topic" "test" {
  name                      = "blob-storage-output"
  stream_analytics_job_name = "${data.azurerm_stream_analytics_job.example.name}"
  resource_group_name       = "${data.azurerm_stream_analytics_job.example.resource_group_name}"
  topic_name                = "${azurerm_servicebus_topic.example.name}"
  servicebus_namespace      = "${azurerm_servicebus_namespace.example.name}"
  shared_access_policy_key  = "${azurerm_servicebus_namespace.example.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    format = "Avro"
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

* `shared_access_policy_key` - (Required) The shared access policy key for the specified shared access policy.

* `shared_access_policy_name` - (Required) The shared access policy name for the Event Hub, Service Bus Queue, Service Bus Topic, etc.

* `serialization` - (Required) A `serialization` block as defined below.

---

A `serialization` block supports the following:

* `type` - (Required) The serialization format used for outgoing data streams. Possible values are `Avro`, `Csv` and `Json`.

* `encoding` - (Optional) The encoding of the incoming data in the case of input and the encoding of outgoing data in the case of output. It currently can only be set to `UTF8`.

-> **NOTE:** This is required when `type` is set to `Csv` or `Json`.

* `field_delimiter` - (Optional) The delimiter that will be used to separate comma-separated value (CSV) records. Possible values are ` ` (space), `,` (comma), `   ` (tab), `|` (pipe) and `;`.

-> **NOTE:** This is required when `type` is set to `Csv`.

* `format` - (Optional) Specifies the format of the JSON the output will be written in. Possible values are `Array` and `LineSeparated`.

-> **NOTE:** This is Required and can only be specified when `type` is set to `Json`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Stream Analytics Output ServiceBus Topic.

## Import

Stream Analytics Output ServiceBus Topic's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_servicebus_topic.test /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.StreamAnalytics/streamingjobs/job1/outputs/output1
```
