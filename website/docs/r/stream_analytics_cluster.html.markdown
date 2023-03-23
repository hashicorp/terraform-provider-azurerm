---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_cluster"
description: |-
  Manages a Stream Analytics Cluster.
---

# azurerm_stream_analytics_cluster

Manages a Stream Analytics Cluster.

## Example Usage

```hcl

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_stream_analytics_cluster" "example" {
  name                = "examplestreamanalyticscluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  streaming_capacity  = 36
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Stream Analytics Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Stream Analytics Cluster should exist. Changing this forces a new resource to be created.

* `streaming_capacity` - (Required) The number of streaming units supported by the Cluster. Accepted values are multiples of `36` in the range of `36` to `216`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Stream Analytics.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Stream Analytics.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics.
* `update` - (Defaults to 90 minutes) Used when updating the Stream Analytics.
* `delete` - (Defaults to 90 minutes) Used when deleting the Stream Analytics.

## Import

Stream Analytics Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_cluster.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/clusters/cluster1
```
