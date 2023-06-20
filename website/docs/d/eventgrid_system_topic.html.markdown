---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_system_topic"
description: |-
  Gets information about an existing EventGrid System Topic

---

# Data Source: azurerm_eventgrid_system_topic

Use this data source to access information about an existing EventGrid System Topic

## Example Usage

```hcl
data "azurerm_eventgrid_system_topic" "example" {
  name                = "eventgrid-system-topic"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the EventGrid System Topic resource.

* `resource_group_name` - The name of the resource group in which the EventGrid System Topic exists.

## Attributes Reference

The following attributes are exported:

* `id` - The EventGrid System Topic ID.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Event Grid System Topic.

* `metric_arm_resource_id` - The Metric ARM Resource ID of the Event Grid System Topic.

* `source_arm_resource_id` - The ID of the Event Grid System Topic ARM Source.

* `topic_type` - The Topic Type of the Event Grid System Topic.

* `tags` - A mapping of tags which are assigned to the Event Grid System Topic.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Event Grid System Topic.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Event Grid System Topic.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Event Grid System Topic.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Event Grid System Topic.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid System Topic.
