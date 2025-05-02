---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_data_collection_endpoint"
description: |-
  Get information about the specified Data Collection Endpoint.

---

# Data Source: azurerm_monitor_data_collection_endpoint

Use this data source to access information about an existing Data Collection Endpoint.

## Example Usage

```hcl
data "azurerm_monitor_data_collection_endpoint" "example" {
  name                = "example-mdce"
  resource_group_name = azurerm_resource_group.example.name
}

output "endpoint_id" {
  value = data.azurerm_monitor_data_collection_endpoint.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Data Collection Endpoint.

* `resource_group_name` - Specifies the name of the resource group the Data Collection Endpoint is located in.

## Attributes Reference

* `id` - The ID of the Resource.

* `configuration_access_endpoint` - The endpoint used for accessing configuration, e.g., `https://mydce-abcd.eastus-1.control.monitor.azure.com`.

* `description` - Specifies a description for the Data Collection Endpoint.

* `immutable_id` - The immutable ID of the Data Collection Endpoint.

* `kind` - The kind of the Data Collection Endpoint. Possible values are `Linux` and `Windows`.

* `location` - The Azure Region where the Data Collection Endpoint should exist.

* `logs_ingestion_endpoint` - The endpoint used for ingesting logs, e.g., `https://mydce-abcd.eastus-1.ingest.monitor.azure.com`.

* `metrics_ingestion_endpoint` - The endpoint used for ingesting metrics, e.g., `https://mydce-abcd.eastus-1.metrics.ingest.monitor.azure.com`.

* `public_network_access_enabled` - Whether network access from public internet to the Data Collection Endpoint are allowed. Possible values are `true` and `false`.

* `tags` - A mapping of tags which should be assigned to the Data Collection Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Collection Endpoint.
