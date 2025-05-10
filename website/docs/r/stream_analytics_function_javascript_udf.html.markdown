---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_function_javascript_udf"
description: |-
  Manages a JavaScript UDF Function within Stream Analytics Streaming Job.
---

# azurerm_stream_analytics_function_javascript_udf

Manages a JavaScript UDF Function within Stream Analytics Streaming Job.

## Example Usage

```hcl
data "azurerm_resource_group" "example" {
  name = "example-resources"
}

data "azurerm_stream_analytics_job" "example" {
  name                = "example-job"
  resource_group_name = data.azurerm_resource_group.example.name
}

resource "azurerm_stream_analytics_function_javascript_udf" "example" {
  name                      = "example-javascript-function"
  stream_analytics_job_name = data.azurerm_stream_analytics_job.example.name
  resource_group_name       = data.azurerm_stream_analytics_job.example.resource_group_name

  script = <<SCRIPT
function getRandomNumber(in) {
  return in;
}
SCRIPT


  input {
    type = "bigint"
  }

  output {
    type = "bigint"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the JavaScript UDF Function. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Job exists. Changing this forces a new resource to be created.

* `stream_analytics_job_name` - (Required) The name of the Stream Analytics Job where this Function should be created. Changing this forces a new resource to be created.

* `input` - (Required) One or more `input` blocks as defined below.

* `output` - (Required) An `output` blocks as defined below.

* `script` - (Required) The JavaScript of this UDF Function.

---

A `input` block supports the following:

* `type` - (Required) The Data Type for the Input Argument of this JavaScript Function. Possible values include `array`, `any`, `bigint`, `datetime`, `float`, `nvarchar(max)` and `record`.

* `configuration_parameter` - (Optional) Is this input parameter a configuration parameter? Defaults to `false`.

---

A `output` block supports the following:

* `type` - (Required) The Data Type output from this JavaScript Function. Possible values include `array`, `any`, `bigint`, `datetime`, `float`, `nvarchar(max)` and `record`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics JavaScript UDF Function.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics JavaScript UDF Function.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics JavaScript UDF Function.
* `update` - (Defaults to 30 minutes) Used when updating the Stream Analytics JavaScript UDF Function.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stream Analytics JavaScript UDF Function.

## Import

Stream Analytics JavaScript UDF Functions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_function_javascript_udf.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StreamAnalytics/streamingJobs/job1/functions/func1
```
