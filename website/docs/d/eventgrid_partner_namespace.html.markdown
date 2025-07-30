---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_namespace"
description: |-
  Gets information about an existing EventGrid Partner Namespace

---

# Data Source: azurerm_eventgrid_partner_namespace

Use this data source to access information about an existing EventGrid Partner Namespace

## Example Usage

```hcl
data "azurerm_eventgrid_partner_namespace" "example" {
  name                = "my-eventgrid-partner-namespace"
  resource_group_name = "example-resources"
}

output "eventgrid_partner_namespace_endpoint" {
  value = data.azurerm_eventgrid_partner_namespace.example.endpoint
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the EventGrid Partner Namespace resource.

* `resource_group_name` - (Required) The name of the resource group in which the EventGrid Partner Namespace exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the EventGrid Partner Namespace.

* `location` - The Azure Region where the EventGrid Partner Namespace exists.

* `endpoint` - The endpoint for the EventGrid Partner Namespace.

* `partner_registration_id` - The fully qualified ARM Id of the partner registration that should be associated with this partner namespace.

* `partner_topic_routing_mode` - The partner topic routing mode.

* `tags` - A mapping of tags which are assigned to the EventGrid Partner Namespace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Partner Namespace.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.EventGrid`: 2022-06-15
