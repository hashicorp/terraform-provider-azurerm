---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_data_collection_endpoint"
description: |-
  Manages a Data Collection Endpoint.
---

# azurerm_monitor_data_collection_endpoint

Manages a Data Collection Endpoint.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_monitor_data_collection_endpoint" "example" {
  name                          = "example-mdce"
  resource_group_name           = azurerm_resource_group.example.name
  location                      = azurerm_resource_group.example.location
  kind                          = "Windows"
  public_network_access_enabled = true
  description                   = "monitor_data_collection_endpoint example"
  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Data Collection Endpoint should exist. Changing this forces a new Data Collection Endpoint to be created.

* `name` - (Required) The name which should be used for this Data Collection Endpoint. Changing this forces a new Data Collection Endpoint to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Data Collection Endpoint should exist. Changing this forces a new Data Collection Endpoint to be created.

---

* `description` - (Optional) Specifies a description for the Data Collection Endpoint.

* `kind` - (Optional) The kind of the Data Collection Endpoint. Possible values are `Linux` and `Windows`.

* `public_network_access_enabled` - (Optional) Whether network access from public internet to the Data Collection Endpoint are allowed. Possible values are `true` and `false`. Default to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Data Collection Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Collection Endpoint.

* `configuration_access_endpoint` - The endpoint used for accessing configuration, e.g., `https://mydce-abcd.eastus-1.control.monitor.azure.com`.

* `immutable_id` - The immutable ID of the Data Collection Endpoint.

* `logs_ingestion_endpoint` - The endpoint used for ingesting logs, e.g., `https://mydce-abcd.eastus-1.ingest.monitor.azure.com`.

* `metrics_ingestion_endpoint` - The endpoint used for ingesting metrics, e.g., `https://mydce-abcd.eastus-1.metrics.ingest.monitor.azure.com`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Collection Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Collection Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Data Collection Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Collection Endpoint.

## Import

Data Collection Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_data_collection_endpoint.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/dataCollectionEndpoints/endpoint1
```
