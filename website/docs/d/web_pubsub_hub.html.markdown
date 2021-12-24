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
  name          = "exampleWpsHub"
  web_pubsub_id = "example-wps"
}

output "web_pubsub_hub_id" {
  value = data.azurerm_web_pubsub_hub.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specify the name of this Web Pubsub Hub.

* `web_pubsub_id` - (Required) Specify the id of the Web Pubsub. Changing this forces a new resource to be created.

## Attributes Reference

* `id` - The ID of the Resource.

* `name` -  The Name of the Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Web Pubsub Hub.

