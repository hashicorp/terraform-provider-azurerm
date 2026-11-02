---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_custom_domain"
description: |-
  Lists Web PubSub Custom Domain resources.
---

# List resource: azurerm_web_pubsub_custom_domain

Lists Web PubSub Custom Domain resources.

## Example Usage

### List all Web PubSub Custom Domains for a Web PubSub

```hcl
list "azurerm_web_pubsub_custom_domain" "example" {
  provider = azurerm
  config {
    web_pubsub_id = azurerm_web_pubsub.example.id
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `web_pubsub_id` - (Required) The ID of the Web PubSub for which Custom Domains should be listed.
