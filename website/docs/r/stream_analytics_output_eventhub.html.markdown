---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_output_eventhub"
sidebar_current: "docs-azurerm-resource-stream-analytics-output-eventhub"
description: |-
  Manages a Stream Analytics Output to an EventHub.
---

# azurerm_stream_analytics_output_eventhub

Manages a Stream Analytics Output to an EventHub.

## Example Usage

```hcl
data "azurerm_resource_group" "example" {
  name = "example-resources"
}

data "azurerm_stream_analytics_job" "example" {
  name                = "example-job"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

XXXX

resource "azurerm_stream_analytics_output_blob" "example" {
  name                      = "output-to-blob-storage"
  stream_analytics_job_name = "${data.azurerm_stream_analytics_job.example.name}"
  resource_group_name       = "${data.azurerm_stream_analytics_job.example.resource_group_name}"
  storage_account_name      = "${azurerm_storage_account.example.name}"
  storage_account_key       = "${azurerm_storage_account.example.primary_access_key}"
  storage_container_name    = "${azurerm_storage_container.example.name}"
  path_pattern              = "some-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stream Output. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job. Changing this forces a new resource to be created.

* `eventhub_name` - (Required) The name of the Event Hub.

* `servicebus_namespace` - (Required) The namespace that is associated with the desired Event Hub, Service Bus Queue, Service Bus Topic, etc.

* `shared_access_policy_key` - (Required) The shared access policy key for the specified shared access policy.

* `shared_access_policy_name` - (Required) The shared access policy name for the Event Hub, Service Bus Queue, Service Bus Topic, etc.

* `serialization` - (Required) A `serialization` block as defined below.

---

A `serialization` block supports the following:

* `type` - (Required) The serialization format used for outgoing data streams. Possible values are `Avro`, `Csv` and `Json`.

* `encoding` - (Optional) The encoding of the incoming data in the case of input and the encoding of outgoing data in the case of output.

-> **NOTE:** This is required when `type` is set to `Csv` or `Json`.

* `field_delimiter` - (Optional) The delimiter that will be used to separate comma-separated value (CSV) records. Possible values are ` ` (space), `,` (comma), `   ` (tab), `|` (pipe) and `;`.

-> **NOTE:** This is required when `type` is set to `Csv`.

* `format` - (Optional) Specifies the format of the JSON the output will be written in.

-> **NOTE:** This is Required and can only be specified when `type` is set to `Json`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Stream Analytics Output EventHub.

## Import

Stream Analytics Outputs to an EventHub can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_output_eventhub.test /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.StreamAnalytics/streamingjobs/job1/outputs/output1
```
