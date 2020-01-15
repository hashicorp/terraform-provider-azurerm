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

* `name` - (Required) The name of the EventGrid Topic resource.

* `resource_group_name` - (Required) The name of the resource group in which the EventGrid Topic exists.

## Attributes Reference

The following attributes are exported:

* `id` - The EventGrid Topic ID.

* `endpoint` - The Endpoint associated with the EventGrid Topic.

* `primary_access_key` - The Primary Shared Access Key associated with the EventGrid Topic.

* `secondary_access_key` - The Secondary Shared Access Key associated with the EventGrid Topic.
