---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_partner_namespace"
description: |-
  Gets information about an existing Event Grid Partner Namespace

---

# Data Source: azurerm_eventgrid_partner_namespace

Use this data source to access information about an existing Event Grid Partner Namespace

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

* `name` - (Required) The name of the Event Grid Partner Namespace resource.

* `resource_group_name` - (Required) The name of the resource group in which the Event Grid Partner Namespace exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Event Grid Partner Namespace.

* `location` - The Azure Region where the Event Grid Partner Namespace exists.

* `endpoint` - The endpoint for the Event Grid Partner Namespace.

* `inbound_ip_rule` - One or more `inbound_ip_rule` blocks as defined below.

* `local_authentication_enabled` - Whether local authentication methods are enabled for the Event Grid Partner Namespace.

* `partner_registration_id` - The resource Id of the partner registration associated with this Event Grid Partner Namespace.

* `partner_topic_routing_mode` - The partner topic routing mode.

* `public_network_access` - Whether or not public network access is allowed for this server.

* `tags` - A mapping of tags which are assigned to the Event Grid Partner Namespace.

---

A `inbound_ip_rule` block supports the following:

* `ip_mask` - The IP mask (CIDR) to match on.

* `action` - The action to take when the rule is matched.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Event Grid Partner Namespace.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.EventGrid` - 2022-06-15
