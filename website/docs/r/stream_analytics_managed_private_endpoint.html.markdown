---
subcategory: "Stream Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stream_analytics_managed_private_endpoint"
description: |-
  Manages a Stream Analytics Managed Private Endpoint.
---

# azurerm_stream_analytics_managed_private_endpoint

Manages a Stream Analytics Managed Private Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_stream_analytics_cluster" "example" {
  name                = "examplestreamanalyticscluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  streaming_capacity  = 36
}

resource "azurerm_stream_analytics_managed_private_endpoint" "example" {
  name                          = "exampleprivateendpoint"
  resource_group_name           = azurerm_resource_group.example.name
  stream_analytics_cluster_name = azurerm_stream_analytics_cluster.example.name
  target_resource_id            = azurerm_storage_account.example.id
  subresource_name              = "blob"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Stream Analytics Managed Private Endpoint. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Stream Analytics Managed Private Endpoint should exist. Changing this forces a new resource to be created.

* `stream_analytics_cluster_name` - (Required) The name of the Stream Analytics Cluster where the Managed Private Endpoint should be created. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the Private Link Enabled Remote Resource which this Stream Analytics Private endpoint should be connected to. Changing this forces a new resource to be created.

* `subresource_name` - (Required) Specifies the sub resource name which the Stream Analytics Private Endpoint is able to connect to. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Stream Analytics.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Stream Analytics.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stream Analytics.
* `delete` - (Defaults to 5 minutes) Used when deleting the Stream Analytics.

## Import

Stream Analytics Private Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stream_analytics_managed_private_endpoint.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/clusters/cluster1/privateEndpoints/endpoint1
```
