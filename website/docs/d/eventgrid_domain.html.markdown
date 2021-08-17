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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the EventGrid Domain resource.

* `resource_group_name` - (Required) The name of the resource group in which the EventGrid Domain exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the EventGrid Domain.

* `location` - The Azure Region in which this EventGrid Domain exists.

* `endpoint` - The Endpoint associated with the EventGrid Domain.

* `input_schema` - The schema in which incoming events will be published to this domain. Possible values are `CloudEventSchemaV1_0`, `CustomEventSchema`, or `EventGridSchema`.

* `input_mapping_fields` - 
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Domain.

