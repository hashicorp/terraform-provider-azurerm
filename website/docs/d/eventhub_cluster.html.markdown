---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_eventhub_cluster"
description: |-
  Gets information about an existing EventHub Cluster.
---

# Data Source: azurerm_eventhub_cluster

Use this data source to access information about an existing EventHub.

## Example Usage

```hcl
data "azurerm_eventhub_cluster" "example" {
  name                = "search-eventhub"
  resource_group_name = "search-service"
}

output "eventhub_id" {
  value = data.azurerm_eventhub_cluster.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this EventHub Cluster.

* `resource_group_name` - (Required) The name of the Resource Group where the EventHub Cluster exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the EventHub Cluster.

* `sku_name` - SKU name of the EventHub Cluster.

* `location` - Location of the EventHub Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Cluster.
