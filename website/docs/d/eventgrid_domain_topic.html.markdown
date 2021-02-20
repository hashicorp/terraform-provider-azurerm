---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_domain_topic"
description: |-
  Gets information about an existing EventGrid Domain Topic

---

# Data Source: azurerm_eventgrid_domain_topic

Use this data source to access information about an existing EventGrid Domain Topic

## Example Usage

```hcl
data "azurerm_eventgrid_domain_topic" "example" {
  name                = "my-eventgrid-domain-topic"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the EventGrid Domain Topic resource.

* `domain_name` - The name of the EventGrid Domain Topic domain.

* `resource_group_name` - The name of the resource group in which the EventGrid Domain Topic exists.

## Attributes Reference

The following attributes are exported:

* `id` - The EventGrid Domain Topic ID.

* `domain_name` - The EventGrid Domain Topic Domain name.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Domain Topic.
