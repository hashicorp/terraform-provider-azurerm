---
subcategory: "Web Pubsub"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_hub"
description: |-
  Gets information about an existing Azure Web Pubsub Hub service.
---

# Data Source: azurerm_web_pubsub_hub

Use this data source to access information about an existing Azure Web Pubsub Hub service.

## Example Usage

```hcl
data "azurerm_web_pubsub_hub" "test" {
  name                = "exampleWpsHub"
  web_pubsub_name     = "example-wps"
  resource_group_name = "example-RG"
}

output "web_pubsub_hub_id" {
  value = data.azurerm_web_pubsub_hub.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Web Pubsub Hub.

* `resource_group_name` - (Required) The name of the Resource Group where the Web Pubsub Hub exists.

* `web_pubsub_name` - (Required) The name of the Web Pubsub where the Web Pubsub Hub exists.

## Attributes Reference

* `id` - The ID of the Resource.

* `name` -  The Name of the Resource.

* `event_handler` - A list of the event handler configured for this Resource.

* `anonymous_connect_policy` - The anonymous connections configured for this Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Web Pubsub Hub.

