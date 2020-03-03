---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_topic"
description: |-
  Gets information about an existing EventGrid Topic

---

# Data Source: azurerm_eventgrid_topic

Use this data source to access information about an existing EventGrid Topic

## Example Usage

```hcl
data "azurerm_eventgrid_topic" "example" {
  name                = "my-eventgrid-topic"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the EventGrid Topic resource.

* `resource_group_name` - The name of the resource group in which the EventGrid Topic exists.

## Attributes Reference

The following attributes are exported:

* `id` - The EventGrid Topic ID.

* `endpoint` - The Endpoint associated with the EventGrid Topic.

* `primary_access_key` - The Primary Shared Access Key associated with the EventGrid Topic.

* `secondary_access_key` - The Secondary Shared Access Key associated with the EventGrid Topic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Topic.
