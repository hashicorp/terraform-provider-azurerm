---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_system_topic"
description: |-
  Gets information about an existing Event Grid System Topic

---

# Data Source: azurerm_eventgrid_system_topic

Use this data source to access information about an existing Event Grid System Topic

## Example Usage

```hcl
data "azurerm_eventgrid_system_topic" "example" {
  name                = "eventgrid-system-topic"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Event Grid System Topic resource.

* `resource_group_name` - (Required) The name of the resource group in which the Event Grid System Topic exists.

## Attribute Reference

The following attributes are exported:

* `id` - The Event Grid System Topic ID.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Event Grid System Topic.

* `metric_resource_id` - The Metric Resource ID of the Event Grid System Topic.

-> **Note:** This is **not** an Azure RM ID ("/subscription/..."), but rather an Azure-internal identifier for this metric in the form of a GUID. For consumption in Azure Monitor resources, generally the system topic's Azure RM ID is used.

* `source_resource_id` - The ID of the Event Grid System Topic ARM Source.

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Event Grid System Topic.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.EventGrid` - 2025-02-15
