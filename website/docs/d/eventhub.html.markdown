---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_eventhub"
description: |-
  Gets information about an existing Event Hub.
---

# Data Source: azurerm_eventhub

Use this data source to access information about an existing Event Hub.

## Example Usage

```hcl
data "azurerm_eventhub" "example" {
  name = "search-eventhub"
  resource_group_name = "search-service"
  namespace_name = "search-eventhubns"
}

output "eventhub_id" {
  value = data.azurerm_eventhub.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Event Hub.

* `resource_group_name` - (Required) The name of the Resource Group where the Event Hub exists.

* `namespace_name` - (Required) The name of the EventHub Namespace where the Event Hub exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Event Hub.

* `partition_count` - The number of partitions in the Event Hub.

* `partition_ids` - The identifiers for partitions of this Event Hub.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Event Hub.