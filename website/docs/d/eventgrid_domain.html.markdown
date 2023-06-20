---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_domain"
description: |-
  Gets information about an existing EventGrid Domain

---

# Data Source: azurerm_eventgrid_domain

Use this data source to access information about an existing EventGrid Domain

## Example Usage

```hcl
data "azurerm_eventgrid_domain" "example" {
  name                = "my-eventgrid-domain"
  resource_group_name = "example-resources"
}

output "eventgrid_domain_mapping_topic" {
  value = data.azurerm_eventgrid_domain.example.input_mapping_fields.0.topic
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the EventGrid Domain resource.

* `resource_group_name` - (Required) The name of the resource group in which the EventGrid Domain exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the EventGrid Domain.

* `location` - The Azure Region in which this EventGrid Domain exists.

* `endpoint` - The Endpoint associated with the EventGrid Domain.

* `primary_access_key` - The primary access key associated with the EventGrid Domain.

* `secondary_access_key` - The secondary access key associated with the EventGrid Domain.

* `input_schema` - The schema in which incoming events will be published to this domain. Possible values are `CloudEventSchemaV1_0`, `CustomEventSchema`, or `EventGridSchema`.

* `input_mapping_fields` - A `input_mapping_fields` block as defined below.

* `input_mapping_default_values` - A `input_mapping_default_values` block as defined below.

* `public_network_access_enabled` - Whether or not public network access is allowed for this server.

* `inbound_ip_rule` - One or more `inbound_ip_rule` blocks as defined below.

* `tags` - A mapping of tags assigned to the EventGrid Domain.

---

A `input_mapping_fields` supports the following:

* `id` - Specifies the id of the EventGrid Event associated with the domain.

* `topic` - Specifies the topic of the EventGrid Event associated with the domain.

* `event_type` - Specifies the event type of the EventGrid Event associated with the domain.

* `event_time` - Specifies the event time of the EventGrid Event associated with the domain.

* `data_version` - Specifies the data version of the EventGrid Event associated with the domain.

* `subject` - Specifies the subject of the EventGrid Event associated with the domain.

---

A `input_mapping_default_values` supports the following:

* `event_type` - Specifies the default event type of the EventGrid Event associated with the domain.

* `data_version` - Specifies the default data version of the EventGrid Event associated with the domain.

* `subject` - Specifies the default subject of the EventGrid Event associated with the domain.

---

A `inbound_ip_rule` block supports the following:

* `ip_mask` - The IP mask (CIDR) to match on.

* `action` - The action to take when the rule is matched. Possible values are `Allow`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Domain.
