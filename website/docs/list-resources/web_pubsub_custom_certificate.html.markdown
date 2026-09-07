---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_custom_certificate"
description: |-
    Lists Web PubSub Custom Certificate resources.
---

# List resource: azurerm_web_pubsub_custom_certificate

Lists Web PubSub Custom Certificate resources.

## Example Usage

### List Web PubSub Custom Certificates in a Web PubSub service

```hcl
list "azurerm_web_pubsub_custom_certificate" "example" {
  provider = azurerm
  config {
    web_pubsub_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.SignalRService/webPubSub/mywps"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `web_pubsub_id` - (Required) The ID of the Web PubSub service to query.
